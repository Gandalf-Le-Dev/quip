package nanoid_test

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/Gandalf-Le-Dev/quip/internal/pkg/utils/nanoid"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	t.Run("short alphabet", func(t *testing.T) {
		alphabet := ""
		_, err := nanoid.Generate(alphabet, 32)
		assert.Error(t, err, "should return error if the alphabet is too small")
	})

	t.Run("long alphabet", func(t *testing.T) {
		alphabet := strings.Repeat("a", 256)
		_, err := nanoid.Generate(alphabet, 32)
		assert.Error(t, err, "should return error if the alphabet is too long")
	})

	t.Run("negative ID length", func(t *testing.T) {
		_, err := nanoid.Generate("abcdef", -1)
		assert.Error(t, err, "should return error if the requested ID length is invalid")
	})

	t.Run("happy path", func(t *testing.T) {
		alphabet := "abcdef"
		id, err := nanoid.Generate(alphabet, 6)
		assert.NoError(t, err, "shouldn't return error")
		assert.Len(t, id, 6, "should return ID of requested length")
		for _, r := range id {
			assert.True(t, strings.ContainsRune(alphabet, r), "should use given alphabet")
		}
	})

	t.Run("works with unicode", func(t *testing.T) {
		alphabet := "🚀💩🦄🤖"
		id, err := nanoid.Generate(alphabet, 6)
		assert.NoError(t, err, "shouldn't return error")
		assert.Equal(t, utf8.RuneCountInString(id), 6, "should return ID of requested length")
		for _, r := range id {
			assert.True(t, strings.ContainsRune(alphabet, r), "should use given alphabet")
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("negative ID length", func(t *testing.T) {
		_, err := nanoid.New(-1)
		assert.Error(t, err, "should return error if the requested ID length is invalid")
	})

	t.Run("happy path", func(t *testing.T) {
		id, err := nanoid.New()
		assert.NoError(t, err, "shouldn't return error")
		assert.Len(t, id, 21, "should return ID of default length")
	})

	t.Run("custom length", func(t *testing.T) {
		id, err := nanoid.New(6)
		assert.NoError(t, err, "shouldn't return error")
		assert.Len(t, id, 6, "should return ID of requested length")
	})
}