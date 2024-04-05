package murmur_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/styledjsx/internal/murmur"
)

func TestHash32(t *testing.T) {
	is := is.New(t)
	is.Equal(murmur.Hash32("test"), uint32(3127628307))
	is.Equal(murmur.Hash32("Hello, world!"), uint32(3224780355))
	is.Equal(murmur.Hash32("The quick brown fox jumps over the lazy dog"), uint32(776992547))
}

func TestHash(t *testing.T) {
	is := is.New(t)
	is.Equal(murmur.Hash("test"), "mvnku")
	is.Equal(murmur.Hash("Hello, world!"), "1elBxm")
	is.Equal(murmur.Hash("The quick brown fox jumps over the lazy dog"), "EPQyW")
	is.Equal(murmur.Hash("The quick brown fox jumps over the lazy dogs"), "3B9DNK")
}
