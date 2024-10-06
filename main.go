package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	"github.com/marelinaa/cipher-algorithms/decrypt"
	"github.com/marelinaa/cipher-algorithms/encrypt"
	"github.com/marelinaa/cipher-algorithms/verify"
)

const (
	errorString     = "error: %v"
	defaultAlphabet = "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ "
	alphabetFile    = "alphabet.txt"
	textFile        = "in.txt"
	keyFile         = "key.txt"
	encryptFile     = "encrypt.txt"
	decryptFile     = "decrypt.txt"
)

var (
	alphabetMap map[rune]int
	power       int
	input       string
	result      string
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
	var text string
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

func WriteToFile(fileName, text string) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()

	file.WriteString(text)
}

func initializeData() (string, error) {
	// open the alphabet file
	alphabet, err := openAndExtractText(alphabetFile)
	if err != nil {
		fmt.Printf("error: %v => now using default alphabet\n", err)
		alphabet = defaultAlphabet
	}

	// check the alphabet for accuracy and make map of alphabet characters
	alphabetMap, err = verify.Alphabet(alphabet)
	if err != nil {
		return "", err
	}

	power = utf8.RuneCountInString(alphabet) // мощность алфавита
	fmt.Printf("Your alphabet: %s It's power: %d\n", alphabet, power)

	// open the file with the text
	input, err = openAndExtractText(textFile)
	if err != nil {
		return "", err
	}

	// check the text for accuracy
	err = verify.Text(input, alphabetMap)
	if err != nil {
		return "", err
	}

	// open the file with the key
	keyString, err := openAndExtractText(keyFile)
	if err != nil {
		return "", err
	}

	return keyString, nil
}

func main() {

	keyString, err := initializeData()
	if err != nil {
		log.Fatalf("error during initialization: %v", err)
	}

	// main logic
	for {
		// choosing a cryptosystem
		fmt.Println("---------------------------")
		fmt.Println("Choose the cryptographic system:")
		fmt.Println("1: Caesar cipher")
		fmt.Println("2: Affine cipher")
		fmt.Println("3: Simple substitution cipher")
		fmt.Println("4: Hill cipher")
		fmt.Println("5: Permutation cipher")
		fmt.Println("6: Vigenere cipher")
		fmt.Println("0: Exit")

		var cipherChoice int
		for {
			fmt.Scan(&cipherChoice)
			if cipherChoice >= 0 && cipherChoice <= 6 {
				break
			}
			fmt.Println("the wrong choice of cryptosystem, try again:")
		}

		if cipherChoice == 0 {
			fmt.Println("Ending process")
			break
		}

		// Выбор операции
		fmt.Println("Choose the operation:")
		fmt.Println("1: Encryption")
		fmt.Println("2: Decryption")

		var operationChoice int
		for {
			fmt.Scan(&operationChoice)
			if operationChoice == 1 || operationChoice == 2 {
				break
			}
			fmt.Println("the wrong choice of operation, try again:")
		}

		switch cipherChoice {
		case 1:
			// Caesar cipher
			key, err := verify.CaesarKey(keyString, alphabetMap)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if operationChoice == 1 {
				result = encrypt.Caesar(input, key, alphabetMap, power)
				WriteToFile(encryptFile, result)

			} else {
				result = decrypt.Caesar(input, key, alphabetMap, power)
				WriteToFile(decryptFile, result)
			}
		case 2:
			key, err := verify.AffineKey(keyString, alphabetMap)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// Affine cipher
			if operationChoice == 1 {
				result = encrypt.Affine(input, key, alphabetMap, power)
				WriteToFile(encryptFile, result)
			} else {
				result = decrypt.Affine(input, key, alphabetMap, power)
				WriteToFile(decryptFile, result)
			}
		case 3:
			key, err := verify.SubstitutionKey(keyString, alphabetMap, power)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// Simple substitution
			if operationChoice == 1 {
				result = encrypt.Substitution(input, key, alphabetMap, power)
				WriteToFile(encryptFile, result)
			} else {
				result = decrypt.Substitution(input, []rune(keyString), alphabetMap, power)
				WriteToFile(decryptFile, result)

			}
		case 4:
			// Hill cipher
			if operationChoice == 1 {
				result = hillEncrypt(input, keyString)
				WriteToFile(encryptFile, result)
			} else {
				result = hillDecrypt(input, keyString)
				WriteToFile(decryptFile, result)
			}
		case 5:
			err := verify.PermutationKey(keyString, alphabetMap, power)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// Permutation cipher
			if operationChoice == 1 {
				result = encrypt.Permutation(input, keyString)
				WriteToFile(encryptFile, result)
			} else {
				result = decrypt.Permutation(input, keyString)
				WriteToFile(decryptFile, result)
			}
		case 6:
			err := verify.VigenereKey(keyString, alphabetMap, power)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// Vigenere cipher
			if operationChoice == 1 {
				result = encrypt.Vigenere(input, keyString, alphabetMap, power)
				WriteToFile(encryptFile, result)
			} else {
				result = decrypt.Vigenere(input, keyString, alphabetMap, power)
				WriteToFile(decryptFile, result)
			}
		default:
			fmt.Println("Wrong choise")
			continue
		}
	}
}

func permutationEncrypt(input string, key []int, alphabetMap map[rune]int, power int) string {
	// Реализация шифра простой замены
	return input
}

func permutationDecrypt(input string, key []int, alphabetMap map[rune]int, power int) string {
	// Реализация дешифрования шифра простой замены
	return input
}

func hillEncrypt(input, key string) string {
	// Реализация шифра Хилла
	return input
}

func hillDecrypt(input, key string) string {
	// Реализация дешифрования шифра Хилла
	return input
}
