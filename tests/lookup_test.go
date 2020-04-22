package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	emoji "github.com/tmdvs/Go-Emoji-Utils"
)

func BenchmarkLookupByShortcode(b *testing.B) {
	codes := []string{
		":copyright:", ":mahjong:", ":black_joker:", ":flag_antigua__barbuda:", ":flag_australia:",
		":foggy:", ":waning_gibbous_moon:", ":tornado:",
	}

	for _, code := range codes {
		_, err := emoji.LookupEmojiByCode(code)
		require.NoError(b, err)
	}
}
