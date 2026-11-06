package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessChunk(t *testing.T) {
	res, err := processChunk("2 1 abcde")
	require.NoError(t, err)
	assert.Equal(t, []string{"ab", "c", "de"}, res.(map[string]interface{})["result"])
}

func TestProcessChunk_rejectsInvalidSizes(t *testing.T) {
	for _, input := range []string{"0 a", "-1 a", "3 0 abcdef", "3 -2 abcdef"} {
		t.Run(input, func(t *testing.T) {
			_, err := processChunk(input)
			assert.Error(t, err)
		})
	}
}
