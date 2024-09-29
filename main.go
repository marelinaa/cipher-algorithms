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
	defaultAlphabet = "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ "
	alphabetFile    = "alphabet.txt"
	textFile        = "in.txt"
	keyFile         = "key.txt"
)

var (
	alphabetMap map[rune]int
	power       int
	text        string
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
	i := 1
	for scanner.Scan() {
		if i != 1 {
			return "", fmt.Errorf("file %s has more than one row", fileName)
		}
		text = scanner.Text()
		if err := scanner.Err(); err != nil {
			return "", err
		}
		i++
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
		return -1, fmt.Errorf("key must contain one symbol from the alphabet")
	}

	k, ok := alphabetMap[keyRunes[0]]
	if !ok {
		return -1, fmt.Errorf("key must contain one symbol from the alphabet")
	}

	return k, nil
}

func main() {
	// open the alphabet file
	alphabet, err := openAndExtractText(alphabetFile) //todo: string? ask Ilya! (i had global var, but it was shadowed by this alphabet)
	if err != nil {
		fmt.Printf("error: %v => now using default alphabet\n", err)
		alphabet = defaultAlphabet
	}

	// check the alphabet for accuracy and make map of alphabet characters
	alphabetMap, err = verifyAlphabet(alphabet) //todo:global var or pass to func
	if err != nil {
		log.Fatalf(errorString, err)
	}

	power = utf8.RuneCountInString(alphabet) // power of the alphabet
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

	//todo: implement choosing ciphering method

	// open the file with the key
	keyString, err := openAndExtractText(keyFile)
	if err != nil {
		log.Fatalf(errorString, err)
	}

	key, err := verifyCaesarKey(keyString)
	if err != nil {
		log.Fatalf(errorString, err)
	}
	fmt.Println(key)

}
