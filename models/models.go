package models

// Order represents a pizza order
type Order struct {
	OrderID   int
	Planet    string
	PizzaType string
	Next      *Order
}

// Task represents a unit of work (order or ingredient operation)
type Task struct {
	Type       int // Task type (Add, Remove, Insert, Search)
	OrderID    int // Used for order operations
	Planet     string
	PizzaType  string
	Ingredient int // Used for ingredient operations
}
