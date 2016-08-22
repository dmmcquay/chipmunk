package chipmunk

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdd(t *testing.T) {
	defer clearContacts()
	sm := http.NewServeMux()
	AddRoutes(sm)
	ts := httptest.NewServer(sm)

	//adding normal user
	u := fmt.Sprintf("%s%s?name=derek&number=1234", ts.URL, prefix["add"])
	req, err := http.NewRequest("POST", u, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("couldn't POST: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		buf := &bytes.Buffer{}
		io.Copy(buf, resp.Body)
		t.Logf("%s", buf.Bytes())
		t.Fatalf("bad request got incorrect status: got %d, want %d", got, want)
	}

	//testing poorly formed query
	u = fmt.Sprintf("%s%s?this=derek&shouldfail=1234", ts.URL, prefix["add"])
	req, err = http.NewRequest("POST", u, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("couldn't POST: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusBadRequest; got != want {
		buf := &bytes.Buffer{}
		io.Copy(buf, resp.Body)
		t.Logf("%s", buf.Bytes())
		t.Fatalf("bad request got incorrect status: got %d, want %d", got, want)
	}

	//adding same user twice
	u = fmt.Sprintf("%s%s?name=derek&number=1234", ts.URL, prefix["add"])
	req, err = http.NewRequest("POST", u, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("couldn't POST: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusBadRequest; got != want {
		buf := &bytes.Buffer{}
		io.Copy(buf, resp.Body)
		t.Logf("%s", buf.Bytes())
		t.Fatalf("bad request got incorrect status: got %d, want %d", got, want)
	}
}

func TestDel(t *testing.T) {
	defer clearContacts()
	sm := http.NewServeMux()
	AddRoutes(sm)
	ts := httptest.NewServer(sm)

	//adding normal user
	u := fmt.Sprintf("%s%s?name=derek&number=1234", ts.URL, prefix["add"])
	req, err := http.NewRequest("POST", u, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("couldn't POST: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		buf := &bytes.Buffer{}
		io.Copy(buf, resp.Body)
		t.Logf("%s", buf.Bytes())
		t.Fatalf("bad request got incorrect status: got %d, want %d", got, want)
	}

	//remove user that doesn't exists
	u = fmt.Sprintf("%s%s?name=doesnotexist", ts.URL, prefix["delete"])
	req, err = http.NewRequest("DELETE", u, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("couldn't DELETE: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusNoContent; got != want {
		buf := &bytes.Buffer{}
		io.Copy(buf, resp.Body)
		t.Logf("%s", buf.Bytes())
		t.Fatalf("bad request got incorrect status: got %d, want %d", got, want)
	}

	//remove user that exists
	u = fmt.Sprintf("%s%s?name=derek", ts.URL, prefix["delete"])
	req, err = http.NewRequest("DELETE", u, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("couldn't DELETE: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		buf := &bytes.Buffer{}
		io.Copy(buf, resp.Body)
		t.Logf("%s", buf.Bytes())
		t.Fatalf("bad request got incorrect status: got %d, want %d", got, want)
	}
	// attempt to get deleted user, should fail
	u = fmt.Sprintf("%s%s?name=derek", ts.URL, prefix["search"])
	req, err = http.NewRequest("GET", u, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("couldn't GET: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusNoContent; got != want {
		buf := &bytes.Buffer{}
		io.Copy(buf, resp.Body)
		t.Logf("%s", buf.Bytes())
		t.Fatalf("bad request got incorrect status: got %d, want %d", got, want)
	}
}

func TestEdit(t *testing.T) {
	defer clearContacts()
	sm := http.NewServeMux()
	AddRoutes(sm)
	ts := httptest.NewServer(sm)

	//adding normal user
	u := fmt.Sprintf("%s%s?name=derek&number=1234", ts.URL, prefix["add"])
	req, err := http.NewRequest("POST", u, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("couldn't POST: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		buf := &bytes.Buffer{}
		io.Copy(buf, resp.Body)
		t.Logf("%s", buf.Bytes())
		t.Fatalf("bad request got incorrect status: got %d, want %d", got, want)
	}
	//edit user
	u = fmt.Sprintf("%s%s?name=derek&edit=drock", ts.URL, prefix["edit"])
	req, err = http.NewRequest("POST", u, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("couldn't POST: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		buf := &bytes.Buffer{}
		io.Copy(buf, resp.Body)
		t.Logf("%s", buf.Bytes())
		t.Fatalf("bad request got incorrect status: got %d, want %d", got, want)
	}
	// get the edited user
	u = fmt.Sprintf("%s%s?name=drock", ts.URL, prefix["search"])
	req, err = http.NewRequest("GET", u, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("couldn't GET: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		buf := &bytes.Buffer{}
		io.Copy(buf, resp.Body)
		t.Logf("%s", buf.Bytes())
		t.Fatalf("bad request got incorrect status: got %d, want %d", got, want)
	}
}
