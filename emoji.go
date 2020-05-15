package emoji

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/desmos-labs/Go-Emoji-Utils/utils"
)

// Emoji - Struct representing Emoji
type Emoji struct {
	Key        string   `json:"key"`
	Value      string   `json:"value"`
	Descriptor string   `json:"descriptor"`
	Shortcodes []string `json:"shortcodes"`
}

func (emoji Emoji) Equals(other Emoji) bool {
	if len(emoji.Shortcodes) != len(other.Shortcodes) {
		return false
	}

	for index, code := range emoji.Shortcodes {
		if code != other.Shortcodes[index] {
			return false
		}
	}

	return emoji.Key == other.Key && emoji.Value == other.Value && emoji.Descriptor == other.Descriptor
}

// Unmarshal the emoji JSON into the Emojis map
func init() {
	// Work out where we are in relation to the caller
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	// Open the Emoji definition JSON and Unmarshal into map
	jsonFile, err := os.Open(path.Dir(filename) + "/data/emoji.json")
	if jsonFile != nil {
		defer jsonFile.Close()
	}
	if err != nil && len(Emojis) < 1 {
		fmt.Println(err)
	}

	byteValue, e := ioutil.ReadAll(jsonFile)
	if e != nil {
		if len(Emojis) > 0 { // Use build-in emojis data (from emojidata.go)
			populateEmojisList()
			return
		}
		panic(e)
	}

	err = json.Unmarshal(byteValue, &Emojis)
	if err != nil {
		panic(e)
	}
	populateEmojisList()
}

// populateEmojisList - takes the emojis stored inside the Emojis map and puts them into the EmojisList as well
func populateEmojisList() {
	for _, emoji := range Emojis {
		EmojisList = append(EmojisList, emoji)
	}
}

// LookupEmoji - Lookup a single emoji definition
func LookupEmoji(emojiString string) (emoji Emoji, err error) {

	hexKey := utils.StringToHexKey(emojiString)

	// If we have a definition for this string we'll return it,
	// else we'll return an error
	if e, ok := Emojis[hexKey]; ok {
		emoji = e
	} else {
		err = fmt.Errorf("No record for \"%s\" could be found", emojiString)
	}

	return emoji, err
}

// LookupEmojis - Lookup definitions for each emoji in the input
func LookupEmojis(emoji []string) (matches []interface{}) {
	for _, emoji := range emoji {
		if match, err := LookupEmoji(emoji); err == nil {
			matches = append(matches, match)
		} else {
			matches = append(matches, err)
		}
	}

	return
}

// LookupEmojiByCode - Lookup a single emoji definition by its shortcode
func LookupEmojiByCode(shortcode string) (emoji Emoji, err error) {
	for _, emoji := range Emojis {
		for _, s := range emoji.Shortcodes {
			if s == shortcode {
				return emoji, nil
			}
		}
	}

	return Emoji{}, fmt.Errorf("No emoji found for shortcode \"%s\"", shortcode)
}

// RemoveAll - Remove all emoji
func RemoveAll(input string) string {

	// Find all the emojis in this string
	matches := FindAll(input)

	for _, item := range matches {
		emo := item.Match.(Emoji)
		rs := []rune(emo.Value)
		for _, r := range rs {
			input = strings.ReplaceAll(input, string([]rune{r}), "")
		}
	}

	// Remove and trim and left over whitespace
	return strings.TrimSpace(strings.Join(strings.Fields(input), " "))
	//return input
}
