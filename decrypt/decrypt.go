package decrypt

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/marelinaa/cipher-algorithms/encrypt"
	"github.com/marelinaa/cipher-algorithms/keys"
	"golang.org/x/exp/rand"
)

// Caesar осуществляет дешифрование Цезаря
func Caesar(input string, key int, alphabetMap map[rune]int, power int) string {
	return encrypt.Caesar(input, -key, alphabetMap, power)
}

func modInverse(a, m int) (int, error) {
	a = a % m
	for x := 1; x < m; x++ {
		if (a*x)%m == 1 {
			return x, nil
		}
	}
	return 0, errors.New("no modular inverse exists")
}

func Affine(input string, key keys.Affine, alphabetMap map[rune]int, power int) string {
	reverseAlphabetMap := make(map[int]rune)

	// Creating a map: numeric representation -> alphabet character
	for char, index := range alphabetMap {
		reverseAlphabetMap[index] = char
	}

	// Find modular inverse of key.K1
	k1Inverse, err := modInverse(key.K1, power)
	if err != nil {
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

func Permutation(input, keyword string, alphabetMap map[rune]int) string {
	cols := utf8.RuneCountInString(keyword)
	paddingLen := cols - (utf8.RuneCountInString(input) % cols)
	rows := (utf8.RuneCountInString(input) + paddingLen) / cols
	// Проверяем, нужно ли добавить символы для выравнивания
	if utf8.RuneCountInString(input)%cols != 0 {
		log.Println("длина шифртекста не кратна длине ключа")
	}

	// Получаем порядок перестановки
	order := getKeywordOrder(keyword, alphabetMap)
	fmt.Println("Порядок перестановки:", order)

	// Получаем обратный порядок перестановки
	reverseOrder := make([]int, len(order))
	for i, pos := range order {
		reverseOrder[pos] = i
	}
	fmt.Println("Обратный порядок перестановки:", reverseOrder)

	// Создаем таблицу для расшифровки
	table := make([][]rune, rows)
	inputRunes := []rune(input)
	idx := 0
	for i := 0; i < rows; i++ {
		table[i] = make([]rune, cols)
		for j := 0; j < cols; j++ {
			table[i][j] = inputRunes[idx]
			idx++
		}
	}

	// Восстанавливаем исходный порядок в каждой строке
	for i := 0; i < rows; i++ {
		table[i] = rearrangeRow(table[i], reverseOrder)
	}

	// Собираем исходный текст
	var plainText strings.Builder
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if table[i][j] != 'X' { // Игнорируем символы заполнения
				plainText.WriteRune(table[i][j])
			}
		}
	}

	return plainText.String()
}

// Функция для получения порядка перестановки по алфавиту
func getKeywordOrder(keyword string, alphabetMap map[rune]int) []int {
	runes := []rune(keyword)
	n := len(runes)

	type letterIndex struct {
		letter rune
		index  int
	}

	letters := make([]letterIndex, n)
	for i, r := range runes {
		letters[i] = letterIndex{r, i}
	}

	sort.Slice(letters, func(i, j int) bool {
		return alphabetMap[letters[i].letter] < alphabetMap[letters[j].letter]
	})

	result := make([]int, n)
	for sortedIndex, li := range letters {
		result[li.index] = sortedIndex
	}

	return result
}

// Функция для перестановки элементов строки в соответствии с порядком order
func rearrangeRow(row []rune, order []int) []rune {
	rearranged := make([]rune, len(row))
	for i, pos := range order {
		rearranged[pos] = row[i]
	}
	return rearranged
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

func randomRune(alphabetMap map[rune]int, power int) rune {
	var alphabet []rune
	for r, count := range alphabetMap {
		for i := 0; i < count; i++ {
			alphabet = append(alphabet, r)
		}
	}

	i := rand.Intn(power - 1)
	return alphabet[i]
}

func determinant2x2(key [2][2]int) int {
	k11 := key[0][0]
	k12 := key[0][1]
	k21 := key[1][0]
	k22 := key[1][1]

	return (k11*k22 - k12*k21)
}

func inverseMatrix(key [2][2]int, power int) ([2][2]int, error) {
	det := determinant2x2(key)

	if det == 0 {
		return [2][2]int{}, errors.New("matrix is singular, no inverse exists")
	}

	invDet, err := modInverse(det, power) // Модульное обратное определителя, предполагая, что работаем в поле 26
	if err != nil {
		return [2][2]int{}, errors.New("no modular inverse found for determinant")
	}

	a := key[0][0]
	b := key[0][1]
	c := key[1][0]
	d := key[1][1]

	inv := [2][2]int{
		{d * invDet % 26, (-b) * invDet % power},
		{(-c) * invDet % 26, a * invDet % power},
	}

	return inv, nil
}

func Mod(x, y int) int {
	if x < 0 {
		a := -x / y

		return ((a)+1)*y + x
	}

	return x - (x/y)*y
}

func hillEncryptPair(k11, k12, k21, k22, p1, p2, power int) (int, int) {
	c1 := Mod(k11*p1+k21*p2, power)
	c2 := Mod(k12*p1+k22*p2, power)

	return c1, c2
}

func Hill(input string, key [2][2]int, alphabetMap map[rune]int, power int) string {
	key, err := inverseMatrix(key, power)
	if err != nil {
		fmt.Println(err)
	}

	reverseAlphabetMap := make(map[int]rune)
	for char, idx := range alphabetMap {
		reverseAlphabetMap[idx] = char
	}

	input = strings.ToUpper(input)
	if utf8.RuneCountInString(input)%2 != 0 {
		rand := randomRune(alphabetMap, power)
		input += string(rand) // Добавляем символ для выравнивания
	}

	var ciphertext strings.Builder
	text := []rune(input)
	for i := 0; i < utf8.RuneCountInString(input); i += 2 {
		p1 := alphabetMap[text[i]]
		p2 := alphabetMap[text[i+1]]

		c1, c2 := hillEncryptPair(key[0][0], key[0][1], key[1][0], key[1][1], p1, p2, power)
		ciphertext.WriteRune(reverseAlphabetMap[c1])
		ciphertext.WriteRune(reverseAlphabetMap[c2])
	}

	return ciphertext.String()
}
