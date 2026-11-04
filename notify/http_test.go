package notify

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func staticToken(context.Context) (string, error) { return "svc-token", nil }

func TestHTTPSender_SendsAndReadsTheReceipt(t *testing.T) {
	var got Message
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/notifications" || r.Header.Get("Authorization") != "Bearer svc-token" ||
			r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("request %s %s auth %q", r.Method, r.URL.Path, r.Header.Get("Authorization"))
		}
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"id":"n-1","status":"accepted"}`))
	}))
	defer srv.Close()

	s, err := NewHTTPSender(srv.URL+"/", staticToken)
	if err != nil {
		t.Fatal(err)
	}
	r, err := s.Send(context.Background(), validMessage())
	if err != nil || r != (Receipt{ID: "n-1", Status: StatusAccepted}) {
		t.Fatalf("Send: %+v, %v", r, err)
	}
	if got.IdempotencyKey != "identity:reenrolment:42" || got.Recipient.Subject != "00u1abcd" || got.Params["deviceName"] != "MacBook" {
		t.Fatalf("posted %+v", got)
	}
}

func TestHTTPSender_StatusAnswersToTheRelay(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
		want   error
	}{
		{"repeat", http.StatusOK, `{"id":"n-1","status":"sent"}`, nil},
		{"invalid", http.StatusBadRequest, `{"error":"template must be a dotted id","status":400}`, ErrInvalidMessage},
		{"key reused", http.StatusConflict, `{"error":"this idempotency key was used for a different notice","status":409}`, ErrRejected},
		{"unknown template", http.StatusUnprocessableEntity, `{"error":"unknown template","status":422}`, ErrRejected},
		{"no scope", http.StatusForbidden, `{"error":"platform:notify is required"}`, ErrRejected},
		{"unauthenticated", http.StatusUnauthorized, "unauthorized", ErrRejected},
		{"off", http.StatusServiceUnavailable, `{"error":"notifications are not configured","code":"","status":503}`, ErrUnavailable},
		{"database down", http.StatusBadGateway, `{"error":"the notification couldn't be recorded"}`, ErrUnavailable},
		{"rate limited", http.StatusTooManyRequests, "", ErrUnavailable},
		{"garbled receipt", http.StatusAccepted, `{"status":"accepted"}`, ErrUnavailable},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			}))
			defer srv.Close()
			s, _ := NewHTTPSender(srv.URL, staticToken)
			_, err := s.Send(context.Background(), validMessage())
			if tc.want == nil {
				if err != nil {
					t.Fatalf("got %v", err)
				}
				return
			}
			if !errors.Is(err, tc.want) {
				t.Fatalf("got %v, want %v", err, tc.want)
			}
			if strings.Contains(tc.body, `"error"`) && !strings.Contains(err.Error(), "status") && !strings.Contains(err.Error(), ":") {
				t.Fatalf("the platform's message is lost: %v", err)
			}
		})
	}
}

func TestHTTPSender_RefusesAnInvalidMessageWithoutCalling(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true }))
	defer srv.Close()
	s, _ := NewHTTPSender(srv.URL, staticToken)
	bad := validMessage()
	bad.Template = "Hello"
	if _, err := s.Send(context.Background(), bad); !errors.Is(err, ErrInvalidMessage) || called {
		t.Fatalf("got %v (called %v)", err, called)
	}
}

func TestHTTPSender_TransportFailuresAreRetryable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	url := srv.URL
	srv.Close()
	s, _ := NewHTTPSender(url, staticToken)
	if _, err := s.Send(context.Background(), validMessage()); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("a closed server: %v, want ErrUnavailable", err)
	}

	noToken, _ := NewHTTPSender(url, func(context.Context) (string, error) { return "", errors.New("okta down") })
	if _, err := noToken.Send(context.Background(), validMessage()); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("no token: %v, want ErrUnavailable", err)
	}

	slow := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { time.Sleep(200 * time.Millisecond) }))
	defer slow.Close()
	s2, _ := NewHTTPSender(slow.URL, staticToken)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if _, err := s2.Send(ctx, validMessage()); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("the caller's deadline: %v", err)
	}
}

func TestHTTPSender_Status(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/notifications/n-1" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(`{"id":"n-1","status":"delivered"}`))
	}))
	defer srv.Close()
	s, _ := NewHTTPSender(srv.URL, staticToken, WithHTTPClient(srv.Client()))
	if r, err := s.Status(context.Background(), "n-1"); err != nil || r.Status != StatusDelivered {
		t.Fatalf("Status: %+v, %v", r, err)
	}
	if _, err := s.Status(context.Background(), "n-2"); !errors.Is(err, ErrRejected) {
		t.Fatalf("someone else's: %v", err)
	}
	if _, err := s.Status(context.Background(), ""); !errors.Is(err, ErrInvalidMessage) {
		t.Fatalf("no id: %v", err)
	}
}

func TestNewHTTPSender_Refusals(t *testing.T) {
	for _, u := range []string{"", "platform.internal", "ftp://platform", "http://"} {
		if _, err := NewHTTPSender(u, staticToken); err == nil {
			t.Errorf("NewHTTPSender(%q) accepted", u)
		}
	}
	if _, err := NewHTTPSender("https://platform.internal", nil); err == nil {
		t.Error("no token source accepted")
	}
}
