package errorsx

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestHTTPErrorError(t *testing.T) {
	var he *HTTPError
	if he.Error() != "<nil>" {
		t.Fatalf("nil receiver error != <nil>")
	}
	he = &HTTPError{Method: "GET", URL: "u", StatusCode: 500, Status: "500", FBTraceID: "abc"}
	if !strings.Contains(he.Error(), "fb-trace-id=abc") {
		t.Fatalf("expected fb-trace-id in error: %s", he.Error())
	}
	he.FBTraceID = ""
	if strings.Contains(he.Error(), "fb-trace-id") {
		t.Fatalf("did not expect fb-trace-id in error: %s", he.Error())
	}
}

func TestIsRetryable(t *testing.T) {
	cases := []struct {
		code int
		want bool
	}{
		{http.StatusTooManyRequests, true},
		{500, true},
		{http.StatusNotImplemented, false},
		{http.StatusHTTPVersionNotSupported, false},
		{400, false},
	}
	for _, tc := range cases {
		if got := IsRetryable(&HTTPError{StatusCode: tc.code}); got != tc.want {
			t.Fatalf("code %d: want %v got %v", tc.code, tc.want, got)
		}
	}
	if IsRetryable(errors.New("x")) {
		t.Fatalf("non HTTPError should not be retryable")
	}
}

func TestNewHTTPErrorFromResponse(t *testing.T) {
	he := NewHTTPErrorFromResponse(nil, []byte("b"))
	if he.StatusCode != 0 || he.Status != "<nil response>" {
		t.Fatalf("unexpected nil response handling: %+v", he)
	}

	req, _ := http.NewRequest(http.MethodGet, "http://x", nil)
	resp := &http.Response{StatusCode: 404, Status: "404", Header: http.Header{"X-Fb-Trace-Id": {"id"}}, Request: req}
	he = NewHTTPErrorFromResponse(resp, []byte("b"))
	if he.FBTraceID != "id" {
		t.Fatalf("expected fbtrace id set, got %s", he.FBTraceID)
	}
}

func TestGraphError(t *testing.T) {
	var ge *GraphError
	if ge.Error() != "<nil>" {
		t.Fatalf("nil graph error != <nil>")
	}
	req, _ := http.NewRequest(http.MethodGet, "http://x", nil)
	resp := &http.Response{StatusCode: 400, Status: "400", Header: http.Header{"Content-Type": {"application/json"}}, Request: req}
	body := []byte(`{"error":{"message":"bad","type":"x","code":10,"fbtrace_id":"tid"}}`)
	ge = TryParseGraphError(resp, body)
	if ge.Detail.Message != "bad" || ge.HTTP.FBTraceID != "tid" {
		t.Fatalf("unexpected parse: %+v", ge)
	}
	if !errors.Is(ge, ge.HTTP) {
		t.Fatalf("unwrap not working")
	}
	if !strings.Contains(ge.Error(), "bad") {
		t.Fatalf("error string missing detail: %s", ge.Error())
	}

	// Invalid JSON path -> Raw preserved and message empty
	badBody := []byte("not-json")
	ge2 := TryParseGraphError(resp, badBody)
	if ge2.Detail.Message != "" || string(ge2.Raw) != string(badBody) {
		t.Fatalf("unexpected fallback parse: %+v", ge2)
	}
	if !strings.Contains(ge2.Error(), "undecoded") {
		t.Fatalf("unexpected error string: %s", ge2.Error())
	}
}

func TestValidationErrorError(t *testing.T) {
	var v *ValidationError
	if v.Error() != "<nil>" {
		t.Fatalf("nil validation error != <nil>")
	}
	v = &ValidationError{Field: "f", Reason: "r", Op: "op"}
	if !strings.Contains(v.Error(), "field=f") {
		t.Fatalf("unexpected error: %s", v.Error())
	}
	v = &ValidationError{Reason: "r", Op: "op"}
	if strings.Contains(v.Error(), "field=") {
		t.Fatalf("did not expect field info: %s", v.Error())
	}
}
