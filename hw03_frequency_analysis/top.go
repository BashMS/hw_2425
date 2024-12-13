package hw03frequencyanalysis

import "strings"

func Top10(inText string) []string {
	if len(inText) == 0 {
		return nil
	}

	// Разобьем текст на слова
	allFields := strings.Fields(inText)
	listCnt := make(map[string]int, len(allFields))
	for _, field := range allFields {
		if cnt, ok := listCnt[field]; ok {
			listCnt[field] = cnt + 1
		} else {
			listCnt[field] = 1
		}
	}

	return nil
}
