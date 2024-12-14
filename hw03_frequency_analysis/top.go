package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(inText string) []string {
	type wordFrequency struct {
		field  string
		fields []string
		value  int
	}
	if len(inText) == 0 {
		return nil
	}

	// Разобьем текст на слова
	allFields := strings.Fields(inText)
	// Слова по частоте вхождения
	listCnt := make(map[string]int, len(allFields))
	for _, field := range allFields {
		if cnt, ok := listCnt[field]; ok {
			listCnt[field] = cnt + 1
		} else {
			listCnt[field] = 1
		}
	}

	//  Перегоним в структуру для сортировки
	listFreq := make([]wordFrequency, len(allFields))
	i := 0
	for key, value := range listCnt {
		listFreq[i].value = value
		listFreq[i].field = key
		i++
	}
	// Отсортируем
	sort.Slice(listFreq, func(i, j int) bool {
		return listFreq[i].value > listFreq[j].value
	})
	// Заберем первые 10 слов
	prepResult := make(map[int][]string)
	for i = 0; i < 10; i++ {
		idx := listFreq[i].value
		if strs, ok := prepResult[idx]; ok {
			prepResult[idx] = append(strs, listFreq[i].field)
		} else {
			prepResult[idx] = []string{listFreq[i].field}
		}
	}
	listFreq = make([]wordFrequency, len(prepResult))
	i = 0
	for key, value := range prepResult {
		listFreq[i].value = key
		listFreq[i].fields = value
		i++
	}
	// Итоговая сортировка
	sort.Slice(listFreq, func(i, j int) bool {
		return listFreq[i].value > listFreq[j].value
	})
	var result []string
	for _, value := range listFreq {
		sort.Strings(value.fields)
		result = append(result, value.fields...)
	}
	return result
}
