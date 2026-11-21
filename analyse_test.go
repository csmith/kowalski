package kowalski

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyseRunLengthEncodingLargeCounts(t *testing.T) {
	for _, input := range []string{"9999999999x", "99999999999999999999x"} {
		t.Run(
			input,
			func(t *testing.T) {
				results := analyseRunLengthEncoding(testChecker, input)
				assert.Len(t, results, 1)
				assert.Len(t, results[0], 250)
			},
		)
	}
}
