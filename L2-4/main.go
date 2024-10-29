package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

/*
L2.4 «Утилита sort»
Отсортировать строки в файле по аналогии с консольной утилитой sort (man sort — смотрим описание и основные параметры): на входе подается файл из несортированными строками, на выходе — файл с отсортированными.

Реализовать поддержку утилитой следующих ключей:

-k — указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел);

-n — сортировать по числовому значению;

-r — сортировать в обратном порядке;

-u — не выводить повторяющиеся строки.

Дополнительно
Реализовать поддержку утилитой следующих ключей:

-M — сортировать по названию месяца;

-b — игнорировать хвостовые пробелы;

-c — проверять отсортированы ли данные;

-h — сортировать по числовому значению с учетом суффиксов.
*/

// Структура для хранения опций сортировки
type SortOptions struct {
	Key        int  // Номер колонки для сортировки (начинается с 1)
	Numeric    bool // Сортировать по числовому значению (-n)
	Reverse    bool // Обратный порядок (-r)
	Unique     bool // Удалить дубликаты (-u)
	Month      bool // Сортировать по названию месяца (-M)
	IgnoreTail bool // Игнорировать хвостовые пробелы (-b)
	Check      bool // Проверить, отсортирован ли файл (-c)
	Human      bool // Сортировка с учётом суффиксов (-h)
}

func main() {
	// Определяем флаги
	key := flag.Int("k", 0, "Указать номер колонки для сортировки (начинается с 1)")
	numeric := flag.Bool("n", false, "Сортировать по числовому значению")
	reverse := flag.Bool("r", false, "Сортировать в обратном порядке")
	unique := flag.Bool("u", false, "Не выводить повторяющиеся строки")
	month := flag.Bool("M", false, "Сортировать по названию месяца")
	ignoreTail := flag.Bool("b", false, "Игнорировать хвостовые пробелы")
	check := flag.Bool("c", false, "Проверить, отсортированы ли данные")
	human := flag.Bool("h", false, "Сортировать по числовому значению с учётом суффиксов")

	flag.Parse()

	// Проверяем наличие входного файла
	if flag.NArg() < 1 {
		fmt.Println("Использование: sort [опции] <файл>")
		os.Exit(1)
	}

	inputFile := flag.Arg(0)

	// Читаем строки из файла
	lines, err := readLines(inputFile)
	if err != nil {
		fmt.Printf("Ошибка чтения файла: %v\n", err)
		os.Exit(1)
	}

	options := SortOptions{
		Key:        *key,
		Numeric:    *numeric,
		Reverse:    *reverse,
		Unique:     *unique,
		Month:      *month,
		IgnoreTail: *ignoreTail,
		Check:      *check,
		Human:      *human,
	}

	if options.Check {
		sorted, err := isSorted(lines, options)
		if err != nil {
			fmt.Printf("Ошибка проверки сортировки: %v\n", err)
			os.Exit(1)
		}
		if sorted {
			fmt.Println("Файл отсортирован.")
			os.Exit(0)
		} else {
			fmt.Println("Файл не отсортирован.")
			os.Exit(1)
		}
	}

	// Удаляем хвостовые пробелы, если указано
	if options.IgnoreTail {
		for i, line := range lines {
			lines[i] = strings.TrimRightFunc(line, unicode.IsSpace)
		}
	}

	// Сортируем строки
	sort.Slice(lines, func(i, j int) bool {
		a := getSortKey(lines[i], options)
		b := getSortKey(lines[j], options)

		// Сравнение
		var less bool
		if options.Numeric || options.Human {
			aNum, errA := parseNumber(a, options)
			bNum, errB := parseNumber(b, options)
			if errA != nil || errB != nil {
				// Если не удалось преобразовать в число, сравниваем как строки
				less = a < b
			} else {
				less = aNum < bNum
			}
		} else if options.Month {
			aMonth := monthToNumber(a)
			bMonth := monthToNumber(b)
			less = aMonth < bMonth
		} else {
			less = a < b
		}

		if options.Reverse {
			return !less
		}
		return less
	})

	// Удаляем дубликаты, если необходимо
	if options.Unique {
		lines = uniqueLines(lines)
	}

	// Выводим отсортированные строки
	for _, line := range lines {
		fmt.Println(line)
	}
}

// Чтение всех строк из файла
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		line = strings.TrimRight(line, "\n\r")
		lines = append(lines, line)
		if err == io.EOF {
			break
		}
	}
	return lines, nil
}

// Получение ключа сортировки для строки
func getSortKey(line string, options SortOptions) string {
	if options.Key <= 0 {
		return line
	}
	fields := strings.Fields(line)
	if options.Key > len(fields) {
		return ""
	}
	return fields[options.Key-1]
}

// Преобразование месяца в номер
func monthToNumber(month string) int {
	months := map[string]int{
		"January":   1,
		"February":  2,
		"March":     3,
		"April":     4,
		"May":       5,
		"June":      6,
		"July":      7,
		"August":    8,
		"September": 9,
		"October":   10,
		"November":  11,
		"December":  12,
	}
	monthCap := strings.Title(strings.ToLower(month))
	if num, exists := months[monthCap]; exists {
		return num
	}
	return 0
}

// Парсинг числового значения с учётом суффиксов
func parseNumber(s string, options SortOptions) (float64, error) {
	if options.Human {
		// Регулярное выражение для чисел с суффиксами
		re := regexp.MustCompile(`^([+-]?[\d\.]+)([KMGTP]?)$`)
		matches := re.FindStringSubmatch(s)
		if len(matches) != 3 {
			return 0, errors.New("неверный формат числа с суффиксом")
		}
		num, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, err
		}
		switch matches[2] {
		case "K":
			num *= 1e3
		case "M":
			num *= 1e6
		case "G":
			num *= 1e9
		case "T":
			num *= 1e12
		case "P":
			num *= 1e15
		}
		return num, nil
	} else {
		return strconv.ParseFloat(s, 64)
	}
}

// Удаление дубликатов из отсортированного среза строк
func uniqueLines(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	unique := []string{lines[0]}
	for _, line := range lines[1:] {
		if line != unique[len(unique)-1] {
			unique = append(unique, line)
		}
	}
	return unique
}

// Проверка, отсортирован ли срез строк
func isSorted(lines []string, options SortOptions) (bool, error) {
	for i := 1; i < len(lines); i++ {
		a := getSortKey(lines[i-1], options)
		b := getSortKey(lines[i], options)

		var cmp int
		if options.Numeric || options.Human {
			aNum, errA := parseNumber(a, options)
			bNum, errB := parseNumber(b, options)
			if errA != nil || errB != nil {
				cmp = strings.Compare(a, b)
			} else {
				if aNum < bNum {
					cmp = -1
				} else if aNum > bNum {
					cmp = 1
				} else {
					cmp = 0
				}
			}
		} else if options.Month {
			aMonth := monthToNumber(a)
			bMonth := monthToNumber(b)
			if aMonth < bMonth {
				cmp = -1
			} else if aMonth > bMonth {
				cmp = 1
			} else {
				cmp = 0
			}
		} else {
			cmp = strings.Compare(a, b)
		}

		if options.Reverse {
			cmp = -cmp
		}

		if cmp > 0 {
			return false, nil
		}
	}
	return true, nil
}
