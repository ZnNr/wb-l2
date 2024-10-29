package main

import "fmt"

/*
Паттерн Команда преобразует все параметры операции или события в объект-команду.
Впоследствии можно выполнить эту операцию, вызвав соответствующий метод объекта.
Объект-команда заключает в себе всё необходимое для проведения операции,
поэтому её легко выполнять, логировать и отменять.

Паттерн Команда применяется, когда:
нужно преобразовать операции в объекты, которые можно обрабатывать и хранить:
использование объектов вместо операций позволяет создавать очереди,
передавать команды дальше или выполнять их в нужный момент;

требуется реализовать операцию отмены выполненных действий.



*/
// Command Interface
type Command interface {
	Execute()
}

// Receiver
type Light struct {
	IsOn bool
}

func (l *Light) TurnOn() {
	l.IsOn = true
	fmt.Println("Light is turned ON")
}

func (l *Light) TurnOff() {
	l.IsOn = false
	fmt.Println("Light is turned OFF")
}

// Concrete Commands
type TurnOnCommand struct {
	light *Light
}

func (c *TurnOnCommand) Execute() {
	c.light.TurnOn()
}

type TurnOffCommand struct {
	light *Light
}

func (c *TurnOffCommand) Execute() {
	c.light.TurnOff()
}

// Invoker
type RemoteControl struct {
	command Command
}

func (r *RemoteControl) SetCommand(c Command) {
	r.command = c
}

func (r *RemoteControl) PressButton() {
	r.command.Execute()
}

func main() {
	light := &Light{}

	onCommand := &TurnOnCommand{light: light}
	offCommand := &TurnOffCommand{light: light}

	remote := &RemoteControl{}

	remote.SetCommand(onCommand)
	remote.PressButton()

	remote.SetCommand(offCommand)
	remote.PressButton()
}
