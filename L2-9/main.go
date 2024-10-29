package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

/*
L2.9 «Взаимодействие с ОС»
Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:

- cd <args> — смена директории (в качестве аргумента могут быть то-то и то);

- pwd — показать путь до текущего каталога;

- echo <args> — вывод аргумента в STDOUT;

- kill <args> — «убить» процесс, переданный в качесте аргумента (пример: такой-то пример);

- ps — выводит общую информацию по запущенным процессам в формате такой-то формат.

Так же требуется поддерживать функционал fork/exec-команд.

Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).

*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике и при необходимости выводит результат на экран. Интерактивный сеанс поддерживается до тех пор, пока не будет введена команда выхода (например \quit).
*/

// Функция для выполнения команды cd
func changeDirectory(args []string) {
	if len(args) < 1 {
		fmt.Println("должен быть указан путь")
		return
	}
	dir := args[0]
	err := os.Chdir(dir)
	if err != nil {
		fmt.Println(err)
	}
}

// Функция для выполнения команды pwd
func printWorkingDirectory() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(dir)
	}
}

// Функция для выполнения команды echo
func echo(args []string) {
	fmt.Println(strings.Join(args, " "))
}

// Функция для завершения процесса
func killProcess(args []string) {
	if len(args) < 1 {
		fmt.Println("должен быть указан идентификатор процесса")
		return
	}
	pid := args[0]
	err := exec.Command("kill", pid).Run()
	if err != nil {
		fmt.Println(err)
	}
}

// Функция для вывода запущенных процессов
func listProcesses() {
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(output))
	}
}

// Функция для обработки команд
func executeCommand(command string) {
	args := strings.Fields(command)

	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "cd":
		changeDirectory(args[1:])
	case "pwd":
		printWorkingDirectory()
	case "echo":
		echo(args[1:])
	case "kill":
		killProcess(args[1:])
	case "ps":
		listProcesses()
	default:
		executeExternalCommand(args)
	}
}

// Функция для выполнения внешних команд
func executeExternalCommand(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// Основной цикл
func main() {
	currentUser, _ := user.Current()
	fmt.Printf("Добро пожаловать, %s! Для выхода введите 'quit'.\n", currentUser.Username)

	for {
		// Выводим приглашение
		fmt.Printf("%s$ ", currentUser.HomeDir)
		var input string
		_, err := fmt.Scanln(&input)

		if err != nil {
			continue
		}

		if input == "quit" {
			break
		}

		executeCommand(input)
	}
}
