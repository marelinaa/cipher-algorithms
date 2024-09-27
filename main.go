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
	defaultAlphabet = "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ "
)

var (
	alphabet string //todo: string?
	text     string
)

func openFileIfExists(fileName string) (*os.File, error) {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, err
	}

	return os.Open(alphabetFile)
}

func checkAlphabet(alphabet string) (map[rune]int, error) {
	if alphabet == "" {
		return nil, fmt.Errorf("Your alphabet is empty!")
	}

	alphMap := make(map[rune]int)
	n := 0
	for _, char := range alphabet {
		_, ok := alphMap[char] // Check if the key already exists in map
		if ok {
			return nil, fmt.Errorf("Your alphabet has repeated characters!")
		}
		alphMap[char] = n
		n++
	}

	return alphMap, nil
}

func main() {
	// open the alphabet file
	file, err := openFileIfExists(alphabetFile)
	if err != nil {
		log.Fatalf(errorString, err)
	}
	defer file.Close()

	// read the first line from the file using a scanner
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	alphabet = scanner.Text()
	if err := scanner.Err(); err != nil {
		log.Fatalf("error: %v", err)
	}
	_, err = checkAlphabet(alphabet)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("Your alphabet: %s It's power: %d", alphabet, utf8.RuneCountInString(alphabet))
}
