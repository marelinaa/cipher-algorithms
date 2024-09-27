package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode/utf8"
)

const (
	errorString     = "error: %v"
	alphabetFile    = "alphabet.txt"
	textFile        = "text.txt"
	defaultAlphabet = "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ "
)

var (
	text string
)

// openAndExtractText opens the file with. Returns the first line from that file
func openAndExtractText(fileName string) (string, error) {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return "", err
	}

	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// read the first line from the file using a scanner
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	text = scanner.Text()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return text, nil
}

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
func verifyText(text string, alphMap map[rune]int) error {
	if text == "" {
		return fmt.Errorf("text can not be empty")
	}

	for _, char := range text {
		_, ok := alphMap[char] // check if the text contains allowable characters
		if !ok {
			return fmt.Errorf("text contains characters not from the alphabet")
		}
	}

	return nil
}

func main() {
	// open the alphabet file
	alphabet, err := openAndExtractText(alphabetFile) //todo: string? ask Ilya! (i had global var, but it was shadowed by this alphabet)
	if err != nil {
		fmt.Printf("error: %v => now using default alphabet\n", err)
		alphabet = defaultAlphabet
	}

	// check the alphabet for accuracy and make map of alphabet characters
	alphabetMap, err := verifyAlphabet(alphabet)
	if err != nil {
		log.Fatalf(errorString, err)
	}

	power := utf8.RuneCountInString(alphabet) // power of the alphabet
	fmt.Printf("Your alphabet: %s It's power: %d\n", alphabet, power)

	// open the file with the text
	text, err = openAndExtractText(textFile)
	if err != nil {
		log.Fatalf(errorString, err)
	}

	// check the text for accuracy
	err = verifyText(text, alphabetMap)
	if err != nil {
		log.Fatalf(errorString, err)
	}
}
