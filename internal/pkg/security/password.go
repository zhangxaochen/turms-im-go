package security

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	PrefixBcrypt       = "{bcrypt}"
	PrefixSaltedSha256 = "{salted_sha256}"
	PrefixNoop         = "{noop}"
)

// MatchesPassword simulates Java Turms PasswordManager.matchesPassword()
func MatchesPassword(rawPassword string, encodedPassword string) bool {
	if encodedPassword == "" {
		return rawPassword == ""
	}

	if strings.HasPrefix(encodedPassword, PrefixBcrypt) {
		hash := encodedPassword[len(PrefixBcrypt):]
		return bcrypt.CompareHashAndPassword([]byte(hash), []byte(rawPassword)) == nil
	} else if strings.HasPrefix(encodedPassword, PrefixSaltedSha256) {
		// Mock implementation for SALTED_SHA256: `{salted_sha256}salt:hash` or similar.
		// Without exact Java format, we do a basic fallback or skip if format unknown.
		// If custom turms implementation does `salt+rawPassword` hashing:
		parts := strings.SplitN(encodedPassword[len(PrefixSaltedSha256):], ":", 2)
		if len(parts) == 2 {
			salt := parts[0]
			expectedHash := parts[1]
			hash := sha256.Sum256([]byte(salt + rawPassword))
			return hex.EncodeToString(hash[:]) == expectedHash
		}
		// Fallback for safety
		return false
	} else if strings.HasPrefix(encodedPassword, PrefixNoop) {
		return encodedPassword[len(PrefixNoop):] == rawPassword
	}

	// Default fallback (like Turms does): plain equal check
	return encodedPassword == rawPassword
}
