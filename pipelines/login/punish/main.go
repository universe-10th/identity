package punish

import (
	"fmt"
	"time"
	"github.com/universe-10th/identity/credentials"
	"github.com/universe-10th/identity/credentials/traits/deniable"
	"github.com/universe-10th/identity/credentials/traits/identified"
)

// An instance of this type is returned by the
// PunishmentCheckStep for a credential if it
// is banned.
type PunishedError struct {
	TimeFormat  string
	PunishedOn  time.Time
	PunishedFor *time.Duration
	Reason      interface{}
	PunishedBy  credentials.Credential
}

func (error *PunishedError) Error() string {
	bannedChunk := fmt.Sprintf("banned on: %s", error.PunishedOn.Format(error.TimeFormat))
	untilChunk := " permanently"
	if error.PunishedFor != nil {
		endTime := error.PunishedOn.Add(*error.PunishedFor)
		untilChunk = fmt.Sprintf(" until: %s", endTime.Format(error.TimeFormat))
	}
	reasonChunk := ""
	if error.Reason != nil {
		if stringer, ok := error.Reason.(fmt.Stringer); ok {
			reasonChunk = fmt.Sprintf(" with reason: %s", stringer)
		} else {
			reasonChunk = fmt.Sprintf(" with reason: %v", error.Reason)
		}
	}
	byChunk := ""
	if error.PunishedBy != nil {
		if identifiedTrait, ok := error.PunishedBy.(identified.Identified); ok {
			identification := identifiedTrait.Identification()
			if stringer, ok := identification.(fmt.Stringer); ok {
				reasonChunk = fmt.Sprintf(" by: %s", stringer)
			} else {
				reasonChunk = fmt.Sprintf(" by: %v", error.Reason)
			}
		}
	}

	return bannedChunk + untilChunk + byChunk + reasonChunk
}

// This pipeline step tells when a credential could not
// login because it counts as punished.
type PunishmentCheckStep struct {
	TimeFormat string
}

// Attempts a log-in step which would fail if the credential
// counts as punished.
func (step *PunishmentCheckStep) Login(credential credentials.Credential, password string) error {
	if punishable, ok := credential.(deniable.Punishable); ok {
		if punishedOn, punishedFor, reason, punishedBy := punishable.PunishedFor(); punishedOn != nil {
			return &PunishedError{step.TimeFormat, *punishedOn, punishedFor, reason, punishedBy}
		} else {
			return nil
		}
	} else {
		return nil
	}
}
