package verify

import (
	"fmt"
	"unicode"

	"github.com/marelinaa/cipher-algorithms/keys"
)

// Alphabet checks the alphabet for accuracy.
// If alphabet is correct, creates map with alphabet characters as keys and numerical representations as values.
func Alphabet(alphabet string) (map[rune]int, error) {
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

// Text checks the text for accuracy
func Text(text string, alphabetMap map[rune]int) error {
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

func CaesarKey(key string, alphabetMap map[rune]int) (int, error) {
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

func gcd(a, b int) int {
	if b == 0 {
		return a
	} else {
		return gcd(b, a%b)
	}
}

func areCoprime(a, b int) bool {
	// If the GCD of two numbers is 1, they are coprime
	return gcd(a, b) == 1
}

func AffineKey(key string, alphabetMap map[rune]int) (keys.Affine, error) {
	err := "affine key must be a pair of symbols from the alphabet, witout delimiters"
	keyRunes := []rune(key)
	if len(keyRunes) != 2 {
		return keys.Affine{}, fmt.Errorf(err)
	}

	k1, ok := alphabetMap[keyRunes[0]]
	if !ok {
		return keys.Affine{}, fmt.Errorf(err)
	}
	k2, ok := alphabetMap[keyRunes[1]]
	if !ok {
		return keys.Affine{}, fmt.Errorf(err)
	}

	if !areCoprime(k1, k2) {
		return keys.Affine{}, fmt.Errorf("numeric representations of symbols must be coprime")
	}

	affineKey := keys.Affine{
		K1: k1,
		K2: k2,
	}

	return affineKey, nil
}

func IsControl(r rune) bool {
	cs := []rune{'\a', '\b', '\f', '\n', '\r', '\t', '\v'}

	for _, c := range cs {
		if r == c {
			return true
		}
	}

	return false
}

func ContainsControlCharacter(s string) bool {
	for _, r := range s {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}
