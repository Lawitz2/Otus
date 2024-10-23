package hw03frequencyanalysis

import (
	"cmp"
	"regexp"
	"slices"
	"strings"
)

var regWords = regexp.MustCompile(`(?:[a-zа-я]+[[:punct:]]*)+[a-zа-я]+`)

type word struct {
	word  string
	count int
}

func Top10(text string) []string {
	text = strings.ToLower(text)
	textSlc := strings.Fields(text)
	freqMap := make(map[string]int)

	for _, v := range textSlc {
		if v == "-" {
			continue
		}

		switch {
		case regWords.MatchString(v):
			freqMap[regWords.FindString(v)]++
		default:
			freqMap[v]++
		}
	}

	words := make([]word, 0, len(freqMap))

	for key, val := range freqMap {
		words = append(words, word{
			word:  key,
			count: val,
		})
	}

	slices.SortFunc(words, func(a, b word) int {
		return cmp.Or(
			cmp.Compare(b.count, a.count),
			cmp.Compare(a.word, b.word))
	})

	wordsOutput := make([]string, 0, min(len(words), 10))

	for _, w := range words[:min(len(words), 10)] {
		wordsOutput = append(wordsOutput, w.word)
	}

	return wordsOutput
}
