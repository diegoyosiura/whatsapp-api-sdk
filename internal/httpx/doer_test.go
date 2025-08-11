package httpx

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

/*** Helpers ***/

// tinyBackoffDoer cria um Doer com backoff mínimo para testes rápidos.
func tinyBackoffDoer(rt http.RoundTripper, maxRetries int) *Doer {
	return New(Options{
		Transport:   rt,
		MaxRetries:  maxRetries,
		BaseBackoff: 1 * time.Millisecond,
		MaxBackoff:  2 * time.Millisecond,
	})
}

// timeoutNetErr é um erro que implementa net.Error com Timeout()=true.
type timeoutNetErr struct{ msg string }

func (e timeoutNetErr) Error() string   { return e.msg }
func (e timeoutNetErr) Timeout() bool   { return true }
func (e timeoutNetErr) Temporary() bool { return true }

// flipFlopRT retorna timeout na 1ª chamada e uma resposta 200 na 2ª.
type flipFlopRT struct {
	calls int
}

func (rt *flipFlopRT) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.calls++
	if rt.calls == 1 {
		return nil, timeoutNetErr{"dial tcp: i/o timeout"}
	}
	// Segunda chamada: sucesso sem rede
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(`ok`)),
		Request:    req,
	}, nil
}

/*** Tests ***/

func TestDefaultRetryPolicy(t *testing.T) {
	// 429 => retry
	if !DefaultRetryPolicy(http.StatusTooManyRequests) {
		t.Fatalf("expected retry for 429")
	}
	// 500/502/503 => retry
	for _, st := range []int{500, 502, 503} {
		if !DefaultRetryPolicy(st) {
			t.Fatalf("expected retry for %d", st)
		}
	}
	// 501/505 => NO retry
	for _, st := range []int{http.StatusNotImplemented, http.StatusHTTPVersionNotSupported} {
		if DefaultRetryPolicy(st) {
			t.Fatalf("did not expect retry for %d", st)
		}
	}
	// 400/401/403/404 => NO retry
	for _, st := range []int{400, 401, 403, 404} {
		if DefaultRetryPolicy(st) {
			t.Fatalf("did not expect retry for %d", st)
		}
	}
}

func TestDo_RetriesOn429Then200(t *testing.T) {
	var calls int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			io.WriteString(w, `{"error":"rate limited"}`)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"ok":true}`)
	}))
	defer ts.Close()

	doer := tinyBackoffDoer(ts.Client().Transport, 2)

	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	resp, err := doer.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Do() error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 after retry, got %d", resp.StatusCode)
	}
	if calls < 2 {
		t.Fatalf("expected at least 2 calls, got %d", calls)
	}
}

func TestDo_StopsOnNonRetryable4xx(t *testing.T) {
	var calls int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error":"bad request"}`)
	}))
	defer ts.Close()

	doer := tinyBackoffDoer(ts.Client().Transport, 3)

	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	resp, err := doer.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Do() error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
	if calls != 1 {
		t.Fatalf("expected no retries on 400, calls=%d", calls)
	}
}

func TestDo_MaxRetriesAndReturnsLastResponse(t *testing.T) {
	var calls int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusServiceUnavailable) // 503 (retryable)
		io.WriteString(w, `{"error":"unavailable"}`)
	}))
	defer ts.Close()

	doer := tinyBackoffDoer(ts.Client().Transport, 2) // 1ª + 2 retries = 3 chamadas

	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	resp, err := doer.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Do() error: %v", err)
	}
	defer resp.Body.Close()

	if calls != 3 {
		t.Fatalf("expected 3 attempts, got %d", calls)
	}
	if resp.StatusCode != 503 {
		t.Fatalf("expected final 503, got %d", resp.StatusCode)
	}
}

func TestDo_ErrorOnNonRewindableBodyWhenRetrying(t *testing.T) {
	var calls int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		// Força retry (500)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `oops`)
	}))
	defer ts.Close()

	doer := tinyBackoffDoer(ts.Client().Transport, 1)

	// Body SEM GetBody => não dá para rebobinar; deve falhar ao tentar retry.
	body := io.NopCloser(bytes.NewBufferString(`payload`))
	req, _ := http.NewRequest(http.MethodPost, ts.URL, body)
	// (req.GetBody == nil)

	_, err := doer.Do(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error due to non-rewindable body, got nil")
	}
	if !strings.Contains(err.Error(), "not rewindable") {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("should fail before 2nd attempt, calls=%d", calls)
	}
}

func TestDo_RewindsBodyAndSucceeds(t *testing.T) {
	var calls int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			w.WriteHeader(http.StatusInternalServerError) // força retry
			io.WriteString(w, `first fail`)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `ok`)
	}))
	defer ts.Close()

	doer := tinyBackoffDoer(ts.Client().Transport, 2)

	orig := []byte(`payload`)
	req, _ := http.NewRequest(http.MethodPost, ts.URL, io.NopCloser(bytes.NewReader(orig)))
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(orig)), nil
	}

	resp, err := doer.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Do() error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 after retry, got %d", resp.StatusCode)
	}
	if calls < 2 {
		t.Fatalf("expected retry path, calls=%d", calls)
	}
}

func TestDo_RespectsContextDeadline(t *testing.T) {
	// Servidor que demora além do timeout do contexto.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	doer := tinyBackoffDoer(ts.Client().Transport, 3)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	_, err := doer.Do(ctx, req)
	if err == nil {
		t.Fatalf("expected context deadline error")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestDo_RetryOnNetworkTimeoutThenOK(t *testing.T) {
	rt := &flipFlopRT{}
	doer := tinyBackoffDoer(rt, 2)

	req, _ := http.NewRequest(http.MethodGet, "http://example.local/whatever", nil)
	resp, err := doer.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error after retry: %v", err)
	}
	resp.Body.Close()
	if rt.calls != 2 {
		t.Fatalf("expected 2 calls (timeout then ok), got %d", rt.calls)
	}
}

/*** Safety: ensure isTempOrTimeout only flags net.Error timeouts ***/
func Test_isTempOrTimeout_OnlyForNetError(t *testing.T) {
	if isTempOrTimeout(errors.New("plain error")) {
		t.Fatalf("plain error should not be considered temp/timeout")
	}
	var ne net.Error = timeoutNetErr{"x"}
	if !isTempOrTimeout(ne) {
		t.Fatalf("net.Error timeout should be considered temp/timeout")
	}
}

func TestDo_NilRequest(t *testing.T) {
	d := New(Options{})
	if _, err := d.Do(context.Background(), nil); err == nil || !strings.Contains(err.Error(), "nil request") {
		t.Fatalf("expected nil request error, got %v", err)
	}
}

func TestDo_RewindBodyError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	d := tinyBackoffDoer(ts.Client().Transport, 1)
	req, _ := http.NewRequest(http.MethodPost, ts.URL, io.NopCloser(bytes.NewBufferString("x")))
	req.GetBody = func() (io.ReadCloser, error) { return nil, errors.New("boom") }
	if _, err := d.Do(context.Background(), req); err == nil || !strings.Contains(err.Error(), "rewind body") {
		t.Fatalf("expected rewind error, got %v", err)
	}
}

type errRT struct{}

func (errRT) RoundTrip(*http.Request) (*http.Response, error) { return nil, errors.New("fail") }

func TestDo_NonRetryableNetworkError(t *testing.T) {
	d := tinyBackoffDoer(errRT{}, 3)
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if _, err := d.Do(context.Background(), req); err == nil || !strings.Contains(err.Error(), "fail") {
		t.Fatalf("expected network error, got %v", err)
	}
}

func TestDo_ReachesFinalReturn(t *testing.T) {
	d := New(Options{})
	d.maxRetries = -1 // force loop skip
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	resp, err := d.Do(context.Background(), req)
	if resp != nil || err != nil {
		t.Fatalf("expected nil resp and err, got %v %v", resp, err)
	}
}
