package cosmicorder

import (
	"sync"

	"github.com/gleb-korostelev/CosmicPizza.git/models"
	"github.com/gleb-korostelev/CosmicPizza.git/tools/logger"
)

// CosmicOrderList represents a linked list of orders
type CosmicOrderList struct {
	head *models.Order
	mu   sync.Mutex
}

// NewSerice creates a new instance of the CosmicOrder service
func NewService() *CosmicOrderList {
	return &CosmicOrderList{}
}

// AddOrder adds an order to the end of the list
func (s *CosmicOrderList) AddOrder(orderID int, planet, pizzaType string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newOrder := &models.Order{OrderID: orderID, Planet: planet, PizzaType: pizzaType}
	if s.head == nil {
		s.head = newOrder
		return
	}

	current := s.head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newOrder
}

// InsertOrder inserts an order at a specific index
func (s *CosmicOrderList) InsertOrder(index, orderID int, planet, pizzaType string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newOrder := &models.Order{OrderID: orderID, Planet: planet, PizzaType: pizzaType}
	if index == 0 {
		newOrder.Next = s.head
		s.head = newOrder
		return
	}

	current := s.head
	for i := 0; i < index-1 && current != nil; i++ {
		current = current.Next
	}

	if current == nil {
		logger.Infof("Cannot insert the order - index is bigger than the list")
		return
	}

	newOrder.Next = current.Next
	current.Next = newOrder
}

// RemoveOrder removes an order by its orderID
func (s *CosmicOrderList) RemoveOrder(orderID int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.head == nil {
		return
	}

	if s.head.OrderID == orderID {
		s.head = s.head.Next
		return
	}

	current := s.head
	for current.Next != nil && current.Next.OrderID != orderID {
		current = current.Next
	}

	if current.Next != nil {
		current.Next = current.Next.Next
	}
}
