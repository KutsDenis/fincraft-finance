package domain

import "fmt"

// Money представляет денежную сумму в тиынах, центах и т.д.
type Money int64

// NewMoneyFromFloat создает новый объект денег из float64
func NewMoneyFromFloat(amount float64) Money {
	return Money(amount * 100)
}

// ToFloat возвращает значение денег в float64
func (m Money) ToFloat() float64 {
	return float64(m) / 100
}

// String возвращает значение денег в виде строки
func (m Money) String() string {
	return fmt.Sprintf("%.2f", m.ToFloat())
}
