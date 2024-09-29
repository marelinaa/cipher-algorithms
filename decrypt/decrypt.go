package decrypt

import (
	"github.com/marelinaa/cipher-algorithms/encrypt"
)

// CaesarDecrypt осуществляет дешифрование Цезаря
func CaesarDecrypt(input string, key int, alphabetMap map[rune]int, power int) string {
	// Для дешифрования мы просто применяем обратный сдвиг
	return encrypt.CaesarEncrypt(input, -key, alphabetMap, power)
}
