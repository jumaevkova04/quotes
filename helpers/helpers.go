package helpers

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// RemoveAllSpaces - remove all blank spaces from string
func RemoveAllSpaces(str string) string {
	return strings.Join(strings.Fields(str), "")
}

// IsValidUUID - check is string valid uuid
func IsValidUUID(s string) bool {
	if _, err := uuid.Parse(s); err != nil {
		return false
	}
	return true
}

// IsTimePassed - if time passed returns true
func IsTimePassed(check, date time.Time) bool {
	return check.After(date)
}
