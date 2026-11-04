package notify

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func validMessage() Message {
	return Message{
		IdempotencyKey: "identity:reenrolment:42",
		TenantID:       "ps-acme",
		Channel:        ChannelEmail,
		Recipient:      Recipient{Subject: "00u1abcd"},
		Template:       "identity.reenrolment-started",
		Locale:         "fr-CA",
		Params:         map[string]string{"deviceName": "MacBook", "cancel_url": "https://example.test/c/1"},
	}
}

func TestValidate_AcceptsAWellFormedMessage(t *testing.T) {
	if err := validMessage().Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	sms := validMessage()
	sms.Channel, sms.Recipient, sms.Locale, sms.Params = ChannelSMS, Recipient{AuditRef: "ref:1:s:0123456789abcdef"}, "", nil
	if err := sms.Validate(); err != nil {
		t.Fatalf("Validate (sms, audit ref, no locale or params): %v", err)
	}
}

func TestValidate_Refusals(t *testing.T) {
	cases := []struct {
		name string
		edit func(*Message)
		want string
	}{
		{"no idempotency key", func(m *Message) { m.IdempotencyKey = "" }, "idempotencyKey is required"},
		{"idempotency key with spaces", func(m *Message) { m.IdempotencyKey = "a b" }, "idempotencyKey must be"},
		{"idempotency key too long", func(m *Message) { m.IdempotencyKey = strings.Repeat("k", MaxKeyBytes+1) }, "idempotencyKey must be"},
		{"no tenant", func(m *Message) { m.TenantID = " " }, "tenantId is required"},
		{"unknown channel", func(m *Message) { m.Channel = "fax" }, "channel must be"},
		{"no recipient", func(m *Message) { m.Recipient = Recipient{} }, "needs a subject or an auditRef"},
		{"both recipient refs", func(m *Message) { m.Recipient.AuditRef = "ref:1:s:00" }, "not both"},
		{"recipient too long", func(m *Message) { m.Recipient.Subject = strings.Repeat("s", MaxKeyBytes*2+1) }, "too long"},
		{"free-text template", func(m *Message) { m.Template = "Hello there" }, "template must be"},
		{"undotted template", func(m *Message) { m.Template = "reenrolment" }, "template must be"},
		{"bad locale", func(m *Message) { m.Locale = "french!" }, "locale must be"},
		{"bad param name", func(m *Message) { m.Params = map[string]string{"1x": "v"} }, `param name "1x"`},
		{"param too long", func(m *Message) { m.Params = map[string]string{"x": strings.Repeat("v", MaxParamValueBytes+1)} }, `param "x" is over`},
		{"param with a newline", func(m *Message) { m.Params = map[string]string{"x": "a\nBcc: b"} }, "control characters"},
		{"param with invalid utf-8", func(m *Message) { m.Params = map[string]string{"x": "\xff"} }, "invalid UTF-8"},
		{"too many params", func(m *Message) {
			m.Params = map[string]string{}
			for i := 0; i <= MaxParams; i++ {
				m.Params["p"+strings.Repeat("a", i)] = "v"
			}
		}, "at most"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := validMessage()
			tc.edit(&m)
			err := m.Validate()
			if !errors.Is(err, ErrInvalidMessage) || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("Validate: got %v, want ErrInvalidMessage naming %q", err, tc.want)
			}
		})
	}
}

func TestValidate_NamesEveryBrokenRule(t *testing.T) {
	err := Message{}.Validate()
	for _, want := range []string{"idempotencyKey", "tenantId", "channel", "recipient", "template"} {
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Fatalf("Validate(empty): %v, want it to name %s", err, want)
		}
	}
}

func TestIsDevelopment(t *testing.T) {
	for env, want := range map[string]bool{"dev": true, "Development": true, " local ": true, "test": true, "production": false, "staging": false, "": false} {
		if got := IsDevelopment(env); got != want {
			t.Errorf("IsDevelopment(%q) = %v, want %v", env, got, want)
		}
	}
}

func TestNewMockSender_RefusedOutsideDevelopment(t *testing.T) {
	for _, env := range []string{"production", "staging", ""} {
		if s, err := NewMockSender(env); !errors.Is(err, ErrNotInDev) || s != nil {
			t.Fatalf("NewMockSender(%q): got %v, %v; want ErrNotInDev", env, s, err)
		}
	}
}

func TestMockSender_RecordsOncePerIdempotencyKey(t *testing.T) {
	s, err := NewMockSender("dev")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	first, err := s.Send(ctx, validMessage())
	if err != nil || first.ID == "" || first.Status != StatusAccepted {
		t.Fatalf("Send: %+v, %v", first, err)
	}
	again, err := s.Send(ctx, validMessage())
	if err != nil || again != first {
		t.Fatalf("a repeated key: got %+v, %v; want the first receipt %+v", again, err, first)
	}
	other := validMessage()
	other.IdempotencyKey = "identity:reenrolment:43"
	if r, _ := s.Send(ctx, other); r.ID == first.ID {
		t.Fatalf("a new key got the first message's id %s", r.ID)
	}
	if n := len(s.Sent()); n != 2 {
		t.Fatalf("Sent: %d messages, want 2", n)
	}
}

func TestMockSender_SentIsACopy(t *testing.T) {
	s, _ := NewMockSender("test")
	msg := validMessage()
	if _, err := s.Send(context.Background(), msg); err != nil {
		t.Fatal(err)
	}
	msg.Params["deviceName"] = "changed by the caller"
	got := s.Sent()
	got[0].Params["deviceName"] = "changed by the reader"
	if s.Sent()[0].Params["deviceName"] != "MacBook" {
		t.Fatal("the recorded message changed through a caller's or reader's map")
	}
}

func TestMockSender_ValidatesAndFails(t *testing.T) {
	s, _ := NewMockSender("dev")
	bad := validMessage()
	bad.Template = "free text"
	if _, err := s.Send(context.Background(), bad); !errors.Is(err, ErrInvalidMessage) {
		t.Fatalf("an invalid message: %v, want ErrInvalidMessage", err)
	}

	s.FailWith(ErrUnavailable)
	if _, err := s.Send(context.Background(), validMessage()); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("FailWith: %v, want ErrUnavailable", err)
	}
	s.FailWith(nil)
	if _, err := s.Send(context.Background(), validMessage()); err != nil {
		t.Fatalf("after FailWith(nil): %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := s.Send(ctx, validMessage()); !errors.Is(err, context.Canceled) {
		t.Fatalf("a cancelled context: %v", err)
	}
	if n := len(s.Sent()); n != 1 {
		t.Fatalf("Sent: %d, want only the one accepted message", n)
	}
}
