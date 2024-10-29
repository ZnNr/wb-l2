package main

import "fmt"

/*
Aбстрактная фабрика — это порождающий паттерн проектирования, который позволяет создавать семейство связанных объектов. Паттерн представляет собой абстракцию над шаблоном Фабрика.
Паттерн используется, когда:
бизнес-логика программы должна работать с семействами связанных продуктов;
в программе уже используется Фабричный метод, но нужно добавить новые типы продуктов.

*/

// ChairMaker — интерфейс для кресел.
type ChairMaker interface {
	SetStyle(string)
	SetWood(string)
	Print() string
}

// TableMaker — интерфейс для столиков.
type TableMaker interface {
	SetStyle(string)
	SetWood(string)
	Print() string
}

// Chair — абстрактное кресло.
type Chair struct {
	Style string
	Wood  string
}

func (c *Chair) SetStyle(style string) {
	c.Style = style
}

func (c *Chair) SetWood(wood string) {
	c.Wood = wood
}

func (c *Chair) Print() string {
	return fmt.Sprintf("Кресло [Стиль: %s, Дерево: %s]", c.Style, c.Wood)
}

// Table — абстрактный столик.
type Table struct {
	Style string
	Wood  string
}

func (t *Table) SetStyle(style string) {
	t.Style = style
}

func (t *Table) SetWood(wood string) {
	t.Wood = wood
}

func (t *Table) Print() string {
	return fmt.Sprintf("Столик [Стиль: %s, Дерево: %s]", t.Style, t.Wood)
}

// Factory — абстрактная фабрика.
type Factory interface {
	MakeChair(string) ChairMaker
	MakeTable(string) TableMaker
}

// После того как определили абстрактные кресло и столик, опишем кресло и стол в стиле ар-деко и фабрику для их производства.

// ArtDecoChair — кресло ар-деко.
type ArtDecoChair struct {
	Chair
}

// ArtDecoTable — столик ар-деко.
type ArtDecoTable struct {
	Table
}

// ArtDeco — фабрика ар-деко.
type ArtDeco struct {
}

func (a *ArtDeco) MakeChair(wood string) ChairMaker {
	var chair ArtDecoChair
	chair.SetStyle("ар-деко")
	chair.SetWood(wood)
	return &chair
}

func (a *ArtDeco) MakeTable(wood string) TableMaker {
	var table ArtDecoTable
	table.SetStyle("ар-деко")
	table.SetWood(wood)
	return &table
}

//Точно так же описываем кресло и столик в стиле модерн и производящую их фабрику.

// ModernChair — кресло модерн.
type ModernChair struct {
	Chair
}

// ModernTable — столик модерн.
type ModernTable struct {
	Table
}

// Modern — фабрика модерна.
type Modern struct {
}

func (m *Modern) MakeChair(wood string) ChairMaker {
	var chair ModernChair
	chair.SetStyle("модерн")
	chair.SetWood(wood)
	return &chair
}

func (m *Modern) MakeTable(wood string) TableMaker {
	var table ModernTable
	table.SetStyle("модерн")
	table.SetWood(wood)
	return &table
}

// создаём абстрактную фабрику и проверяем её работу.
// GetFactory — абстрактная фабрика.
func GetFactory(style string) Factory {
	if style == "art-deco" {
		return &ArtDeco{}
	}
	if style == "modern" {
		return &Modern{}
	}
	return nil
}

func main() {
	artdecoFactory := GetFactory("art-deco")
	modernFactory := GetFactory("modern")

	artdecoChair := artdecoFactory.MakeChair("дуб")
	artdecoTable := artdecoFactory.MakeTable("дуб")

	modernChair := modernFactory.MakeChair("ясень")
	modernTable := modernFactory.MakeTable("ясень")

	fmt.Println(artdecoChair.Print(), artdecoTable.Print())
	fmt.Println(modernChair.Print(), modernTable.Print())
}

/*
Выведет:

Кресло [Стиль: ар-деко, Дерево: дуб] Столик [Стиль: ар-деко, Дерево: дуб]
Кресло [Стиль: модерн, Дерево: ясень] Столик [Стиль: модерн, Дерево: ясень]

*/
