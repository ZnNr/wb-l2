package unpack

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
)

/*
L2.3 «Задача на распаковку»
Создать Go-функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы/руны.

Например:

"a4bc2d5e" => "aaaabccddddde"

"abcd" => "abcd"

"45" => "" (некорректная строка)

"" => ""

Дополнительно
Реализовать поддержку escape-последовательностей.

Например:

qwe\4\5 => qwe45 (*)

qwe\45 => qwe44444 (*)

qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка, функция должна возвращать ошибку.

Написать unit-тесты.
*/

// Unpack осуществляет распаковку строки, учитывая повторяющиеся символы и escape-последовательности.
// Возвращает распакованную строку или ошибку, если входная строка некорректна.
// Unpack осуществляет распаковку строки, учитывая повторяющиеся символы и escape-последовательности.
// Возвращает распакованную строку или ошибку, если входная строка некорректна.
func Unpack(input string) (string, error) {
	var result strings.Builder
	var prev rune
	var escape bool
	var numberBuf strings.Builder

	for i, ch := range input {
		if escape {
			// Добавляем символ без обработки
			result.WriteRune(ch)
			prev = ch
			escape = false
			continue
		}

		if ch == '\\' {
			escape = true
			continue
		}

		if unicode.IsDigit(ch) {
			if prev == 0 {
				return "", errors.New("некорректная строка: цифра без предшествующего символа")
			}

			// Собираем все цифры для многозначных чисел
			numberBuf.Reset()
			numberBuf.WriteRune(ch)
			j := i + 1
			for j < len(input) {
				nextCh := rune(input[j])
				if unicode.IsDigit(nextCh) {
					numberBuf.WriteRune(nextCh)
					j++
				} else {
					break
				}
			}

			count, err := strconv.Atoi(numberBuf.String())
			if err != nil {
				return "", errors.New("некорректная строка: не удалось преобразовать цифру")
			}

			if count == 0 {
				// Удаляем последний добавленный символ, так как count = 0
				s := result.String()
				if len(s) > 0 {
					result.Reset()
					result.WriteString(s[:len(s)-1])
				}
				prev = 0
			} else {
				// Повторяем предыдущий символ (count - 1) раз
				for j := 0; j < count-1; j++ {
					result.WriteRune(prev)
				}
			}

			// Пропускаем обработанные цифры в основном цикле
			for k := i + 1; k < j; k++ {
				// Просто пропускаем, так как они уже обработаны
				i = k
			}

			continue
		}

		// Добавляем текущий символ в результат
		result.WriteRune(ch)
		prev = ch
	}

	if escape {
		return "", errors.New("некорректная строка: завершена escape-последовательность")
	}

	return result.String(), nil
}

func main() {
	testStrings := []string{
		`a4bc2d5e`,
		`abcd`,
		`45`,
		`qwe\4\5`,
		`qwe\45`,
		`qwe\\5`,
		`\`,
	}

	for _, s := range testStrings {
		result, err := Unpack(s)
		if err != nil {
			log.Printf("Ошибка распаковки строки %q: %v\n", s, err)
		} else {
			fmt.Printf("Распакованная строка: %q\n", result)
		}
	}
}
