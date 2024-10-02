package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Not enough arguments")
		return
	} else if len(os.Args) > 3 {
		fmt.Println("Too many arguments")
		return
	}
	inputFileName := os.Args[1]
	outputFileName := os.Args[2]

	// Read the input file
	data, err := os.ReadFile(inputFileName)
	if checkError(err) {
		return
	}
	inputStringData := string(data) // Convert bytes to string

	inputStringData = strings.Replace(inputStringData, "(cap, ", "(cap,", -1)
	inputStringData = strings.Replace(inputStringData, "(low, ", "(low,", -1)
	inputStringData = strings.Replace(inputStringData, "(up, ", "(up,", -1)
	inputStringData = strings.Replace(inputStringData, "' ", "'", -1)
	inputStringData = strings.Replace(inputStringData, " '", "'", -1)

	inputWordsArray := strings.Fields(inputStringData)
	var outputWordsArray []string

	// Vowel map for "a" to "an" conversion
	vowelsMap := map[string]bool{
		"a": true, "e": true, "i": true, "o": true, "u": true, "h": true,
		"A": true, "E": true, "I": true, "O": true, "U": true, "H": true,
	}

	// Process each word
	for wordNumber := 0; wordNumber < len(inputWordsArray); wordNumber++ {
		word := inputWordsArray[wordNumber]

		if word == "(hex)" {
			lastWord := outputWordsArray[len(outputWordsArray)-1]
			outputWordsArray[len(outputWordsArray)-1] = convertHexToDecimal(lastWord)
		} else if word == "(bin)" {
			lastWord := outputWordsArray[len(outputWordsArray)-1]
			outputWordsArray[len(outputWordsArray)-1] = convertBinToDecimal(lastWord)
		} else if word == "(up)" {
			lastWord := outputWordsArray[len(outputWordsArray)-1]
			outputWordsArray[len(outputWordsArray)-1] = strings.ToUpper(lastWord)
		} else if strings.HasPrefix(word, "(up,") {
			countToConvert, err := strconv.Atoi(strings.TrimSpace(word[4 : len(word)-1])) // Extract number
			if checkError(err) {
				return
			}
			for i := countToConvert; i > 0; i-- {
				outputWordsArray[len(outputWordsArray)-i] = strings.ToUpper(string(inputWordsArray[wordNumber-i]))
			}
		} else if word == "(low)" {
			lastWord := outputWordsArray[len(outputWordsArray)-1]
			outputWordsArray[len(outputWordsArray)-1] = strings.ToLower(lastWord)
		} else if strings.HasPrefix(word, "(low,") {
			countToConvert, err := strconv.Atoi(strings.TrimSpace(word[5 : len(word)-1])) // Extract number
			if checkError(err) {
				return
			}
			for i := countToConvert; i > 0; i-- {
				outputWordsArray[len(outputWordsArray)-i] = strings.ToLower(string(inputWordsArray[wordNumber-i]))
			}
		} else if word == "(cap)" {
			lastWord := outputWordsArray[len(outputWordsArray)-1]
			outputWordsArray[len(outputWordsArray)-1] = capitalize(lastWord)
		} else if strings.HasPrefix(word, "(cap,") {
			countToConvert, err := strconv.Atoi(strings.TrimSpace(word[5 : len(word)-1])) // Extract number
			if checkError(err) {
				return
			}
			for i := countToConvert; i > 0; i-- {
				outputWordsArray[len(outputWordsArray)-i] = capitalize(string(inputWordsArray[wordNumber-i]))
			}
		} else if len(outputWordsArray) > 0 && strings.ToLower(outputWordsArray[len(outputWordsArray)-1]) == "a" && vowelsMap[string(word[0])] {
			// Change "a" to "an"
			outputWordsArray[len(outputWordsArray)-1] += "n" // Change "a" or "A" to "an" or "An"
			outputWordsArray = append(outputWordsArray, word)
		} else {
			outputWordsArray = append(outputWordsArray, word) // Append the current word normally
		}
	}

	// Handle punctuation and formatting
	outputWordsArray = handlePunctuation(outputWordsArray)

	// Write output to file
	outputString := strings.TrimSpace(strings.Join(outputWordsArray, " "))

	err = os.WriteFile(outputFileName, []byte(outputString), 0644)
	checkError(err)
}

// Convert hex to decimal
func convertHexToDecimal(hexStr string) string {
	decimal, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		return hexStr
	}
	return strconv.FormatInt(decimal, 10)
}

// Convert binary to decimal
func convertBinToDecimal(binStr string) string {
	decimal, err := strconv.ParseInt(binStr, 2, 64)
	if err != nil {
		return binStr
	}
	return strconv.FormatInt(decimal, 10)
}

// Capitalize the first letter of a word
func capitalize(word string) string {
	if len(word) > 1 {
		return strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
	}
	return strings.ToUpper(word)
}

// Check for errors
func checkError(err error) bool {
	if err != nil {
		fmt.Println("Error:", err)
		return true
	}
	return false
}

// Handle punctuation and single quotes
func handlePunctuation(words []string) []string {
	var result []string
	punctuationMarks := map[string]bool{
		".": true, ",": true, "!": true, "?": true, ":": true, ";": true,
	}

	for i := 0; i < len(words); i++ {
		word := strings.TrimSpace(words[i]) // Trim spaces around the word

		// Handle group punctuations
		// If the ("..." or "!?") punctuation is a separate word
		if word == "..." || word == "!?" {
			if len(result) > 0 {
				// Attach it to the previous word if it exists
				result[len(result)-1] = result[len(result)-1] + word
			} else {
				// If it's the first word, add it as-is
				result = append(result, word)
			}
			continue
		}

		if punctuationMarks[word] { // Case 1: Punctuation as a separate word
			// Attach punctuation to the last word
			if len(result) > 0 {
				result[len(result)-1] = result[len(result)-1] + word
			}
		} else if len(word) > 1 && punctuationMarks[string(word[0])] { // Case 2: Punctuation attached to the next word
			punctuation := string(word[0]) // Extract the punctuation from the beginning of the word
			trimmedWord := word[1:]        // Extract the word after punctuation

			// Attach the punctuation to the previous word
			if len(result) > 0 {
				result[len(result)-1] = result[len(result)-1] + punctuation
			}

			// Add the trimmed word to the result
			result = append(result, trimmedWord)
		} else {
			// Normal word append it to the result
			result = append(result, word)
		}
	}

	var inQuotes bool // Track if we're inside quotes
	// Handle single quotes
	for i := 0; i < len(result); i++ {
		if result[i] == "'" { // Case 1: quote as a separate word
			if inQuotes {
				// Attach closing quote to the previous word
				if i > 0 {
					result[i-1] = strings.TrimSpace(result[i-1]) + "'" // No space before
					result[i] = ""
				}
				inQuotes = false
			} else {
				// Open quote, add it to the result
				inQuotes = true
				if i < len(result)-1 {
					result[i+1] = "'" + strings.TrimSpace(result[i+1]) // No space after
					result[i] = ""
				}
			}
		} else if strings.Contains(result[i], "'") { // Case 2: quote is part of the word
			if inQuotes {
				inQuotes = false
			} else {
				inQuotes = true
			}
		}
	}

	return result
}
