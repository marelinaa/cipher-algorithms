package encrypt

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/marelinaa/cipher-algorithms/keys"
	"golang.org/x/exp/rand"
)

func Caesar(input string, key int, alphabetMap map[rune]int, power int) string {
	reverseAlphabetMap := make(map[int]rune)

	for char, index := range alphabetMap {
		reverseAlphabetMap[index] = char
	}

	var encryptedText []rune

	for _, char := range input {
		idx := alphabetMap[char]
		newIdx := (idx + key) % power
		if newIdx < 0 {
			newIdx += power
		}
		encryptedText = append(encryptedText, reverseAlphabetMap[newIdx])
	}

	return string(encryptedText)
}

func Affine(input string, key keys.Affine, alphabetMap map[rune]int, power int) string {
	reverseAlphabetMap := make(map[int]rune)

	for char, index := range alphabetMap {
		reverseAlphabetMap[index] = char
	}

	var encryptedText []rune

	for _, char := range input {
		idx := alphabetMap[char]
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

func Permutation(input, keyword string, alphabetMap map[rune]int, power int) string {
	cols := utf8.RuneCountInString(keyword)
	paddingLen := cols - (utf8.RuneCountInString(input) % cols)
	rows := utf8.RuneCountInString(input) / cols
	if utf8.RuneCountInString(input)%cols != 0 {
		rows = (utf8.RuneCountInString(input) + paddingLen) / cols
	}

	fmt.Println(rows, cols)

	if utf8.RuneCountInString(input)%cols != 0 {
		fmt.Println("я тут")
		paddingChar := randomRune(alphabetMap, power)
		for i := 0; i < paddingLen; i++ {
			input += string(paddingChar)
		}
	}

	runes := []rune(keyword)
	n := utf8.RuneCountInString(keyword)

	// Структура для хранения букв и их исходных индексов
	type letterIndex struct {
		letter rune
		index  int
	}

	// Заполняем структуру буквы и их индексы
	letters := make([]letterIndex, n)
	for i, r := range runes {
		letters[i] = letterIndex{r, i}
	}

	// Сортируем структуру по алфавиту
	sort.Slice(letters, func(i, j int) bool {
		return alphabetMap[letters[i].letter] < alphabetMap[letters[j].letter]
	})

	// Создаем слайс результата
	order := make([]int, n)
	for sortedIndex, li := range letters {
		order[li.index] = sortedIndex
	}

	fmt.Println(order)

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

	// Переставляем элементы каждой строки в соответствии с порядком из слайса order
	for i := 0; i < rows; i++ {
		table[i] = rearrangeRow(table[i], order)
	}

	// Создаем шифротекст, проходя по строкам
	var cipherText strings.Builder
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if table[i][j] != 0 { // Игнорируем пустые ячейки
				cipherText.WriteRune(table[i][j])
			}
		}
	}

	return cipherText.String()
}

func rearrangeRow(row []rune, order []int) []rune {
	rearranged := make([]rune, len(row))
	for i, pos := range order {
		rearranged[pos] = row[i]
	}
	return rearranged
}

func Vigenere(plaintext string, key string, alphabet map[rune]int) string {
	alphabetLength := len(alphabet)
	reverseMap := make(map[int]rune)
	for char, index := range alphabet {
		reverseMap[index] = char
	}

	keyLength := len(key)
	keyIndices := make([]int, keyLength)
	for i, char := range key {
		keyIndices[i] = alphabet[char]
	}

	encryptedText := make([]rune, 0)

	for i, char := range plaintext {
		p := alphabet[char]
		k := keyIndices[i%keyLength]
		encryptedCharIndex := (p + k) % alphabetLength
		encryptedText = append(encryptedText, reverseMap[encryptedCharIndex])
	}

	return string(encryptedText)
}

// Функция для нахождения обратного числа по модулю
func modInverse(a, m int) (int, error) {
	a = a % m
	for x := 1; x < m; x++ {
		if (a*x)%m == 1 {
			return x, nil
		}
	}
	return 0, errors.New("no modular inverse exists")
}

func Mod(x, y int) int {
	if x < 0 {
		a := -x / y
		fmt.Println(a)

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
