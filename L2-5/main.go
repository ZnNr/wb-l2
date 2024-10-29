package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
L2.5 «Поиск анаграмм по словарю»
Написать функцию поиска всех множеств анаграмм по словарю.

Например:

'пятак', 'пятка' и 'тяпка' — принадлежат одному множеству;

'листок', 'слиток' и 'столик' — другому.

Требования
Входные данные для функции: ссылка на массив, каждый элемент которого — слово на русском языке в кодировке utf8.

Выходные данные: ссылка на мапу множеств анаграмм.

Ключ — первое встретившееся в словаре слово из множества. Значение — ссылка на массив, каждый элемент которого, слово из множества.

Массив должен быть отсортирован по возрастанию.

Множества из одного элемента не должны попасть в результат.

Все слова должны быть приведены к нижнему регистру.

В результате каждое слово должно встречаться только один раз.
*/

// FindAnagramSets ищет все множества анаграмм в словаре.
// Входные данные: ссылка на срез слов (каждое слово на русском языке, utf8).
// Выходные данные: ссылка на мапу множеств анаграмм.
// Ключ — первое встретившееся в словаре слово из множества.
// Значение — ссылка на отсортированный срез слов из множества.
func FindAnagramSets(words *[]string) *map[string]*[]string {
	anagramMap := make(map[string][]string)
	uniqueWords := make(map[string]struct{}) // Для удаления дубликатов

	// Итерация по словам
	for _, word := range *words {
		// Приведение к нижнему регистру
		lowerWord := strings.ToLower(word)

		// Удаление дубликатов
		if _, exists := uniqueWords[lowerWord]; exists {
			continue
		}
		uniqueWords[lowerWord] = struct{}{}

		// Создание ключа для анаграмм: сортировка букв
		sortedKey := sortStringByRunes(lowerWord)

		// Группировка слов по ключу
		anagramMap[sortedKey] = append(anagramMap[sortedKey], lowerWord)
	}

	// Формирование результирующей карты
	result := make(map[string]*[]string)
	for _, group := range anagramMap {
		if len(group) < 2 {
			continue // Исключаем группы с одним элементом
		}

		// Сортировка группы по возрастанию
		sort.Strings(group)

		// Используем первое слово группы в качестве ключа
		key := group[0]
		result[key] = &group
	}

	return &result
}

// sortStringByRunes сортирует символы строки по возрастанию.
func sortStringByRunes(s string) string {
	// Преобразуем строку в срез рун
	runes := []rune(s)

	// Сортировка
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})

	return string(runes)
}

// Пример использования функции
func main() {
	words := []string{
		"пятак",
		"пятка",
		"тяпка",
		"листок",
		"слиток",
		"столик",
		"дом",
		"мод",
		"тест",
	}

	anagramSets := FindAnagramSets(&words)

	// Вывод результатов
	for key, group := range *anagramSets {
		fmt.Printf("Анаграммы для '%s': %v\n", key, *group)
	}
}
