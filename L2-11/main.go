package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

/*
L2.11 «Утилита telnet»
Реализовать простейший telnet-клиент.

Примеры вызовов:

go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123
Требования
Программа должна подключаться к указанному хосту (ip или доменное имя + порт) по протоколу TCP. После подключения STDIN программы должен записываться в сокет, а данные полученные из сокета должны выводиться в STDOUT.

Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться. При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

func main() {
	// Определение аргументов командной строки
	timeout := flag.Duration("timeout", 10*time.Second, "Timeout for connection")
	host := flag.String("host", "", "Host (IP or domain name)")
	port := flag.String("port", "", "Port")
	flag.Parse()

	// Проверка на наличие необходимых аргументов
	if *host == "" || *port == "" {
		fmt.Println("Usage: go-telnet --timeout=<duration> host port")
		return
	}

	// Формирование адреса
	address := net.JoinHostPort(*host, *port)

	// Попытка установить соединение с сервером
	conn, err := net.DialTimeout("tcp", address, *timeout)
	if err != nil {
		fmt.Printf("Failed to connect to %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", address)

	// Использование горoutine для передачи данных из сокета в stdout
	go func() {
		io.Copy(os.Stdout, conn)
	}()

	// Основной поток передает данные из stdin в сокет
	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		fmt.Println("Connection closed.")
	}
}

/*
Описание кода
1. Аргументы командной строки: используем пакет flag для обработки аргументов командной строки, включая таймаут, хост и порт.
2. Подключение: С помощью метода net.DialTimeout пытаемся установить TCP-соединение с указанным адресом и проверяем, нет ли ошибок.
3. Копирование данных:
- запускаем горутину, которая перенаправляет данные от сокета в STDOUT при помощи io.Copy.
- Основной поток выполняет io.Copy, чтобы перенаправить данные из STDIN в сокет.
4. Завершение: Программа завершится, если соединение будет закрыто со стороны сервера или если пользователь нажмет Ctrl+D.

Запуск программы

go run go-telnet.go --timeout=10s mysite.ru 8080
Это создаст простейший telnet-клиент, который будет подключаться к указанному серверу и порту, передавать данные между STDIN и STDOUT.
*/
