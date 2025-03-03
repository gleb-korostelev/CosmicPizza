package utils

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/gleb-korostelev/CosmicPizza.git/config"
	"github.com/gleb-korostelev/CosmicPizza.git/models"
	cosmicorder "github.com/gleb-korostelev/CosmicPizza.git/service/cosmicOrder"
	ingredienttree "github.com/gleb-korostelev/CosmicPizza.git/service/ingredientTree"
	"github.com/gleb-korostelev/CosmicPizza.git/service/worker"
	"github.com/gleb-korostelev/CosmicPizza.git/tools/logger"
)

// Random order data
var planets = []string{"Mars", "Venus", "Jupiter", "Saturn", "Neptune", "Pluto", "Andromeda Nebula"}
var pizzaTypes = []string{"BlackHole Pepperoni", "Galactic Cheese", "Quantum Anchoa", "Nebula Deluxe", "Supernova Supreme", "Dark Matter Veggie", "Antimatter Pizza"}

// Task types for fan-out processing
const (
	AddOrderTask    = 1
	RemoveOrderTask = 2
	InsertIngTask   = 3
	SearchIngTask   = 4
)

// ProcessOrder add's order in orderList and processing it
func ProcessOrder(orderList *cosmicorder.CosmicOrderList, order models.Order) worker.Task {
	return worker.Task{
		Action: func(ctx context.Context) error {
			orderList.AddOrder(order.OrderID, order.Planet, order.PizzaType)
			logger.Infof("Processed order #%d from %s: %s", order.OrderID, order.Planet, order.PizzaType)
			time.Sleep(config.OrderProcessTime * time.Millisecond) // Job simulation
			return nil
		},
		Done: make(chan struct{}),
	}
}

// ProcessIngredient adds new ingridient to the tree
func ProcessIngredient(tree *ingredienttree.IngredientTree, ingredient int) worker.Task {
	return worker.Task{
		Action: func(ctx context.Context) error {
			tree.Insert(ingredient + 1)
			logger.Infof("Inserted ingredient: %d", ingredient)
			time.Sleep(config.IngredientProcessTime * time.Millisecond) // Job simulation
			return nil
		},
		Done: make(chan struct{}),
	}
}

// GenerateRandomOrder creates a random order
func GenerateRandomOrder(orderID int) models.Order {
	return models.Order{
		OrderID:   orderID,
		Planet:    planets[rand.Intn(len(planets))],
		PizzaType: pizzaTypes[rand.Intn(len(pizzaTypes))],
		Next:      nil,
	}
}

// GenerateRandomIngredient creates a random magical ingredient
func GenerateRandomIngredient() int {
	return rand.Intn(config.MaxIngredientNumber) + 1 // Random number between 1 and 100
}

// ProcessTasks reads from the output channels and executes corresponding actions
func ProcessTasks(workerPool *worker.WorkerPool, orderList *cosmicorder.CosmicOrderList, ingredientTree *ingredienttree.IngredientTree, outputChs []chan models.Task, wg *sync.WaitGroup) chan models.Task {
	processedCh := make(chan models.Task)

	go func() {
		defer close(processedCh)
		for _, ch := range outputChs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for task := range ch {
					// logger.Infof("ProcessTasks task is being processed %v", task)
					workerPool.AddTask(worker.Task{
						Action: func(ctx context.Context) error {
							err := SwitchProcessTasks(ctx, task, orderList, ingredientTree)
							return err
						},
						Done: make(chan struct{}),
					})
					processedCh <- task

				}
			}()
		}
		wg.Wait()
	}()
	return processedCh
}

func SwitchProcessTasks(ctx context.Context, task models.Task, orderList *cosmicorder.CosmicOrderList, ingredientTree *ingredienttree.IngredientTree) error {
	switch task.Type {
	case AddOrderTask:
		orderList.AddOrder(task.OrderID, task.Planet, task.PizzaType)
		// logger.Infof("Added Order #%d from %s: %s", task.OrderID, task.Planet, task.PizzaType)
	case RemoveOrderTask:
		orderList.RemoveOrder(task.OrderID)
		// logger.Infof("Removed Order #%d", task.OrderID)
	case InsertIngTask:
		ingredientTree.Insert(task.Ingredient)
		// logger.Infof("Inserted Ingredient: %d", task.Ingredient)
	case SearchIngTask:
		_ = ingredientTree.Search(task.Ingredient)
		// logger.Infof("Searched Ingredient %d: Found? %v", task.Ingredient, found)
	}
	return nil
}

// Generate random tasks for orders and ingredients
func GenerateTasks() []models.Task {
	tasks := []models.Task{}

	// Generate random orders
	for i := 1; i <= config.OrderTaskNumber; i++ {
		order := GenerateRandomOrder(i)
		tasks = append(tasks, models.Task{
			Type:      AddOrderTask,
			OrderID:   order.OrderID,
			Planet:    order.Planet,
			PizzaType: order.PizzaType,
		})
	}

	// Remove half random orders
	for i := 1; i <= config.OrderTaskNumber/2; i++ {
		tasks = append(tasks, models.Task{
			Type:    RemoveOrderTask,
			OrderID: rand.Intn(config.OrderTaskNumber) + 1,
		})
	}

	// Generate random ingredient insertions
	for i := 0; i < config.IngredientTaskNumber; i++ {
		tasks = append(tasks, models.Task{
			Type:       InsertIngTask,
			Ingredient: GenerateRandomIngredient(),
		})
		// logger.Infof("Ingredient generated %d", tasks[i].Ingredient)
	}

	// Search for some ingredients
	for i := 0; i < config.IngredientTaskNumber; i++ {
		tasks = append(tasks, models.Task{
			Type:       SearchIngTask,
			Ingredient: GenerateRandomIngredient(),
		})
	}

	return tasks
}

func ChanGenerator(tasks []models.Task, doneCh chan struct{}) chan models.Task {

	inputCh := make(chan models.Task)

	go func() {
		defer close(inputCh)

		for _, task := range tasks {
			select {
			case <-doneCh:
				return
			case inputCh <- task:
				// logger.Infof("ChanGenerator: Sent task to FanOut: %+v", task)
			}
		}
	}()

	return inputCh
}

// CollectResults gathers all results in the main thread
func CollectResults(orderList *cosmicorder.CosmicOrderList, ingredientTree *ingredienttree.IngredientTree, outputCh chan models.Task) {
	remainingOrders := []models.Order{}
	remainingIngredients := []int{}
	antimatterPizzaFound := false

	for task := range outputCh {
		switch task.Type {
		case AddOrderTask:
			remainingOrders = append(remainingOrders, models.Order{
				OrderID:   task.OrderID,
				Planet:    task.Planet,
				PizzaType: task.PizzaType,
				Next:      nil,
			})
		case RemoveOrderTask:
			// Remove from remaining orders
			for i, order := range remainingOrders {
				if order.OrderID == task.OrderID {
					remainingOrders = append(remainingOrders[:i], remainingOrders[i+1:]...)
					break
				}
			}
		case InsertIngTask:
			remainingIngredients = append(remainingIngredients, task.Ingredient)
		}

		// Check if "Antimatter Pizza" is still in the menu
		if task.PizzaType == "Antimatter Pizza" {
			antimatterPizzaFound = true
		}
	}

	logger.Infof("Final Orders in List: %v", remainingOrders)
	logger.Infof("Final Ingredients in Tree: %v", remainingIngredients)
	if antimatterPizzaFound {
		logger.Infof("Antimatter Pizza is still available!")
	} else {
		logger.Infof("Warning: Antimatter Pizza is missing from the menu!")
	}
}
