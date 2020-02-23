package echoserver_test

import (
	"fmt"
	"github.com/brnsampson/echopilot/pkg/echoserver"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

type testLogger struct {
	debugs []string
	infos  []string
	errors []string
}

func (t *testLogger) Debug(is ...interface{}) {
	for _, i := range is {
		t.debugs = append(t.debugs, fmt.Sprint(i))
	}
}

func (t *testLogger) Info(is ...interface{}) {
	for _, i := range is {
		t.infos = append(t.infos, fmt.Sprint(i))
	}
}

func (t *testLogger) Error(is ...interface{}) {
	for _, i := range is {
		t.errors = append(t.errors, fmt.Sprint(i))
	}
}

func (t *testLogger) Debugf(s string, is ...interface{}) {
	t.debugs = append(t.debugs, fmt.Sprintf(s, is...))
}

func (t *testLogger) Infof(s string, is ...interface{}) {
	t.infos = append(t.infos, fmt.Sprintf(s, is...))
}

func (t *testLogger) Errorf(s string, is ...interface{}) {
	t.errors = append(t.errors, fmt.Sprintf(s, is...))
}

func (t *testLogger) Sync() error {
	fmt.Println("test logger sync'd")
	return nil
}
