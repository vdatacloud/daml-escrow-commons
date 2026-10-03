package apierror

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"regexp"
)

// HeaderRequestID carries the request id in both directions.
const HeaderRequestID = "X-Request-Id"

type requestIDKey struct{}

var requestIDRe = regexp.MustCompile(`^[A-Za-z0-9._:-]{1,128}$`)

// RequestID is middleware giving every request an id: the caller's
// X-Request-Id when it is well-formed, else a fresh one. The id is set on
// the response header and in the context (RequestIDFrom), so an error's
// requestId ties a client report to server logs.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(HeaderRequestID)
		if !requestIDRe.MatchString(id) {
			id = newRequestID()
		}
		w.Header().Set(HeaderRequestID, id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey{}, id)))
	})
}

// RequestIDFrom is the request id RequestID put in ctx, or "".
func RequestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

func newRequestID() string {
	var b [12]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

// Writer sends errors as HTTP responses.
type Writer struct {
	// Detail is the level responses carry. The zero value is Full: set it
	// explicitly -- Summary outside developer environments (see "What an
	// error may disclose" in the package doc).
	Detail Detail
	// Log, when set, is called with each error before it is written --
	// log e.LogFields() there (short forms only).
	Log func(r *http.Request, e *Error)
}

// Write sends e: its status (500 if unset), JSON content type, the request
// id (from RequestID middleware, else the request header) in body and
// header, and the canonical body at w's detail level.
func (wr Writer) Write(w http.ResponseWriter, r *http.Request, e *Error) {
	if e.Status == 0 {
		e.Status = http.StatusInternalServerError
	}
	if e.Code == "" {
		e.Code = CodeForStatus(e.Status)
	}
	if e.RequestID == "" && r != nil {
		if e.RequestID = RequestIDFrom(r.Context()); e.RequestID == "" {
			if id := r.Header.Get(HeaderRequestID); requestIDRe.MatchString(id) {
				e.RequestID = id
			}
		}
	}
	if wr.Log != nil {
		wr.Log(r, e)
	}
	body, err := e.Canonical(wr.Detail)
	if err != nil {
		// Details that can't be encoded must not lose the error itself.
		body, _ = e.Canonical(Summary)
	}
	if e.RequestID != "" {
		w.Header().Set(HeaderRequestID, e.RequestID)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(e.Status)
	_, _ = w.Write(body)
}

// Write sends e at Full detail with no logging (see Writer).
func Write(w http.ResponseWriter, r *http.Request, e *Error) {
	Writer{}.Write(w, r, e)
}

// FromResponse reads a non-2xx response from another platform service
// into an Error (see Parse); the caller closes resp.Body.
func FromResponse(resp *http.Response, body []byte) *Error {
	e := Parse(resp.StatusCode, body)
	if e.RequestID == "" {
		e.RequestID = resp.Header.Get(HeaderRequestID)
	}
	return e
}

func decodeInto(raw json.RawMessage, v any) error { return json.Unmarshal(raw, v) }
