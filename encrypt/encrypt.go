package encrypt

import "github.com/marelinaa/cipher-algorithms/keys"

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
