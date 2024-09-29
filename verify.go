package main

import (
	"fmt"
	"unicode"
)

// verifyAlphabet checks the alphabet for accuracy.
// If alphabet is correct, creates map with alphabet characters as keys and numerical representations as values.
func verifyAlphabet(alphabet string) (map[rune]int, error) {
	if alphabet == "" {
		return nil, fmt.Errorf("alphabet can not be empty")
	}

	alphabetMap := make(map[rune]int)
	n := 0
	for _, char := range alphabet {
		_, ok := alphabetMap[char] // check if the key already exists in map
		if ok {
			return nil, fmt.Errorf("your alphabet has repeated characters")
		}
		alphabetMap[char] = n
		n++
	}

	return alphabetMap, nil
}

// verifyText checks the text for accuracy
func verifyText(text string, alphabetMap map[rune]int) error {
	if text == "" {
		return fmt.Errorf("text can not be empty")
	}

	for _, char := range text {
		_, ok := alphabetMap[char] // check if the text contains allowable characters
		if !ok {
			return fmt.Errorf("text contains characters not from the alphabet")
		}
	}

	return nil
}

func verifyCaesarKey(key string) (int, error) {
	keyRunes := []rune(key)
	if len(keyRunes) != 1 {
		return -1, fmt.Errorf("caesar key must contain one symbol from the alphabet")
	}

	k, ok := alphabetMap[keyRunes[0]]
	if !ok {
		return -1, fmt.Errorf("caesar key must contain one symbol from the alphabet")
	}

	return k, nil
}

func isControl(r rune) bool {
	cs := []rune{'\a', '\b', '\f', '\n', '\r', '\t', '\v'}

	for _, c := range cs {
		if r == c {
			return true
		}
	}

	return false
}

func containsControlCharacter(s string) bool {
	for _, r := range s {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}
