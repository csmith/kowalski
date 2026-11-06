package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"testing"

	"code.dny.dev/ssrf"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// allowedServer stands in for a public host, by making the image client permit the address it listens on.
func allowedServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.2:0")
	require.NoError(t, err)

	server := httptest.NewUnstartedServer(handler)
	server.Listener.Close()
	server.Listener = listener
	server.Start()
	t.Cleanup(server.Close)

	original := imageClient
	imageClient = newImageClient(ssrf.New(
		ssrf.WithAnyPort(),
		ssrf.WithAllowedV4Prefixes(netip.MustParsePrefix("127.0.0.2/32")),
	))
	t.Cleanup(func() { imageClient = original })

	return server
}

func restrictedServer(t *testing.T) (*httptest.Server, *bool) {
	t.Helper()

	hit := false
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		hit = true
	}))
	t.Cleanup(server.Close)
	return server, &hit
}

func TestFetchImage_refusesRestrictedAddresses(t *testing.T) {
	server, hit := restrictedServer(t)
	port := server.URL[strings.LastIndex(server.URL, ":"):]

	for _, url := range []string{
		server.URL,
		"http://localhost" + port,
		"http://[::1]" + port,
		"http://10.0.0.1/",
		"http://192.168.0.1/",
		"http://169.254.169.254/latest/meta-data/",
		"http://[fd00::1]/",
	} {
		t.Run(url, func(t *testing.T) {
			r := &fakeReplier{}
			_, ok := fetchImage(url, r)
			assert.False(t, ok)
			assert.Equal(t, []string{"Unable to download image"}, r.replies)
		})
	}
	assert.False(t, *hit)
}

func TestFetchImage_refusesRedirectsToRestrictedAddresses(t *testing.T) {
	target, hit := restrictedServer(t)
	redirector := allowedServer(t, http.RedirectHandler(target.URL, http.StatusFound))

	r := &fakeReplier{}
	_, ok := fetchImage(redirector.URL, r)
	assert.False(t, ok)
	assert.False(t, *hit)
}

func TestFetchImage_limitsDownloadSize(t *testing.T) {
	server := allowedServer(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(make([]byte, maxImageSize+1024))
	}))

	image, ok := fetchImage(server.URL, &fakeReplier{})
	require.True(t, ok)
	defer image.Close()

	data, err := io.ReadAll(image)
	require.NoError(t, err)
	assert.Len(t, data, maxImageSize)
}

func TestFileCommands_doNotEchoDownloadErrors(t *testing.T) {
	server, _ := restrictedServer(t)

	for name, command := range map[string]FileCommand{"colours": Colours, "hidden": HiddenPixels, "rgb": RGB} {
		t.Run(name, func(t *testing.T) {
			r := &fakeReplier{}
			command("", []string{server.URL}, r)
			assert.Equal(t, []string{"Unable to download image"}, r.replies)
		})
	}
}
