package tests

import (
	"testing"

	emoji "github.com/desmos-labs/Go-Emoji-Utils"
	"github.com/stretchr/testify/require"
)

func TestListEmojis(t *testing.T) {
	require.NotEmpty(t, emoji.EmojisList)
}
