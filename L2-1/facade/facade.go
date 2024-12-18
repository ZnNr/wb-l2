package main

import "fmt"

/*
Фасад — это структурный паттерн, который добавляет простой интерфейс к сложной системе для взаимодействия с ней.

Паттерн Фасад используется, когда:
есть сложная система, работу с которой нужно упростить;
хочется уменьшить количество зависимостей между клиентом и сложной системой;
требуется разбить сложную систему на компоненты — применение паттерна к каждому компоненту упростит взаимодействие между ними.

*/

//пример, когда вы заказываете пиццу, вы не углубляетесь в процесс приготовления и доставки,
//а о многих других действиях в цепочке производства даже не задумываетесь.

// Pizza — сложная система с простым интерфейсом.
type Pizza struct {
	site     Site
	kitchen  Kitchen
	delivery Delivery
}

// Site отвечает за получение заказа на сайте.
type Site struct {
}

func (s *Site) GetOrder() {
	fmt.Println("Оформляем заказ.")
}

func (s *Site) RedirectOrder() {
	fmt.Println("Отправляем заказ на кухню.")
}

func (s *Site) CloseOrder() {
	fmt.Println("Заказ выполнен.")
}

// Kitchen готовит пиццу.
type Kitchen struct {
}

func (k *Kitchen) CookPizza() {
	fmt.Println("Готовим пиццу.")
}

// Delivery занимается доставкой пиццы.
type Delivery struct {
}

func (d *Delivery) SearchCourier() {
	fmt.Println("Ищем курьера.")
}

func (d *Delivery) DeliverPizza() {
	fmt.Println("Доставляем пиццу.")
}

func (d *Delivery) GetPayment() {
	fmt.Println("Получаем оплату.")
}

// StartOrder — метод для заказа пиццы.
func (p *Pizza) StartOrder() {
	p.site.GetOrder()
	p.site.RedirectOrder()
	p.kitchen.CookPizza()
	p.delivery.SearchCourier()
	p.delivery.DeliverPizza()
}

// FinishOrder — выполнение заказа.
func (p *Pizza) FinishOrder() {
	p.delivery.GetPayment()
	p.site.CloseOrder()
}

func main() {
	pizza := &Pizza{}

	pizza.StartOrder()
	pizza.FinishOrder()
}
