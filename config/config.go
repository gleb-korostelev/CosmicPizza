package config

// List of constants
const (
	// MaxConcurrentUpdates defines the maximum number of concurrent operations for workerPool.
	MaxConcurrentWorkerPoolOperations = 10

	// OrderProcessTime in milliseconds for concurrency tests
	OrderProcessTime = 500

	// IngredientProcessTime in milliseconds for concurrency tests
	IngredientProcessTime = 300

	// MaxIngredientNumber is a maximum random ingredient number in data type int
	MaxIngredientNumber = 100

	// OrderNumber number  of orders to be randomized for workerPool
	OrderNumber = 5

	// IngredientNumber number of ingredients to be randomized for workerPool
	IngredientNumber = 5

	// OrderTaskNumber number of generated orders for GenerateTasks
	OrderTaskNumber = 5

	// IngredientTaskNumber number of generated ingredients for GenerateTasks
	IngredientTaskNumber = 6

	// NumberOfWorkersForFunOut for fanout
	NumberOfWorkersForFunOut = 5
)
