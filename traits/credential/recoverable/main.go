package recoverable

import "time"

// This trait provides features to set and get
// the recovery token in the case the password
// needs to be recovered. The recovery token
// can be set with a duration considering the
// current date time (if the duration is > 0).
// Setting a duration <= 0. When trying to get
// the current recovery token, it will return
// "" if absent or if the current token expired.
// In that case, the token will also be set to
// "" (removing the token).
type Recoverable interface {
	SetRecoveryToken(string, duration time.Duration)
	RecoveryToken() string
}
