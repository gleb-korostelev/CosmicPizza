package ingredienttree

import (
	"sync"

	"github.com/gleb-korostelev/CosmicPizza.git/tools/logger"
)

// IngredientTree represents a binary search tree for ingredients
type IngredientTree struct {
	value int
	left  *IngredientTree
	right *IngredientTree
	mu    sync.Mutex
}

func NewService() *IngredientTree {
	return &IngredientTree{}
}

// Insert adds a new ingredient to the tree
func (s *IngredientTree) Insert(value int) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// If tree is empty or default initialized (0 root), replace it with the first inserted value
	if s.value == 0 && s.left == nil && s.right == nil {
		s.value = value
		return
	}

	if value < s.value {
		if s.left == nil {
			s.left = &IngredientTree{value: value}
		} else {
			s.left.Insert(value)
		}
	} else if value > s.value {
		if s.right == nil {
			s.right = &IngredientTree{value: value}
		} else {
			s.right.Insert(value)
		}
	} else if value == s.value {
		logger.Infof("Value is already in the tree")

	}
}

// Search checks if an ingredient exists
func (s *IngredientTree) Search(value int) bool {
	if s == nil {
		return false
	}

	if value < s.value {
		return s.left.Search(value)
	} else if value > s.value {
		return s.right.Search(value)
	}

	return true
}

// TraverseInOrder returns sorted ingredient values
func (s *IngredientTree) TraverseInOrder() []int {
	if s == nil {
		return []int{}
	}
	var result []int
	if s.left != nil {
		result = append(result, s.left.TraverseInOrder()...)
	}
	result = append(result, s.value)
	if s.right != nil {
		result = append(result, s.right.TraverseInOrder()...)
	}
	return result
}

// FindMinMaxSum finds the minimum, maximum, and sum of all ingredients in the tree.
func (s *IngredientTree) FindMinMaxSum() (min int, max int, sum int) {
	if s == nil {
		return 0, 0, 0
	}

	// Finding minimum value (leftmost node)
	current := s
	for current.left != nil {
		current = current.left
	}
	min = current.value

	// Finding maximum value (rightmost node)
	current = s
	for current.right != nil {
		current = current.right
	}
	max = current.value

	// Calculating the sum of all values
	sum = s.calculateSum()

	return min, max, sum
}

// calculateSum recursively sums up all values in the tree.
func (s *IngredientTree) calculateSum() int {
	if s == nil {
		return 0
	}
	return s.value + s.left.calculateSum() + s.right.calculateSum()
}
