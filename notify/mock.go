package notify

import (
	"context"
	"fmt"
	"sync"
)

// MockSender records the messages it is sent, for tests and local
// development. It validates each message as a real sender does and answers
// a repeated IdempotencyKey with the first receipt, never a second message.
// It is a stand-in: NewMockSender refuses outside development.
type MockSender struct {
	mu       sync.Mutex
	sent     []Message
	receipts map[string]Receipt
	failWith error
}

// NewMockSender is a MockSender for environment, or ErrNotInDev when
// environment isn't a development one (IsDevelopment).
func NewMockSender(environment string) (*MockSender, error) {
	if !IsDevelopment(environment) {
		return nil, fmt.Errorf("%w: environment %q", ErrNotInDev, environment)
	}
	return &MockSender{receipts: map[string]Receipt{}}, nil
}

// Send validates and records msg.
func (m *MockSender) Send(ctx context.Context, msg Message) (Receipt, error) {
	if err := ctx.Err(); err != nil {
		return Receipt{}, fmt.Errorf("notify: send: %w", err)
	}
	if err := msg.Validate(); err != nil {
		return Receipt{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.failWith != nil {
		return Receipt{}, m.failWith
	}
	if r, ok := m.receipts[msg.IdempotencyKey]; ok {
		return r, nil
	}
	r := Receipt{ID: fmt.Sprintf("mock-%d", len(m.sent)+1), Status: StatusAccepted}
	m.receipts[msg.IdempotencyKey] = r
	m.sent = append(m.sent, cloneMessage(msg))
	return r, nil
}

// Sent is every distinct message recorded, in order.
func (m *MockSender) Sent() []Message {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Message, len(m.sent))
	for i, msg := range m.sent {
		out[i] = cloneMessage(msg)
	}
	return out
}

// FailWith makes every later Send fail with err (nil restores sending), to
// test a caller's retries, e.g. FailWith(ErrUnavailable).
func (m *MockSender) FailWith(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failWith = err
}

func cloneMessage(msg Message) Message {
	if msg.Params != nil {
		p := make(map[string]string, len(msg.Params))
		for k, v := range msg.Params {
			p[k] = v
		}
		msg.Params = p
	}
	return msg
}

var _ Sender = (*MockSender)(nil)
