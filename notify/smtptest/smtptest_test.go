package smtptest

import (
	"bufio"
	"errors"
	"net"
	"net/smtp"
	"net/textproto"
	"strings"
	"testing"

	"github.com/vdatacloud/daml-escrow-commons/notify"
)

func start(t *testing.T, opts ...Option) *MockServer {
	t.Helper()
	s, err := NewMockServer("test", "127.0.0.1:0", opts...)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

const multipartMail = "From: Tripart <no-reply@example.test>\r\n" +
	"To: ada@example.test\r\n" +
	"Subject: =?utf-8?q?Un_nouvel_appareil_=C3=A9tait_ajout=C3=A9?=\r\n" +
	"MIME-Version: 1.0\r\n" +
	"Content-Type: multipart/alternative; boundary=\"b1\"\r\n\r\n" +
	"--b1\r\nContent-Type: text/plain; charset=utf-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n" +
	"Cancel: https://portal.example.test/login\r\n" +
	".a line that starts with a dot\r\n" +
	"--b1\r\nContent-Type: text/html; charset=utf-8\r\nContent-Transfer-Encoding: base64\r\n\r\n" +
	"PHA+PGEgaHJlZj0iaHR0cHM6Ly9w\r\nb3J0YWwuZXhhbXBsZS50ZXN0L2xvZ2luIj5DYW5jZWw8L2E+PC9wPg==\r\n" +
	"--b1--\r\n"

func TestMockServer_RecordsWhatNetSMTPSends(t *testing.T) {
	var seen []Message
	s := start(t, WithOnMessage(func(m Message) { seen = append(seen, m) }))
	if err := smtp.SendMail(s.Addr(), nil, "no-reply@example.test", []string{"ada@example.test", "bob@example.test"}, []byte(multipartMail)); err != nil {
		t.Fatalf("SendMail: %v", err)
	}
	msgs := s.Messages()
	if len(msgs) != 1 || len(seen) != 1 {
		t.Fatalf("messages %d, seen %d", len(msgs), len(seen))
	}
	m := msgs[0]
	if m.From != "no-reply@example.test" || strings.Join(m.To, ",") != "ada@example.test,bob@example.test" {
		t.Fatalf("envelope %q %q", m.From, m.To)
	}
	c, err := m.Content()
	if err != nil {
		t.Fatal(err)
	}
	if c.Subject != "Un nouvel appareil était ajouté" {
		t.Errorf("subject %q", c.Subject)
	}
	if !strings.Contains(c.Text, "https://portal.example.test/login") || !strings.Contains(c.Text, "\n.a line that starts with a dot") {
		t.Errorf("text %q (dot-stuffing undone)", c.Text)
	}
	if c.HTML != `<p><a href="https://portal.example.test/login">Cancel</a></p>` {
		t.Errorf("html %q", c.HTML)
	}
}

func TestMockServer_SinglePart(t *testing.T) {
	s := start(t)
	body := "Subject: hi\r\n\r\nplain body\r\n"
	if err := smtp.SendMail(s.Addr(), nil, "a@example.test", []string{"b@example.test"}, []byte(body)); err != nil {
		t.Fatal(err)
	}
	c, err := s.Messages()[0].Content()
	if err != nil || c.Subject != "hi" || c.Text != "plain body\n" || c.HTML != "" {
		t.Fatalf("%+v, %v", c, err)
	}
}

func TestMockServer_FailWith(t *testing.T) {
	s := start(t)
	for _, code := range []int{550, 451} {
		s.FailWith(code)
		err := smtp.SendMail(s.Addr(), nil, "a@example.test", []string{"b@example.test"}, []byte("Subject: x\r\n\r\nx\r\n"))
		var tp *textproto.Error
		if !errors.As(err, &tp) || tp.Code != code {
			t.Fatalf("FailWith(%d): %v", code, err)
		}
	}
	s.FailWith(0)
	if err := smtp.SendMail(s.Addr(), nil, "a@example.test", []string{"b@example.test"}, []byte("Subject: x\r\n\r\nx\r\n")); err != nil {
		t.Fatal(err)
	}
	if len(s.Messages()) != 1 {
		t.Fatalf("messages %d", len(s.Messages()))
	}
}

func TestMockServer_TooBig(t *testing.T) {
	s := start(t, WithMaxSize(64))
	err := smtp.SendMail(s.Addr(), nil, "a@example.test", []string{"b@example.test"}, []byte("Subject: x\r\n\r\n"+strings.Repeat("y", 200)+"\r\n"))
	var tp *textproto.Error
	if !errors.As(err, &tp) || tp.Code != 552 || len(s.Messages()) != 0 {
		t.Fatalf("want 552 and nothing kept: %v, %d", err, len(s.Messages()))
	}
}

// The conversation's rules, by hand: order, syntax, and what isn't offered.
func TestMockServer_Protocol(t *testing.T) {
	s := start(t)
	conn, err := net.Dial("tcp", s.Addr())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	r := bufio.NewReader(conn)
	read := func() string {
		l, err := r.ReadString('\n')
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimRight(l, "\r\n")
	}
	say := func(cmd, wantPrefix string) {
		t.Helper()
		if _, err := conn.Write([]byte(cmd + "\r\n")); err != nil {
			t.Fatal(err)
		}
		if got := read(); !strings.HasPrefix(got, wantPrefix) {
			t.Fatalf("%s: got %q, want %s", cmd, got, wantPrefix)
		}
	}
	if g := read(); !strings.HasPrefix(g, "220") {
		t.Fatalf("greeting %q", g)
	}
	say("HELO me", "250")
	say("RCPT TO:<b@example.test>", "503")
	say("DATA", "503")
	say("MAIL FROM:a@example.test", "501")
	say("MAIL FROM:<a@example.test> BODY=8BITMIME", "250")
	say("RCPT b@example.test", "501")
	say("STARTTLS", "502")
	say("AUTH PLAIN x", "502")
	say("NOOP", "250")
	say("RSET", "250")
	say("RCPT TO:<b@example.test>", "503")
	say("BOGUS", "502")
	say("QUIT", "221")
}

func TestNewMockServer_DevOnly(t *testing.T) {
	if _, err := NewMockServer("production", "127.0.0.1:0"); !errors.Is(err, notify.ErrNotInDev) {
		t.Fatalf("production: %v", err)
	}
	if _, err := NewMockServer("dev", "not an address"); err == nil {
		t.Fatal("a bad address listened")
	}
	s := start(t)
	if s.Host() != "127.0.0.1" || s.Port() == 0 {
		t.Fatalf("addr %s %d", s.Host(), s.Port())
	}
}
