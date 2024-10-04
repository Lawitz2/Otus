package hw03frequencyanalysis

import (
	"slices"
	"strings"
)

func Top10(text string) []string {
	freqMap := make(map[string]int)

	textSlice := strings.Fields(text)
	slices.Sort(textSlice)

	for _, val := range textSlice {
		freqMap[val]++
	}

	freqSlice := make([]int, 0, len(freqMap))

	for _, val := range freqMap {
		freqSlice = append(freqSlice, val)
	}

	slices.SortFunc(freqSlice, func(a, b int) int {
		return b - a
	})
	freqSlice = freqSlice[:min(10, len(freqSlice))]

	result := make([]string, 0, min(len(freqMap), 10))
	box := make([]string, 0, cap(result))

	for i := 0; i < cap(result); {
		for key, val := range freqMap {
			if val == freqSlice[i] {
				box = append(box, key)
				delete(freqMap, key)
			}
		}
		slices.Sort(box)
		result = append(result, box...)
		i += len(box)
		box = box[len(box):]
	}
	return result
}
