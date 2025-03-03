package main

import (
	"github.com/gleb-korostelev/CosmicPizza.git/config"
	cosmicorder "github.com/gleb-korostelev/CosmicPizza.git/service/cosmicOrder"
	ingredienttree "github.com/gleb-korostelev/CosmicPizza.git/service/ingredientTree"
	testfunctions "github.com/gleb-korostelev/CosmicPizza.git/testFunctions"
	"github.com/gleb-korostelev/CosmicPizza.git/tools/closer"
	"github.com/gleb-korostelev/CosmicPizza.git/tools/logger"
)

func main() {
	// Defer the cleanup of all resources using the closer utility.
	defer func() {
		closer.Wait()
		closer.CloseAll()
	}()
	orderList := cosmicorder.NewService()
	ingredientTree := ingredienttree.NewService()

	// This are some test functions a did
	// testfunctions.WorkerPoolSimulation(config.OrderNumber, config.IngredientNumber, orderList, ingredientTree)
	// testfunctions.TryInsertSameIngredients(ingredientTree)
	// testfunctions.TryInsertBadIndexOrder(orderList)

	// This is the main function of the project
	testfunctions.FullConcurrencySimulation(config.NumberOfWorkersForFunOut, orderList, ingredientTree)

	// calculates min max and the sum of all ingredients
	min, max, sum := ingredientTree.FindMinMaxSum()

	logger.Infof("Final Orders List Processed")
	logger.Infof("Final Ingredient Tree Values: %v", ingredientTree.TraverseInOrder())
	logger.Infof("Final max/min/sum of values: %v, %v, %v", min, max, sum)
}
