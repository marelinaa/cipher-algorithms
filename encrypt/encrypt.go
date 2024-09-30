package encrypt

import (
	"fmt"
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

func Permutation(input string, key []int, blockLen int) string {
	var result []rune

	// Добавляем паддинг, если длина текста не кратна длине ключа
	if len(input)%blockLen != 0 {
		paddingLen := blockLen - (len(input) % blockLen)
		paddingChar := rune('X') // Символ паддинга
		for i := 0; i < paddingLen; i++ {
			input += string(paddingChar)
		}
	}

	// Разбиваем текст на блоки длины ключа
	for i := 0; i < len(input); i += blockLen {
		block := []rune(input[i : i+blockLen])
		encryptedBlock := make([]rune, blockLen)

		// Перестановка символов в блоке по порядку ключа
		for j, pos := range key {
			// Уменьшаем на 1, если ключ начинается с 1
			encryptedBlock[j] = block[pos-1]
		}

		result = append(result, encryptedBlock...)
	}

	return string(result)
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
