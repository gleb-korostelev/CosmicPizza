package fanout

import (
	"sync"

	"github.com/gleb-korostelev/CosmicPizza.git/models"
)

// FanOutService manages multiple workers that process data from an input channel
type FanOutService struct {
	numWorkers int
	inputCh    chan models.Task
	doneCh     chan struct{}
	outputChs  []chan models.Task
	wg         sync.WaitGroup
}

// NewFanOutService initializes a new FanOutService
func NewFanOutService(inputch chan models.Task, numWorkers int) *FanOutService {
	fanOut := &FanOutService{
		numWorkers: numWorkers,
		inputCh:    inputch,
		doneCh:     make(chan struct{}),
		outputChs:  make([]chan models.Task, numWorkers),
	}

	// Create worker goroutines
	for i := 0; i < numWorkers; i++ {
		outputCh := make(chan models.Task)
		fanOut.outputChs[i] = outputCh
		go fanOut.worker(outputCh)
	}

	return fanOut
}

// worker processes data from the input channel and sends it to an output channel
func (s *FanOutService) worker(outputCh chan models.Task) {
	defer close(outputCh)

	for task := range s.inputCh {
		select {
		case <-s.doneCh:
			// logger.Infof("FanOutService: Stopping worker due to shutdown signal.")
			return
		case outputCh <- task:
			// logger.Infof("FanOutService: Task sent to output channel: %+v", task)
		}
	}
}

// AddData sends data to the input channel for processing
func (s *FanOutService) AddData(value models.Task) {
	select {
	case <-s.doneCh: // If the service is stopped, ignore new data
		return
	default:
		s.inputCh <- value
	}
}

// GetOutputChannels returns the output channels of the workers
func (s *FanOutService) GetOutputChannels() []chan models.Task {
	return s.outputChs
}

// Shutdown gracefully stops all workers and closes channels
func (s *FanOutService) Shutdown() {
	close(s.doneCh) // Signal all workers to stop

	s.wg.Wait() // Wait for all workers to finish
	// close(s.inputCh)
}
