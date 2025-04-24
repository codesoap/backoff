package backoff_test

import (
	"testing"

	"github.com/codesoap/backoff"
)

func TestReset(t *testing.T) {
	dummyAction := func() bool { return false }
	limiter := &backoff.FailLimiter{}
	limiter.Try(dummyAction)
	limiter.Reset()
	tried := limiter.Try(dummyAction)
	if !tried {
		t.Fatal("Reset did not lead to immediate retry.")
	}
	tried = limiter.Try(dummyAction)
	if tried {
		t.Fatal("Skipping does not work after reset.")
	}
	tried = limiter.Try(dummyAction)
	if !tried {
		t.Fatal("Reset seemingly did not reset fail count.")
	}
}
