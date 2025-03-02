package backoff_test

import (
	"fmt"

	"github.com/codesoap/backoff"
)

func Example_failProgression() {
	myActionThatMightFail := func() bool {
		// Emulate an action that always fails for demonstration purposes:
		return false
	}

	limiter := backoff.FailLimiter{}
	for i := 1; i <= 12; i++ {
		tried := limiter.Try(myActionThatMightFail)
		if tried {
			fmt.Printf("#%02d: Executed action.\n", i)
		} else {
			fmt.Printf("#%02d: Skipped action.\n", i)
		}
	}
	// Output:
	// #01: Executed action.
	// #02: Skipped action.
	// #03: Executed action.
	// #04: Skipped action.
	// #05: Skipped action.
	// #06: Executed action.
	// #07: Skipped action.
	// #08: Skipped action.
	// #09: Skipped action.
	// #10: Skipped action.
	// #11: Executed action.
	// #12: Skipped action.
}

func Example_cappedFailProgression() {
	myActionThatMightFail := func() bool {
		// Emulate an action that always fails for demonstration purposes:
		return false
	}

	limiter := backoff.FailLimiter{SkipLimit: 2}
	for i := 1; i <= 12; i++ {
		tried := limiter.Try(myActionThatMightFail)
		if tried {
			fmt.Printf("#%02d: Executed action.\n", i)
		} else {
			fmt.Printf("#%02d: Skipped action.\n", i)
		}
	}
	// Output:
	// #01: Executed action.
	// #02: Skipped action.
	// #03: Executed action.
	// #04: Skipped action.
	// #05: Skipped action.
	// #06: Executed action.
	// #07: Skipped action.
	// #08: Skipped action.
	// #09: Executed action.
	// #10: Skipped action.
	// #11: Skipped action.
	// #12: Executed action.
}

func Example_slowFailProgression() {
	myActionThatMightFail := func() bool {
		// Emulate an action that always fails for demonstration purposes:
		return false
	}

	limiter := backoff.FailLimiter{BackoffInterval: 2}
	for i := 1; i <= 12; i++ {
		tried := limiter.Try(myActionThatMightFail)
		if tried {
			fmt.Printf("#%02d: Executed action.\n", i)
		} else {
			fmt.Printf("#%02d: Skipped action.\n", i)
		}
	}
	// Output:
	// #01: Executed action.
	// #02: Skipped action.
	// #03: Executed action.
	// #04: Skipped action.
	// #05: Executed action.
	// #06: Skipped action.
	// #07: Skipped action.
	// #08: Executed action.
	// #09: Skipped action.
	// #10: Skipped action.
	// #11: Executed action.
	// #12: Skipped action.
}

func Example_withError() {
	// Emulate a simple action function for demonstration purposes:
	myActionThatMightFail := func(someParameter int) error {
		if someParameter > 3 {
			return fmt.Errorf("I don't like large numbers")
		}
		return nil
	}

	limiter := backoff.FailLimiter{}
	for i := 1; i <= 12; i++ {
		// If the tried action requires parameters or returns an error or
		// other values, which are needed for further steps, they can be
		// passed out using closures:
		var err error
		tried := limiter.Try(func() bool {
			err = myActionThatMightFail(i)
			return err == nil
		})
		if tried {
			if err != nil {
				fmt.Printf("#%02d: Executed action, but failed: %v\n", i, err)
				continue
			}
			fmt.Printf("#%02d: Executed action and succeeded.\n", i)
			// Do more computations after success...
		} else {
			fmt.Printf("#%02d: Skipped action.\n", i)
		}
	}
	// Output:
	// #01: Executed action and succeeded.
	// #02: Executed action and succeeded.
	// #03: Executed action and succeeded.
	// #04: Executed action, but failed: I don't like large numbers
	// #05: Skipped action.
	// #06: Executed action, but failed: I don't like large numbers
	// #07: Skipped action.
	// #08: Skipped action.
	// #09: Executed action, but failed: I don't like large numbers
	// #10: Skipped action.
	// #11: Skipped action.
	// #12: Skipped action.
}
