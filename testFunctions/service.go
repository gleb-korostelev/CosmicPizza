package testfunctions

import (
	"sync"

	"github.com/gleb-korostelev/CosmicPizza.git/config"
	cosmicorder "github.com/gleb-korostelev/CosmicPizza.git/service/cosmicOrder"
	fanout "github.com/gleb-korostelev/CosmicPizza.git/service/fanOut"
	ingredienttree "github.com/gleb-korostelev/CosmicPizza.git/service/ingredientTree"
	"github.com/gleb-korostelev/CosmicPizza.git/service/worker"
	"github.com/gleb-korostelev/CosmicPizza.git/utils"
)

func WorkerPoolSimulation(orderNumber, ingredientNumber int, orderList *cosmicorder.CosmicOrderList, ingredientTree *ingredienttree.IngredientTree) {
	// Order Worker Pool initialization
	orderWorkerPool := worker.NewWorkerPool(config.MaxConcurrentWorkerPoolOperations)
	defer orderWorkerPool.Shutdown()

	var wg sync.WaitGroup

	for i := 1; i <= orderNumber; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			order := utils.GenerateRandomOrder(id)
			orderWorkerPool.AddTask(utils.ProcessOrder(orderList, order))
		}(i)
	}

	// Ingredient Worker Pool initialization
	ingredientWorkerPool := worker.NewWorkerPool(config.MaxConcurrentWorkerPoolOperations)
	defer ingredientWorkerPool.Shutdown()

	for i := 0; i < ingredientNumber; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ingredient := utils.GenerateRandomIngredient()
			ingredientWorkerPool.AddTask(utils.ProcessIngredient(ingredientTree, ingredient))
		}()
	}

	wg.Wait()
}

func TryInsertBadIndexOrder(orderList *cosmicorder.CosmicOrderList) {
	orderList.InsertOrder(10000000, 3, "Venus", "Quantum Anchoa")
}

func TryInsertSameIngredients(ingredientTree *ingredienttree.IngredientTree) {
	ingredientTree.Insert(7)
	ingredientTree.Insert(7)
}

func FullConcurrencySimulation(fanoutWorkerNumber int, orderList *cosmicorder.CosmicOrderList, ingredientTree *ingredienttree.IngredientTree) {
	// Initialize WorkerPool service
	workerPool := worker.NewWorkerPool(config.MaxConcurrentWorkerPoolOperations)
	defer workerPool.Shutdown()

	var wg sync.WaitGroup

	// Generate tasks dynamically
	tasks := utils.GenerateTasks()

	doneCh := make(chan struct{})
	defer close(doneCh)

	// Create input channel for tasks
	inputCh := utils.ChanGenerator(tasks, doneCh)

	// Initialize FanOutService
	fanOut := fanout.NewFanOutService(inputCh, fanoutWorkerNumber)
	defer fanOut.Shutdown()

	// Start processing tasks
	processed := utils.ProcessTasks(workerPool, orderList, ingredientTree, fanOut.GetOutputChannels(), &wg)

	// Collect final results
	utils.CollectResults(orderList, ingredientTree, processed)
}
