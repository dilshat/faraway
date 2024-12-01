package pkg

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// FindNonce generates the nonce for a given task and hash prefix.
func FindNonce(ctx context.Context, task, hashPrefix string) (string, error) {
	var nonce int
	for {
		select {
		case <-ctx.Done():
			// Return the error from the context if it was canceled or timed out
			return "", ctx.Err()
		default:
			// Proceed with the nonce calculation
			combined := fmt.Sprintf("%s%d", task, nonce)
			hash := sha256.Sum256([]byte(combined))
			hashHex := hex.EncodeToString(hash[:])

			if strings.HasPrefix(hashHex, hashPrefix) {
				return fmt.Sprintf("%d", nonce), nil
			}
			nonce++
		}
	}
}

// VerifySolution verifies if the given solution (nonce) is correct for the task and hash prefix.
func VerifySolution(task, nonce, hashPrefix string) bool {
	hash := sha256.Sum256([]byte(task + nonce))
	hashHex := hex.EncodeToString(hash[:])
	return len(hashHex) >= len(hashPrefix) && hashHex[:len(hashPrefix)] == hashPrefix
}
