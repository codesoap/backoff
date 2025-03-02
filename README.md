This is a minimalist exponential backoff implementation for Go. It
focuses on absolute simplicity and only contains two functions, of which
you may not even need the second one.

It provides high performance and does not use floating-point arithmetic.

See https://pkg.go.dev/github.com/codesoap/backoff for the full
documentation.

# Usage
Here is a simple example; more can be found in
[the documentation](https://pkg.go.dev/github.com/codesoap/backoff):

```go
package main

import (
	"fmt"
	"github.com/codesoap/backoff"
)

func myActionThatMightFail() bool {
	// Emulate an action that always fails for demonstration purposes:
	return false
}

func main() {
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
```
