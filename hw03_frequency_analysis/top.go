package hw03frequencyanalysis

import (
	"regexp"
	"slices"
	"strings"
)

var regWords = regexp.MustCompile(`(?:[a-zа-я]+[[:punct:]]*)+[a-zа-я]+`)

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

	appearancesSlc := make([]int, 0, len(freqMap))

	for _, val := range freqMap {
		appearancesSlc = append(appearancesSlc, val)
	}

	slices.SortFunc(appearancesSlc, func(a, b int) int {
		return b - a
	})

	resSlice := make([]string, 0, len(freqMap))

	for i := 0; i < min(len(appearancesSlc), 10); {
		k := 0
		check := appearancesSlc[i]
		for key, val := range freqMap {
			if check == val {
				resSlice = append(resSlice, key)
				k++
			}
		}
		slices.Sort(resSlice[i:])
		i += k
	}

	return resSlice[:min(len(resSlice), 10)]
}
