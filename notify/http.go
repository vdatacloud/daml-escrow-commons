package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ErrRejected: the platform service refused the message for good (a reused
// idempotency key with a different message, an unknown template, a channel
// with no provider, or the caller lacks platform:notify). Not retryable.
var ErrRejected = errors.New("notify: rejected by the platform service")

// TokenFunc returns the calling service's bearer token for the platform
// service (its own client-credentials token, with platform:notify).
type TokenFunc func(ctx context.Context) (string, error)

// HTTPSender sends through the platform service's
// POST /api/v1/notifications (daml-escrow PLAN.md Phase 101.7). Its errors
// tell a caller's outbox relay what to do: ErrUnavailable means retry with
// the same IdempotencyKey; ErrInvalidMessage and ErrRejected mean stop.
type HTTPSender struct {
	base   *url.URL
	token  TokenFunc
	client *http.Client
}

// HTTPOption configures NewHTTPSender.
type HTTPOption func(*HTTPSender)

// WithHTTPClient sends with c (default: a client with a 15 s timeout).
func WithHTTPClient(c *http.Client) HTTPOption {
	return func(s *HTTPSender) { s.client = c }
}

// NewHTTPSender sends to the platform service at baseURL (for example
// "https://platform.internal"), authenticated by token.
func NewHTTPSender(baseURL string, token TokenFunc, opts ...HTTPOption) (*HTTPSender, error) {
	u, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return nil, fmt.Errorf("notify: the platform service URL must be http(s)://host: %q", baseURL)
	}
	if token == nil {
		return nil, errors.New("notify: a token source is required")
	}
	s := &HTTPSender{base: u, token: token, client: &http.Client{Timeout: 15 * time.Second}}
	for _, o := range opts {
		o(s)
	}
	return s, nil
}

// Send validates msg and posts it. A repeat of an accepted key answers with
// the first receipt.
func (s *HTTPSender) Send(ctx context.Context, msg Message) (Receipt, error) {
	if err := msg.Validate(); err != nil {
		return Receipt{}, err
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return Receipt{}, fmt.Errorf("notify: encoding the message: %w", err)
	}
	return s.do(ctx, http.MethodPost, "/api/v1/notifications", body)
}

// Status is the receipt of a message this caller sent (ErrRejected for an
// id that isn't theirs or doesn't exist).
func (s *HTTPSender) Status(ctx context.Context, id string) (Receipt, error) {
	if id == "" {
		return Receipt{}, fmt.Errorf("%w: an id is required", ErrInvalidMessage)
	}
	return s.do(ctx, http.MethodGet, "/api/v1/notifications/"+url.PathEscape(id), nil)
}

func (s *HTTPSender) do(ctx context.Context, method, path string, body []byte) (Receipt, error) {
	tok, err := s.token(ctx)
	if err != nil {
		return Receipt{}, fmt.Errorf("%w: getting a token: %v", ErrUnavailable, err)
	}
	var rdr io.Reader
	if body != nil {
		rdr = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, s.base.String()+path, rdr)
	if err != nil {
		return Receipt{}, fmt.Errorf("notify: building the request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := s.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return Receipt{}, fmt.Errorf("notify: %w", ctx.Err())
		}
		return Receipt{}, fmt.Errorf("%w: calling the platform service: %v", ErrUnavailable, err)
	}
	defer func() { _ = resp.Body.Close() }()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))

	switch {
	case resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted:
		var r Receipt
		if err := json.Unmarshal(raw, &r); err != nil || r.ID == "" {
			return Receipt{}, fmt.Errorf("%w: an unreadable receipt", ErrUnavailable)
		}
		return r, nil
	case resp.StatusCode == http.StatusBadRequest:
		return Receipt{}, fmt.Errorf("%w: %s", ErrInvalidMessage, refusal(raw, resp.StatusCode))
	case resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500:
		return Receipt{}, fmt.Errorf("%w: %s", ErrUnavailable, refusal(raw, resp.StatusCode))
	default:
		return Receipt{}, fmt.Errorf("%w: %s", ErrRejected, refusal(raw, resp.StatusCode))
	}
}

// refusal is the platform's error envelope message, or the status alone.
func refusal(raw []byte, status int) string {
	var env struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}
	if json.Unmarshal(raw, &env) == nil && env.Error != "" {
		if env.Code != "" {
			return fmt.Sprintf("%d %s: %s", status, env.Code, env.Error)
		}
		return fmt.Sprintf("%d: %s", status, env.Error)
	}
	return fmt.Sprintf("status %d", status)
}

var _ Sender = (*HTTPSender)(nil)
