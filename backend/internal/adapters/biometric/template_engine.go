package biometric

import (
	"context"
	"crypto/sha256"
	"fmt"
)

// TemplateEngine implements ports.BiometricMatcher.
// This is a placeholder implementation. The real fingerprint matching SDK
// (e.g., SourceAFIS, Neurotechnology MegaMatcher) will be integrated later.
type TemplateEngine struct{}

// NewTemplateEngine creates a new biometric template engine.
func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{}
}

// Match compares two fingerprint templates and returns a match score between 0 and 1.
// This placeholder performs a simple byte-level comparison using hash similarity.
// A production implementation would use a proper fingerprint minutiae matching algorithm.
func (e *TemplateEngine) Match(_ context.Context, template1, template2 []byte) (float64, error) {
	if len(template1) == 0 || len(template2) == 0 {
		return 0, fmt.Errorf("templates must not be empty")
	}

	hash1 := sha256.Sum256(template1)
	hash2 := sha256.Sum256(template2)

	// Exact match check via hash comparison.
	if hash1 == hash2 {
		return 1.0, nil
	}

	// Simple byte-level similarity for non-identical templates.
	// In production, this would be replaced by minutiae-based matching.
	minLen := len(template1)
	if len(template2) < minLen {
		minLen = len(template2)
	}

	matching := 0
	for i := 0; i < minLen; i++ {
		if template1[i] == template2[i] {
			matching++
		}
	}

	maxLen := len(template1)
	if len(template2) > maxLen {
		maxLen = len(template2)
	}

	score := float64(matching) / float64(maxLen)
	return score, nil
}

// ExtractTemplate extracts a fingerprint template from a raw image.
// This placeholder returns a SHA-256 hash of the input image as the "template".
// A production implementation would use a fingerprint feature extraction SDK.
func (e *TemplateEngine) ExtractTemplate(_ context.Context, rawImage []byte) ([]byte, error) {
	if len(rawImage) == 0 {
		return nil, fmt.Errorf("raw image must not be empty")
	}

	// Placeholder: use the hash of the raw image as a template.
	// Real implementation would extract minutiae points.
	hash := sha256.Sum256(rawImage)
	return hash[:], nil
}
