package main

import (
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

type fakeReplier struct {
	replies []string
}

func (f *fakeReplier) reply(format string, a ...interface{}) {
	f.replies = append(f.replies, fmt.Sprintf(format, a...))
}

func (f *fakeReplier) replyWithFiles(_ []*discordgo.File, format string, a ...interface{}) {
	f.reply(format, a...)
}

func TestChunk(t *testing.T) {
	r := &fakeReplier{}
	Chunk("2 1 abcde", r)
	assert.Equal(t, []string{"Chunked: ab c de"}, r.replies)
}

func TestChunk_rejectsInvalidSizes(t *testing.T) {
	for _, input := range []string{"0 a", "-1 a", "3 0 abcdef", "3 -2 abcdef"} {
		t.Run(input, func(t *testing.T) {
			r := &fakeReplier{}
			Chunk(input, r)
			assert.Equal(t, []string{"Chunk sizes must be greater than zero"}, r.replies)
		})
	}
}
