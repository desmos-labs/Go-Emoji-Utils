package main

// Import emoji data from Emojipedia.org
// Useful for rebuilding the emoji data found in the `data/emoji.json` file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/tmdvs/Go-Emoji-Utils/utils"

	"github.com/tmdvs/Go-Emoji-Utils"

	"github.com/PuerkitoBio/goquery"
)

var (
	// emojiCodeRegex represents the regex that allows to strip down an emoji descriptor to get its shortcode
	emojiCodeRegex = regexp.MustCompile("[^a-zA-Z0-9\\s-]+")
)

type lookup struct {
	Name string
	URL  string
}

func main() {
	fmt.Println("Updating Emoji Definition using Emojipedia…")

	// Grab the latest Apple Emoji Definitions
	res, err := http.Get("https://emojipedia.org/apple/")
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	// Load the Apple Emoji HTML page into goquery so that we
	// can query the DOM
	doc, docErr := goquery.NewDocumentFromReader(res.Body)
	if docErr != nil {
		panic(docErr)
	}

	// Create a channel for lookups so that we can do this async
	lookups := make(chan lookup)

	go func() {
		// Find all emojis on the page
		doc.Find("ul.emoji-grid li").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			emojiPage, _ := s.Find("a").Attr("href")
			title, _ := s.Find("img").Attr("title")

			fmt.Printf("Adding Emoji %d to lookups: %s - %s\n", i, title, emojiPage)

			// Add this specific emoji to the lookups to complete
			lookups <- lookup{
				Name: title,
				URL:  "https://emojipedia.org" + emojiPage,
			}
		})

		close(lookups)
	}()

	emojis := map[string]emoji.Emoji{}

	// Process a lookup
	for lookup := range lookups {
		fmt.Println("Looking up " + lookup.Name)

		// Grab the emojipedia page for this emoji
		res, err := http.Get(lookup.URL)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Create a new goquery reader
		doc, docErr := goquery.NewDocumentFromReader(res.Body)
		if docErr != nil {
			panic(docErr)
		}

		// Grab the emoji from the "Copy emoji" input field on the HTML page
		emojiString, _ := doc.Find(".copy-paste input[type=text]").Attr("value")

		// Grab the shortcodes from the list
		var shortCodes []string
		doc.Find("ul.shortcodes li").Each(func(i int, s *goquery.Selection) {
			shortCodes = append(shortCodes, s.Text())
		})

		// Generate a shortcode based on the emoji name if no shortcodes were found
		if len(shortCodes) == 0 {
			shortCodes = append(shortCodes, getEmojiShortCode(lookup.Name))
		}

		// Convert the raw Emoji value to our hex key
		hexString := utils.StringToHexKey(emojiString)

		// Add this emoji to our map
		emojis[hexString] = emoji.Emoji{
			Key:        hexString,
			Value:      emojiString,
			Descriptor: lookup.Name,
			Shortcodes: shortCodes,
		}

		// Print our progress to the console
		fmt.Println(emojis[hexString])
	}

	// Marshal the emojis map as JSON and write to the data directory
	s, _ := json.MarshalIndent(emojis, "", "\t")
	ioutil.WriteFile("data/emoji.json", []byte(s), 0644)
}

// getEmojiShortCode takes the given string value and converts it to the appropriate shortcode.
// In particular, it will return a lowercase value with all the special symbols removed and
// the spaces converted to be "_".
func getEmojiShortCode(name string) string {
	shortCode := emojiCodeRegex.ReplaceAllString(name, "")
	shortCode = strings.ToLower(shortCode)
	shortCode = strings.ReplaceAll(shortCode, " ", "_")
	shortCode = strings.ReplaceAll(shortCode, "-", "_")
	shortCode = fmt.Sprintf(":%s:", shortCode)
	return shortCode
}
