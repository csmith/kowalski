package main

import (
	"bytes"
	"testing"

	"github.com/blevesearch/vellum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func useTransducer(t *testing.T, words ...string) {
	var buf bytes.Buffer
	builder, err := vellum.New(&buf, nil)
	require.NoError(t, err)
	for i := range words {
		require.NoError(t, builder.Insert([]byte(words[i]), 1))
	}
	require.NoError(t, builder.Close())

	fst, err := vellum.Load(buf.Bytes())
	require.NoError(t, err)

	old := transducer
	transducer = fst
	t.Cleanup(func() { transducer = old })
}

func TestWordLink_noLeftMatches(t *testing.T) {
	useTransducer(t, "door bell")
	r := &fakeReplier{}
	WordLink("house bell", r)
	assert.Equal(t, []string{"No words found that can follow 'house'"}, r.replies)
}

func TestWordLink_noRightMatches(t *testing.T) {
	useTransducer(t, "house boat")
	r := &fakeReplier{}
	WordLink("house bell", r)
	assert.Equal(t, []string{"No words found that can precede 'bell'"}, r.replies)
}
