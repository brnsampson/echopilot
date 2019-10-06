package echoserver_test

import (
	"fmt"
	"github.com/brnsampson/echopilot/pkg/echoserver"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestWrapHandler(t *testing.T) {
	logger := &testLogger{}
	testString := "testing a thinger"
	data := url.Values{}
	data.Set("data", "foo")
	req, err := http.NewRequest("GET", "/", strings.NewReader(data.Encode()))
	ok(t, err)
	w := httptest.NewRecorder()
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testString)
	}

	newfn := echoserver.WrapHandler(fn, logger)
	newfn(w, req)
	equals(t, w.Code, http.StatusOK)
	assert(t, (1 <= len(logger.debugs)), "WrapHandler should emit at least one debug message!")
}
