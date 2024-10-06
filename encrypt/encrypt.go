package encrypt

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/marelinaa/cipher-algorithms/keys"
)

// Caesar осуществляет шифрование Цезаря с использованием символа-ключа
func Caesar(input string, key int, alphabetMap map[rune]int, power int) string {
	reverseAlphabetMap := make(map[int]rune)

	// Создаем мапу: числовое представление - символ алфавита
	for char, index := range alphabetMap {
		reverseAlphabetMap[index] = char
	}

	var encryptedText []rune

	for _, char := range input {
		idx := alphabetMap[char]
		// Сдвиг символа с учетом алфавита
		newIdx := (idx + key) % power
		if newIdx < 0 {
			newIdx += power // Обрабатываем отрицательные индексы
		}
		encryptedText = append(encryptedText, reverseAlphabetMap[newIdx])
	}

	return string(encryptedText)
}

func Affine(input string, key keys.Affine, alphabetMap map[rune]int, power int) string {
	reverseAlphabetMap := make(map[int]rune)

	// Создаем мапу: числовое представление - символ алфавита
	for char, index := range alphabetMap {
		reverseAlphabetMap[index] = char
	}

	var encryptedText []rune

	for _, char := range input {
		idx := alphabetMap[char]
		// Сдвиг символа с учетом алфавита
		newIdx := (key.K1*idx + key.K2) % power
		encryptedText = append(encryptedText, reverseAlphabetMap[newIdx])
	}

	return string(encryptedText)
}

func Substitution(input string, key []int, alphabetMap map[rune]int, power int) string {
	// Create reverse map to find character by its numeric index
	reverseAlphabetMap := make(map[int]rune)
	for char, idx := range alphabetMap {
		reverseAlphabetMap[idx] = char
	}

	var encryptedText []rune

	for _, char := range input {
		// Find the index of the character in the alphabet
		idx, ok := alphabetMap[char]
		if !ok {
			fmt.Printf("input contains invalid character: '%c'\n", char)
			return ""
		}

		// Encrypt by replacing the character with the one at the index from the key
		encryptedChar := reverseAlphabetMap[key[idx]]
		encryptedText = append(encryptedText, encryptedChar)
	}

	return string(encryptedText)
}

// func Permutation(input string, keyString string, key []int, alphabetMap map[rune]int, blockLen int) string {
// 	var result []rune
// 	reverseAlphabetMap := make(map[int]rune)

// 	for char, index := range alphabetMap {
// 		reverseAlphabetMap[index] = char
// 	}

// 	// если длина текста не кратна длине ключа
// 	if utf8.RuneCountInString(input)%blockLen != 0 {
// 		paddingLen := blockLen - (len(input) % blockLen)
// 		paddingChar := rune('X')
// 		for i := 0; i < paddingLen; i++ {
// 			input += string(paddingChar)
// 		}
// 	}

// 	sort.Ints(key)

// 	for n := range key {
// 		char := reverseAlphabetMap[n]
// 		ind :=
// 	}

// 	return string(result)
// }

// Функция для шифрования методом перестановки с паролем
func Permutation(input, keyword string) string {
	cols := utf8.RuneCountInString(keyword)
	paddingLen := cols - (utf8.RuneCountInString(input) % cols)
	rows := (utf8.RuneCountInString(input) + paddingLen) / cols

	if utf8.RuneCountInString(input)%cols != 0 {
		paddingLen := cols - (utf8.RuneCountInString(input) % cols)
		paddingChar := rune('X') //todo:change
		for i := 0; i < paddingLen; i++ {
			input += string(paddingChar)
		}
	}

	// Заполняем таблицу текста
	table := make([][]rune, rows)
	for i := range table {
		table[i] = make([]rune, cols)
	}

	for i, r := range input {
		row := i / cols
		col := i % cols
		table[row][col] = r
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

	// Переставляем столбцы в таблице
	sortedTable := make([][]rune, rows)
	for i := range sortedTable {
		sortedTable[i] = make([]rune, cols)
		for j, col := range colOrder {
			sortedTable[i][j] = table[i][col]
		}
	}

	// Создаем шифротекст, проходя по строкам
	var cipherText strings.Builder
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if sortedTable[i][j] != 0 { // Игнорируем пустые ячейки
				cipherText.WriteRune(sortedTable[i][j])
			}
		}
	}

	return cipherText.String()
}

func Vigenere(input string, key string, alphabetMap map[rune]int, power int) string {
	var encrypted strings.Builder

	keyIndex := 0
	for _, char := range input {
		if _, ok := alphabetMap[char]; !ok {
			encrypted.WriteString(string(char))
			continue
		}

		keyLen := utf8.RuneCountInString(key)
		keyChar := rune(key[keyIndex%keyLen])
		shift := (alphabetMap[keyChar] + alphabetMap[char]) % power

		for letter, index := range alphabetMap {
			if index == shift {
				encrypted.WriteString(string(letter))
				break
			}
		}

		keyIndex++
	}

	return encrypted.String()
}
