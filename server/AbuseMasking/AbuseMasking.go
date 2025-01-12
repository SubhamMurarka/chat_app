package AbuseMasking

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/SubhamMurarka/chat_app/server/config"
)

var trie *Trie

func Loadfile() []string {
	config.ConnectS3()

	words := strings.Fields(config.FileContent)

	return words
}

func MakeTrie(words []string) {
	trie = NewTrie()
	for _, val := range words {
		trie.insert(val)
	}
}

func Filter(line string) string {
	words := strings.Fields(line) // Split the message into words
	for i, word := range words {
		cleanedWord := cleanWord(word) // Remove any non-alphabetic characters forP Trie search
		fmt.Println("cleanedword: ", cleanedWord)
		if trie.search(cleanedWord) {
			// Calculate the starting index of the cleaned word within the original word
			startIndex := strings.Index(strings.ToLower(word), cleanedWord)
			if startIndex != -1 {
				// Mask only the cleaned word while keeping non-alphabetic characters intact
				masked := strings.Repeat("*", len(cleanedWord))
				words[i] = word[:startIndex] + masked + word[startIndex+len(cleanedWord):]
				fmt.Println("word: ", words[i])
			}
		}
	}
	return strings.Join(words, " ")
}

// Helper function to clean a copy word by removing non-alphabetic characters
func cleanWord(word string) string {
	var result []rune
	for _, ch := range word {
		if unicode.IsLetter(ch) {
			result = append(result, unicode.ToLower(ch)) // Convert to lowercase
		}
	}
	return string(result)
}
