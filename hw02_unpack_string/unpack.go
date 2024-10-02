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
	textRunes := []rune(text)
	if len(textRunes) < 1 { // empty input check
		return "", nil
	}

	numbers := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	if slices.Contains(numbers, textRunes[0]) {
		return "", ErrInvalidString
	}
	strBuilder := strings.Builder{}
	var escape, prtFlag, numFlag bool
	var output rune

	if textRunes[0] == '\\' {
		escape = true
	} else {
		output = textRunes[0]
		prtFlag = true
	}

	for _, symbol := range textRunes[1:] {
		if escape {
			switch {
			case symbol != '\\' && !slices.Contains(numbers, symbol):
				return "", ErrInvalidString
			default:
				output = symbol
				escape = false
				continue
			}
		}

		switch {
		case symbol == '\\':
			numFlag = false
			escape = true
			if !prtFlag {
				prtFlag = true
				continue
			}
			strBuilder.WriteRune(output)
		case slices.Contains(numbers, symbol):
			if numFlag {
				return "", ErrInvalidString
			}
			repeats, err := strconv.Atoi(string(symbol))
			if err != nil {
				slog.Error(err.Error())
			}
			strBuilder.WriteString(strings.Repeat(string(output), repeats))
			prtFlag = false
			numFlag = true
		default:
			numFlag = false
			if !prtFlag {
				prtFlag = true
				output = symbol
				continue
			}
			strBuilder.WriteRune(output)
			output = symbol
		}
	}

	switch { // handling last symbol (print it out if it wasn't, throw error if escaping end of string)
	case escape:
		return "", ErrInvalidString
	case prtFlag:
		strBuilder.WriteRune(output)
	}

	return strBuilder.String(), nil
}
