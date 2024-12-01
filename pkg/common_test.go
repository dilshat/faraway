package pkg

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"
)

func TestFindNonce(t *testing.T) {
	task := "testTask"
	hashPrefix := "abc"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	nonce, _ := FindNonce(ctx, task, hashPrefix)
	if nonce == "" {
		t.Errorf("Expected nonce to be non-empty, but got an empty string")
	}

	// Verify that the nonce produces a hash with the correct prefix
	hash := sha256.Sum256([]byte(task + nonce))
	hashHex := hex.EncodeToString(hash[:])
	if hashPrefix != hashHex[:len(hashPrefix)] {
		t.Errorf("Nonce %s does not produce the correct hash prefix", nonce)
	}
}

func TestVerifySolution(t *testing.T) {
	task := "testTask"
	hashPrefix := "abc"

	// Find the correct nonce using the task and hashPrefix
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	nonce, _ := FindNonce(ctx, task, hashPrefix)

	// Verify that the solution is correct
	if !VerifySolution(task, nonce, hashPrefix) {
		t.Errorf("Expected nonce %s to be correct for task %s and prefix %s", nonce, task, hashPrefix)
	}

	// Test with an incorrect nonce
	incorrectNonce := "9999"
	if VerifySolution(task, incorrectNonce, hashPrefix) {
		t.Errorf("Expected nonce %s to be incorrect for task %s and prefix %s", incorrectNonce, task, hashPrefix)
	}
}
