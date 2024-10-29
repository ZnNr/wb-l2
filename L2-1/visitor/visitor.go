package main

import "fmt"

/*
Посетитель (Visitor) — даёт возможность добавлять новый функционал к объектам, не внося в них изменения;

Шаблон Посетитель позволяет отвязать функциональность от объекта.
Новые методы добавляются не для каждого типа из семейства,
а для промежуточного объекта visitor, аккумулирующего функциональность.
Типам семейства добавляется только один метод accept(visitor).
*/

// Visitor Interface
type Visitor interface {
	VisitConcreteElementA(*ConcreteElementA)
	VisitConcreteElementB(*ConcreteElementB)
}

// Element Interface
type Element interface {
	Accept(Visitor)
}

// Конкретные элементы
type ConcreteElementA struct {
	Name string
}

func (e *ConcreteElementA) Accept(v Visitor) {
	v.VisitConcreteElementA(e)
}

type ConcreteElementB struct {
	Age int
}

func (e *ConcreteElementB) Accept(v Visitor) {
	v.VisitConcreteElementB(e)
}

// Конкретный посетитель
type ConcreteVisitor struct{}

func (v *ConcreteVisitor) VisitConcreteElementA(e *ConcreteElementA) {
	fmt.Printf("Visitor: %s\n", e.Name)
}

func (v *ConcreteVisitor) VisitConcreteElementB(e *ConcreteElementB) {
	fmt.Printf("Visitor: %d\n", e.Age)
}

func main() {
	elements := []Element{
		&ConcreteElementA{Name: "Alice"},
		&ConcreteElementB{Age: 30},
	}

	visitor := &ConcreteVisitor{}

	for _, e := range elements {
		e.Accept(visitor)
	}
}
