package main

import (
	"fmt"
	"sync"
	"time"
)

/*
L2.8 «Or channel»
Реализовать функцию, которая будет объединять один или более done-каналов в single-канал, если один из его составляющих каналов закроется.

Очевидным вариантом решения могло бы стать выражение при использовании select, которое бы реализовывало эту связь, однако иногда неизвестно общее число done-каналов, с которыми вы работаете в рантайме. В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or-каналов, реализовывала бы весь функционал.

Определение функции:

var or func(channels ...<- chan interface{}) <- chan interface{}
Пример использования функции:

sig := func(after time.Duration) <- chan interface{} {
c := make(chan interface{})
go func() {
defer close(c)
time.Sleep(after)
}()
return c
}
start := time.Now()
<-or (
sig(2*time.Hour),
sig(5*time.Minute),
sig(1*time.Second),
sig(1*time.Hour),
sig(1*time.Minute),
)
fmt.Printf(“fone after %v”, time.Since(start))
*/
// or объединяет один или более done-каналов в single-канал.
// Если один из переданных каналов закрывается, результирующий канал также закрывается.
var or func(channels ...<-chan interface{}) <-chan interface{} = func(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	output := make(chan interface{})
	var once sync.Once

	// Функция для закрытия output единожды
	closeOutput := func() {
		once.Do(func() {
			close(output)
		})
	}

	// Запускаем горутину для каждого входного канала
	for _, ch := range channels {
		go func(c <-chan interface{}) {
			select {
			case <-c:
				closeOutput()
			}
		}(ch)
	}

	return output
}

// Пример использования функции or
func main() {
	// Функция sig возвращает канал, который закроется через заданное время
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("Функция завершена через %v\n", time.Since(start))
}

/*
Пояснение
1. Функция or:
- Случай 0 каналов: Возвращаем nil, поскольку нет каналов для отслеживания.
- Случай 1 канала: Возвращаем этот канал напрямую.
- Случай нескольких каналов:
    - Создаем новый канал output.
    - Используем sync.Once для гарантии, что output будет закрыт только один раз.
    - Для каждого входного канала запускаем горутину, которая ждет закрытия своего канала. Как только канал закрывается, вызывается closeOutput, который закрывает output.
2. Функция sig:
- Создает канал, который закроется через заданное время.
- Это используется для имитации различных done-каналов с разными задержками.
3. Функция main:
- Создает несколько done-каналов с разными временными задержками.
- Использует функцию or для объединения этих каналов.
- Блокируется до тех пор, пока любой из переданных каналов не будет закрыт.
- Выводит время, прошедшее с момента запуска, после завершения.
4. Тестирование
Запустив приведенный код, вы получите вывод, указывающий, что функция завершилась через примерно 1 секунду, поскольку самый быстрый канал (sig(1*time.Second)) закрывается первым:

Копировать
Функция завершена через 1.001234567s

*/
