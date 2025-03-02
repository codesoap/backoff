// Package backoff provides a simple and performant exponential backoff
// implementation.
package backoff

// A FailLimiter limits the amount of calls to a function, if the
// function fails. Controlled by its BackoffInterval and SkipLimit
// fields, the amount of calls to the function become increasingly rare,
// if the function continues to fail.
type FailLimiter struct {
	// BackoffInterval determines how many times an action has to fail in a
	// row before the amount of skips between attempts is doubled.
	//
	// If BackoffInterval is 0 or smaller, it will be treated as being 1.
	BackoffInterval int

	// SkipLimit defines the maximum amount of times an action is
	// skipped before retrying it. It will be ignored if it is 0.
	SkipLimit int

	skip        int
	failsInARow int
}

// Try calls action if no limit is currently active. If action returns
// true, any progression in the limit will be reset. If action returns
// false, the next call(s) of Try will not call action. If action
// continuously fails, the calls of action will become increasingly rare
// according to fl.BackoffInterval and fl.SkipLimit.
//
// Try returns true, if action has been called and false, if a limit
// prevented the call.
func (fl *FailLimiter) Try(action func() bool) bool {
	if fl.skip > 0 {
		fl.skip--
		return false
	}
	if success := action(); success {
		fl.Reset()
	} else {
		backoffInterval := max(fl.BackoffInterval, 1)
		fl.skip = 1 << (fl.failsInARow / backoffInterval)
		if fl.SkipLimit > 0 {
			fl.skip = min(fl.skip, fl.SkipLimit)
		}
		fl.failsInARow++
	}
	return true
}

// Reset resets any progression in the limits and ensures that the
// passed action is not skipped at the next call of fl.Try.
func (fl *FailLimiter) Reset() {
	fl.skip = 0
	fl.failsInARow = 0
}
