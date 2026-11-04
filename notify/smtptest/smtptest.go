// Package smtptest is a mock SMTP server: it accepts mail and records it,
// delivering nothing, for tests and local runs of an SMTP provider (the
// platform service's, daml-escrow PLAN.md Phase 101). It speaks enough SMTP
// for Go's net/smtp client and most relays' clients: EHLO/HELO, MAIL, RCPT,
// DATA, RSET, NOOP, QUIT. No TLS and no AUTH: it is a stand-in, so
// NewMockServer refuses outside development.
package smtptest

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"strings"
	"sync"
	"time"

	"github.com/vdatacloud/daml-escrow-commons/notify"
)

// Message is one mail the server accepted.
type Message struct {
	From     string
	To       []string
	Data     []byte // as sent, dot-unstuffed, CRLF line ends
	Received time.Time
}

// Content is a Message's subject and its text and HTML bodies, decoded.
type Content struct {
	Subject string
	Text    string
	HTML    string
}

// Content decodes m: the subject (RFC 2047) and the text/plain and
// text/html parts, single-part or multipart/alternative, quoted-printable
// or base64.
func (m Message) Content() (Content, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(m.Data))
	if err != nil {
		return Content{}, fmt.Errorf("smtptest: reading the message: %w", err)
	}
	var c Content
	if c.Subject, err = new(mime.WordDecoder).DecodeHeader(msg.Header.Get("Subject")); err != nil {
		return Content{}, fmt.Errorf("smtptest: decoding the subject: %w", err)
	}
	err = readPart(&c, msg.Header.Get("Content-Type"), msg.Header.Get("Content-Transfer-Encoding"), msg.Body)
	return c, err
}

func readPart(c *Content, contentType, encoding string, body io.Reader) error {
	if contentType == "" {
		contentType = "text/plain"
	}
	media, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return fmt.Errorf("smtptest: content type: %w", err)
	}
	if strings.HasPrefix(media, "multipart/") {
		r := multipart.NewReader(body, params["boundary"])
		for {
			p, err := r.NextRawPart()
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("smtptest: reading a part: %w", err)
			}
			if err := readPart(c, p.Header.Get("Content-Type"), p.Header.Get("Content-Transfer-Encoding"), p); err != nil {
				return err
			}
		}
	}
	switch strings.ToLower(encoding) {
	case "quoted-printable":
		body = quotedprintable.NewReader(body)
	case "base64":
		body = base64.NewDecoder(base64.StdEncoding, body) // skips line breaks
	}
	raw, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("smtptest: reading a body: %w", err)
	}
	text := strings.ReplaceAll(string(raw), "\r\n", "\n")
	switch media {
	case "text/plain":
		c.Text = text
	case "text/html":
		c.HTML = text
	}
	return nil
}

// MockServer is a running mock SMTP server.
type MockServer struct {
	ln      net.Listener
	maxSize int
	onMsg   func(Message)
	timeout time.Duration

	mu       sync.Mutex
	messages []Message
	failCode int
	conns    sync.WaitGroup
	closed   bool
}

// Option configures NewMockServer.
type Option func(*MockServer)

// WithMaxSize caps a message's size in bytes (default 10 MiB); a bigger one
// is refused with 552.
func WithMaxSize(n int) Option { return func(s *MockServer) { s.maxSize = n } }

// WithOnMessage calls f with each accepted message (on the connection's
// goroutine).
func WithOnMessage(f func(Message)) Option { return func(s *MockServer) { s.onMsg = f } }

// NewMockServer listens on addr ("127.0.0.1:0" picks a port) and serves
// until Close. It refuses outside development (notify.ErrNotInDev).
func NewMockServer(environment, addr string, opts ...Option) (*MockServer, error) {
	if !notify.IsDevelopment(environment) {
		return nil, fmt.Errorf("smtptest: a mock SMTP server: %w", notify.ErrNotInDev)
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("smtptest: listening: %w", err)
	}
	s := &MockServer{ln: ln, maxSize: 10 << 20, timeout: time.Minute}
	for _, o := range opts {
		o(s)
	}
	go s.serve()
	return s, nil
}

// Addr is the address the server listens on.
func (s *MockServer) Addr() string { return s.ln.Addr().String() }

// Host is the host part of Addr, for an SMTP client's config.
func (s *MockServer) Host() string {
	h, _, _ := net.SplitHostPort(s.Addr())
	return h
}

// Port is the port part of Addr.
func (s *MockServer) Port() int { return s.ln.Addr().(*net.TCPAddr).Port }

// Messages are the messages accepted so far, oldest first.
func (s *MockServer) Messages() []Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Message(nil), s.messages...)
}

// FailWith makes every later RCPT answer code (e.g. 550 for a permanent
// refusal, 451 for a temporary one); 0 accepts again.
func (s *MockServer) FailWith(code int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failCode = code
}

// Close stops listening and waits for open sessions to end.
func (s *MockServer) Close() error {
	s.mu.Lock()
	s.closed = true
	s.mu.Unlock()
	err := s.ln.Close()
	s.conns.Wait()
	return err
}

func (s *MockServer) serve() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return
		}
		s.mu.Lock()
		if s.closed {
			s.mu.Unlock()
			_ = conn.Close()
			return
		}
		s.conns.Add(1)
		s.mu.Unlock()
		go func() {
			defer s.conns.Done()
			defer func() { _ = conn.Close() }()
			s.session(conn)
		}()
	}
}

// session is one SMTP conversation.
func (s *MockServer) session(conn net.Conn) {
	r := bufio.NewReader(io.LimitReader(conn, int64(s.maxSize)*2+1<<20))
	w := bufio.NewWriter(conn)
	reply := func(lines ...string) bool {
		for _, l := range lines {
			_, _ = w.WriteString(l + "\r\n")
		}
		return w.Flush() == nil
	}
	if !reply("220 smtptest mock ESMTP (delivers nothing)") {
		return
	}
	var from string
	var to []string
	for {
		_ = conn.SetDeadline(time.Now().Add(s.timeout))
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimRight(line, "\r\n")
		verb, arg, _ := strings.Cut(line, " ")
		switch strings.ToUpper(verb) {
		case "EHLO":
			reply("250-smtptest", "250-8BITMIME", fmt.Sprintf("250 SIZE %d", s.maxSize))
		case "HELO":
			reply("250 smtptest")
		case "MAIL":
			addr, ok := path(arg, "FROM:")
			if !ok {
				reply("501 syntax: MAIL FROM:<address>")
				continue
			}
			from, to = addr, nil
			reply("250 ok")
		case "RCPT":
			addr, ok := path(arg, "TO:")
			switch {
			case !ok:
				reply("501 syntax: RCPT TO:<address>")
			case from == "":
				reply("503 MAIL first")
			default:
				s.mu.Lock()
				code := s.failCode
				s.mu.Unlock()
				if code != 0 {
					reply(fmt.Sprintf("%d refused (smtptest FailWith)", code))
					continue
				}
				to = append(to, addr)
				reply("250 ok")
			}
		case "DATA":
			if len(to) == 0 {
				reply("503 RCPT first")
				continue
			}
			reply("354 end with <CRLF>.<CRLF>")
			data, tooBig, err := readData(r, s.maxSize)
			if err != nil {
				return
			}
			if tooBig {
				reply("552 message too big")
			} else {
				m := Message{From: from, To: to, Data: data, Received: time.Now()}
				s.mu.Lock()
				s.messages = append(s.messages, m)
				n := len(s.messages)
				s.mu.Unlock()
				if s.onMsg != nil {
					s.onMsg(m)
				}
				reply(fmt.Sprintf("250 ok: queued as smtptest-%d", n))
			}
			from, to = "", nil
		case "RSET":
			from, to = "", nil
			reply("250 ok")
		case "NOOP":
			reply("250 ok")
		case "QUIT":
			reply("221 bye")
			return
		case "STARTTLS", "AUTH":
			reply("502 not offered by the mock")
		default:
			reply("502 unknown command")
		}
	}
}

// path is the <address> after prefix in a MAIL or RCPT argument, ignoring
// parameters (BODY=8BITMIME, SIZE=...).
func path(arg, prefix string) (string, bool) {
	if len(arg) < len(prefix) || !strings.EqualFold(arg[:len(prefix)], prefix) {
		return "", false
	}
	rest := strings.TrimSpace(arg[len(prefix):])
	if !strings.HasPrefix(rest, "<") {
		return "", false
	}
	end := strings.IndexByte(rest, '>')
	if end < 0 {
		return "", false
	}
	return rest[1:end], true
}

// readData reads a DATA body up to the lone ".", undoing dot-stuffing.
// Past max bytes it keeps reading to the end but reports tooBig.
func readData(r *bufio.Reader, max int) (data []byte, tooBig bool, err error) {
	var b bytes.Buffer
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, false, err
		}
		if line == ".\r\n" || line == ".\n" {
			return b.Bytes(), tooBig, nil
		}
		line = strings.TrimPrefix(line, ".")
		if b.Len()+len(line) > max {
			tooBig = true
			continue
		}
		b.WriteString(line)
	}
}
