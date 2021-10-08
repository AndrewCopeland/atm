package atm_test

import (
	"errors"
	"strings"
	"testing"
)

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error was returned and was not expected. %s", err.Error())
	}
}

func assertError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Error was not returned and was expected")
	}
}

func assertErrorIsError(t *testing.T, err error, errIs error) {
	if err != errIs {
		if err == nil {
			err = errors.New("nil")
		}
		t.Errorf("Unexpected error was returned. '%s' should be '%s'", err.Error(), errIs.Error())
	}
}

func assertErrorContains(t *testing.T, err error, contains string) {
	if !strings.Contains(err.Error(), contains) {
		t.Errorf("Incorrect error, should contains '%s' but is '%s'", contains, err)
	}
}
