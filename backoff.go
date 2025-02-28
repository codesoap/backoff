// Package backoff provides a simple and performant exponential backoff
// implementation.
package backoff

// A Limiter limits the amount of calls to a function, if the function
// fails. Controlled by its DoubleAfterNFails and MaxSkipsBeforeRetry
// fields, the amount of calls to the function become increasingly rare,
// if the function continues to fail.
type Limiter struct {
	// DoubleAfterNFails determines how many times an action has to fail in
	// a row before the amount of skips between attempts is doubled.
	//
	// If DoubleAfterNFails is 0 or smaller, it will be treated as being 1.
	DoubleAfterNFails int

	// MaxSkipsBeforeRetry defines the maximum amount of times an action is
	// skipped before retrying it. Will be ignored if it is 0.
	MaxSkipsBeforeRetry int

	skip        int
	failsInARow int
}

// Try calls action if no limit is currently active. If action returns
// true, any progression in the limit will be reset. If action returns
// false, the next call(s) of Try will not call action. If action
// continuously fails, the calls of action will become increasingly rare
// according to l.DoubleAfterNFails and l.MaxSkipsBeforeRetry.
//
// Try returns true, if action has been called and false, if a limit
// prevented the call.
func (l *Limiter) Try(action func() bool) bool {
	if l.skip > 0 {
		l.skip--
		return false
	}
	if success := action(); success {
		l.Reset()
	} else {
		doubleAfterNFails := max(l.DoubleAfterNFails, 1)
		l.skip = 1 << (l.failsInARow / doubleAfterNFails)
		if l.MaxSkipsBeforeRetry > 0 {
			l.skip = min(l.skip, l.MaxSkipsBeforeRetry)
		}
		l.failsInARow++
	}
	return true
}

// Reset resets any progression in the limits and ensures that the
// passed action is not skipped at the next call of l.Try.
func (l *Limiter) Reset() {
	l.skip = 0
	l.failsInARow = 0
}
