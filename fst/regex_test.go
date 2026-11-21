package fst

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegexAutomaton_acceptsOrdinaryPatterns(t *testing.T) {
	for _, pattern := range []string{`[a-z]{3,}ing`, `.*(ough|augh).*`, `[a-z]{0,40}`, `\pL{5}`} {
		_, err := NewRegexAutomaton(pattern)
		assert.NoError(t, err, pattern)
	}
}

func TestNewRegexAutomaton_rejectsExpensivePatternsQuickly(t *testing.T) {
	for _, pattern := range []string{`(\pL?\pL?){0,20}`, `(.?.?){0,500}`, `.{0,1000}.{0,1000}`} {
		start := time.Now()
		_, err := NewRegexAutomaton(pattern)
		require.Error(t, err, pattern)
		assert.Less(t, time.Since(start), time.Second, pattern)
	}
}
