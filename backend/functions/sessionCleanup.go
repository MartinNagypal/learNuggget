package functions

import "time"

// StartSessionCleanup deletes expired sessions once right away (in case the
// server was offline at midnight) and then every night at 00:00 server time.
// It runs in its own goroutine and returns immediately.
func StartSessionCleanup() {
	go func() {
		// Errors are logged inside DeleteExpiredSessions; the next run tries again.
		DeleteExpiredSessions()
		for {
			time.Sleep(time.Until(nextMidnight(time.Now())))
			DeleteExpiredSessions()
		}
	}()
}

// nextMidnight returns the next 00:00 in the local time zone after now.
// time.Date normalizes day+1 across month/year ends and handles DST changes.
func nextMidnight(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
}
