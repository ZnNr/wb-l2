package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
L2.7 «Утилита cut»
Реализовать утилиту аналог консольной команды cut (man cut). Утилита должна принимать строки через STDIN, разбивать по разделителю (TAB) на колонки и выводить запрошенные.

Реализовать поддержку утилитой следующих ключей:

-f — "fields": выбрать поля (колонки);

-d — "delimiter": использовать другой разделитель;

-s — "separated": только строки с разделителем.
*/
func main() {
	// Определение флагов
	fieldsFlag := flag.String("f", "", "Выбранные поля (колонки), разделенные запятыми.")
	delimiterFlag := flag.String("d", "\t", "Разделитель, по умолчанию - табуляция.")
	separatedFlag := flag.Bool("s", false, "Выводить только строки с разделителем.")
	flag.Parse()

	// Проверка, указаны ли поля
	if *fieldsFlag == "" {
		fmt.Println("Ошибка: нужно указать поля для вывода с помощью -f")
		os.Exit(1)
	}

	fields := strings.Split(*fieldsFlag, ",")
	delimiter := *delimiterFlag

	// Чтение из стандартного ввода
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// Проверка на наличие разделителя, если указано
		if *separatedFlag && !strings.Contains(line, delimiter) {
			continue
		}

		// Разделение строки по разделителю
		columns := strings.Split(line, delimiter)

		// Формируем вывод
		output := []string{}
		for _, field := range fields {
			index, err := parseField(field, len(columns))
			if err == nil {
				output = append(output, columns[index])
			}
		}

		// Выводим результат
		fmt.Println(strings.Join(output, delimiter))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка при чтении входных данных:", err)
	}
}

// parseField парсит указанные поля и возвращает индекс
func parseField(field string, maxColumns int) (int, error) {
	index := 0
	var err error
	if strings.HasPrefix(field, "-") {
		return -1, errors.New("недопустимое поле: " + field)
	}

	// Преобразование в целое значение
	if field != "*" {
		index, err = strconv.Atoi(field)
		if err != nil {
			return -1, err
		}
		index-- // Поскольку индексы начинаются с 0
	}

	if index < 0 || index >= maxColumns {
		return -1, errors.New("индекс вне диапазона")
	}

	return index, nil
}

/*
Как использовать утилиту:
Компилируйте программу, например, с именем cut.
Запустите ее, передавая входные данные через стандартный ввод:

echo -e "col1\tcol2\tcol3\nval1\tval2\tval3" | ./cut -f 1,3 -d $'\t'

Указывает поля -f 1,3, которые необходимо вывести, и использует разделитель табуляции по умолчанию.
*/
