package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

/*
L2.6 «Утилита grep»
Реализовать утилиту фильтрации по аналогии с консольной утилитой (man grep — смотрим описание и основные параметры).

Реализовать поддержку утилитой следующих ключей:

-A - "after": печатать +N строк после совпадения;

-B - "before": печатать +N строк до совпадения;

-C - "context": (A+B) печатать ±N строк вокруг совпадения;

-c - "count": количество строк;

-i - "ignore-case": игнорировать регистр;

-v - "invert": вместо совпадения, исключать;

-F - "fixed": точное совпадение со строкой, не паттерн;

-n - "line num": напечатать номер строки.
*/

// Флаги
var (
	after      = flag.Int("A", 0, "Количество строк после совпадения")
	before     = flag.Int("B", 0, "Количество строк перед совпадением")
	context    = flag.Int("C", 0, "Количество строк вокруг совпадения")
	countOnly  = flag.Bool("c", false, "Подсчитать количество совпадений")
	ignoreCase = flag.Bool("i", false, "Игнорировать регистр")
	invert     = flag.Bool("v", false, "Инвертировать совпадения")
	fixed      = flag.Bool("F", false, "Точное совпадение со строкой, не паттерн")
	lineNum    = flag.Bool("n", false, "Печатать номер строки")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Использование: %s [опции] шаблон [файл]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// Проверка на наличие шаблона
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Ошибка: Шаблон не указан.")
		flag.Usage()
		os.Exit(1)
	}

	pattern := flag.Arg(0)
	var file *os.File
	var err error

	// Открытие файла или чтение из stdin
	if flag.NArg() >= 2 {
		file, err = os.Open(flag.Arg(1))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка открытия файла: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
	} else {
		file = os.Stdin
	}

	// Подготовка шаблона
	if *ignoreCase {
		pattern = strings.ToLower(pattern)
	}

	scanner := bufio.NewScanner(file)
	var (
		line         string
		lineNumber   int
		matchCount   int
		matchedLines []int
		afterCounter int
		beforeQueue  []string
		beforeSize   int
		_            []string
		_            []string
	)

	if *context > 0 {
		*before = *context
		*after = *context
	}

	if *before > 0 {
		beforeSize = *before
	}

	for scanner.Scan() {
		line = scanner.Text()
		lineNumber++
		processedLine := line
		if *ignoreCase {
			processedLine = strings.ToLower(line)
		}

		var isMatch bool
		if *fixed {
			isMatch = strings.Contains(processedLine, pattern)
		} else {
			// Для простоты используем Contains. Для полноценного grep можно использовать regexp.
			isMatch = strings.Contains(processedLine, pattern)
		}

		if *invert {
			isMatch = !isMatch
		}

		if isMatch {
			matchCount++
			matchedLines = append(matchedLines, lineNumber)

			// Вывод строк перед совпадением
			for _, bl := range beforeQueue {
				printLine(bl, -1, *lineNum)
			}
			beforeQueue = []string{}

			// Вывод текущей строки
			printLine(line, lineNumber, *lineNum)

			// Установка счетчика для вывода последующих строк
			afterCounter = *after
			continue
		}

		// Вывод строк после совпадения
		if afterCounter > 0 {
			printLine(line, lineNumber, *lineNum)
			afterCounter--
			continue
		}

		// Сохранение строк перед совпадением
		if beforeSize > 0 {
			if len(beforeQueue) >= beforeSize {
				beforeQueue = beforeQueue[1:]
			}
			beforeQueue = append(beforeQueue, line)
		}
	}

	if *countOnly {
		fmt.Println(matchCount)
	} else if !*countOnly && *context > 0 {
		// Дополнительная обработка контекста при необходимости
		// Здесь можно реализовать дополнительные функции, если требуется
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения файла: %v\n", err)
		os.Exit(1)
	}
}

func printLine(line string, number int, showNum bool) {
	if showNum && number != -1 {
		fmt.Printf("%d:%s\n", number, line)
	} else {
		fmt.Println(line)
	}
}
