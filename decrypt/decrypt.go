package decrypt

import (
	"github.com/marelinaa/cipher-algorithms/encrypt"
	"github.com/marelinaa/cipher-algorithms/keys"
)

// Caesar осуществляет дешифрование Цезаря
func Caesar(input string, key int, alphabetMap map[rune]int, power int) string {
	// Для дешифрования применяем обратный сдвиг
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
