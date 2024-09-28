package hw02unpackstring

import (
	"errors"
	"log/slog"
	"slices"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(text string) (string, error) {
	if len(text) < 1 { // empty string check
		return "", nil
	}

	textRunes := []rune(text)

	numbers := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

	if slices.Contains(numbers, rune(text[0])) { // if the 1st symbol is a number - invalid string
		return "", ErrInvalidString
	}

	strBuilder := strings.Builder{}
	var symbol rune
	var index int
	afterDigit := false

	for index, symbol = range textRunes[1:] {
		if slices.Contains(numbers, symbol) {
			if afterDigit { // if two numbers in a row - invalid string
				return "", ErrInvalidString
			}

			repeats, err := strconv.Atoi(string(symbol))
			if err != nil {
				slog.Error(err.Error())
				return "", err
			}

			strBuilder.WriteString(strings.Repeat(string(textRunes[index]), repeats))
			afterDigit = true
		} else {
			if afterDigit {
				afterDigit = false
				continue
			}
			strBuilder.WriteRune(textRunes[index])
		}
	}

	if !slices.Contains(numbers, symbol) { // append last symbol of the input text if it is not a number
		strBuilder.WriteRune(symbol) // (if it is a number - it'll be handled by the for loop)
	}

	return strBuilder.String(), nil
}
