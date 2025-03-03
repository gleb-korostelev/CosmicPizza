# **Cosmic Pizza - Task Processing System**

## **Overview**

Cosmic Pizza is a concurrent task processing system that utilizes a **fan-out, worker pool, and fan-in** approach to efficiently manage and process orders and ingredients. It ensures high throughput by distributing tasks across multiple workers and aggregating results effectively.

## **Architecture**

The system consists of the following components:

### **1. FanOutService**

- Distributes incoming tasks (orders/ingredients) among multiple worker channels.
- Ensures load balancing by spreading tasks across multiple goroutines.

### **2. WorkerPool**

- Processes each task in a concurrent manner.
- Executes `SwitchProcessTasks` to ensure tasks are fully completed before forwarding them.

### **3. Task Flow**

1. Tasks are **generated** and sent to `FanOutService`.
2. `FanOutService` **distributes tasks** across multiple worker channels.
3. `WorkerPool` **processes** each task and ensures its completion.
4. FanIn flow **aggregates results** and makes them available for collection.
5. `CollectResults()` **finalizes the processed data**.

---

## **Installation & Setup**

### **1. Clone the Repository**

```sh
git clone https://github.com/gleb-korostelev/CosmicPizza.git
cd CosmicPizza
```

### **2. Install Dependencies**

Ensure you have **Go 1.23+** installed.

```sh
go mod tidy
```

### **3. Run the Simulation**

```sh
go run cmd/main.go
```

---

## **Usage**

### **Running Full Simulation**

The main function to start the concurrency simulation:

```go
	testfunctions.FullConcurrencySimulation(config.NumberOfWorkersForFunOut, orderList, ingredientTree)
```