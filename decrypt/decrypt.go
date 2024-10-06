package decrypt

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/marelinaa/cipher-algorithms/encrypt"
	"github.com/marelinaa/cipher-algorithms/keys"
)

// Caesar осуществляет дешифрование Цезаря
func Caesar(input string, key int, alphabetMap map[rune]int, power int) string {
	return encrypt.Caesar(input, -key, alphabetMap, power)
}

func modInverse(a, m int) int {
	for x := 1; x < m; x++ {
		if (a*x)%m == 0 {
			return x
		}
	}
	// Return -1 if no modular inverse is found, though this would ideally not happen if K1 and power are coprime.
	return -1
}

func Affine(input string, key keys.Affine, alphabetMap map[rune]int, power int) string {
	reverseAlphabetMap := make(map[int]rune)

	// Creating a map: numeric representation -> alphabet character
	for char, index := range alphabetMap {
		reverseAlphabetMap[index] = char
	}

	// Find modular inverse of key.K1
	k1Inverse := modInverse(key.K1, power)
	if k1Inverse == -1 {
		panic("K1 has no modular inverse, decryption is not possible!")
	}

	var decryptedText []rune

	for _, char := range input {
		idx := alphabetMap[char]
		// Decryption formula: P = K1^{-1} * (C - K2) mod power
		newIdx := (k1Inverse * (idx - key.K2 + power)) % power
		decryptedText = append(decryptedText, reverseAlphabetMap[newIdx])
	}

	return string(decryptedText)
}

func Substitution(input string, key []rune, alphabetMap map[rune]int, power int) string {
	reverseAlphabetMap := make(map[int]rune)

	// Creating a map: numeric representation -> alphabet character
	for char, index := range alphabetMap {
		reverseAlphabetMap[index] = char
	}

	keyMap := make(map[rune]int)
	for i, char := range key {
		keyMap[char] = i
	}

	var decryptedText []rune

	for _, char := range input {
		// Check if the character is valid in the reverse key map
		i := keyMap[char]
		decryptedChar, ok := reverseAlphabetMap[i]
		if !ok {
			fmt.Printf("input contains invalid character: '%c'\n", char)
			return ""
		}

		// Decrypt by replacing the character with the one at the index from the reverse key
		decryptedText = append(decryptedText, decryptedChar)
	}

	return string(decryptedText)
}

// Функция расшифрования методом перестановки с паролем
func Permutation(input, keyword string) string {
	cols := utf8.RuneCountInString(keyword)
	rows := utf8.RuneCountInString(input) / cols

	if utf8.RuneCountInString(input)%cols != 0 {
		fmt.Println("the length of the ciphertext must be a multiple of the length of the key")

		return ""
	}

	// Заполняем таблицу зашифрованного текста
	sortedTable := make([][]rune, rows)
	for i := range sortedTable {
		sortedTable[i] = make([]rune, cols)
	}

	// Добавляем буквы зашифрованного текста в таблицу
	for i, r := range input {
		row := i / cols
		col := i % cols
		sortedTable[row][col] = r
	}

	// Добавляем буквы пароля в таблицу
	keywordRunes := []rune(keyword)
	colOrder := make([]int, len(keywordRunes))
	for i := range keywordRunes {
		colOrder[i] = i
	}

	// Сортируем индексы столбцов по алфавиту
	sort.Slice(colOrder, func(i, j int) bool {
		return keywordRunes[colOrder[i]] < keywordRunes[colOrder[j]]
	})

	// Создаем обратный порядок для расшифровки
	inverseColOrder := make([]int, len(colOrder))
	for i, col := range colOrder {
		inverseColOrder[col] = i
	}

	// Восстанавливаем оригинальную таблицу
	originalTable := make([][]rune, rows)
	for i := range originalTable {
		originalTable[i] = make([]rune, cols)
		for j, col := range inverseColOrder {
			originalTable[i][j] = sortedTable[i][col]
		}
	}

	// Извлекаем исходный текст, проходя по строкам
	var originalText strings.Builder
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if originalTable[i][j] != 0 {
				originalText.WriteRune(originalTable[i][j])
			}
		}
	}

	return originalText.String()
}

func Vigenere(input string, key string, alphabetMap map[rune]int, power int) string {
	var decrypted strings.Builder

	keyIndex := 0
	for _, char := range input {
		if _, ok := alphabetMap[char]; !ok {
			decrypted.WriteString(string(char))
			continue
		}

		keyLen := utf8.RuneCountInString(key)
		keyChar := rune(key[keyIndex%keyLen])
		shift := (alphabetMap[char] - alphabetMap[keyChar] + power) % power

		for letter, index := range alphabetMap {
			if index == shift {
				decrypted.WriteString(string(letter))
				break
			}
		}

		keyIndex++
	}

	return decrypted.String()
}
