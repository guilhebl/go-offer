package offer

import (
	"github.com/gorilla/mux"
	"github.com/guilhebl/go-worker-pool"
	"log"
	"runtime"
	"sync"
)

// centralized module manager which holds references to JobQueue and other global app scoped objects
// Singleton enforcing the module will be initialized at max. once per app.
type Module struct {
	Dispatcher job.WorkerPool
	JobQueue   chan job.Job
	Router     *mux.Router
}

var instance *Module
var once sync.Once

func BuildInstance(router *mux.Router) *Module {
	once.Do(func() {
		instance = newModule(router)
	})
	return instance
}

func GetInstance() *Module {
	return instance
}

func newModule(router *mux.Router) *Module {
	log.Printf("New Module")

	// fetch ENV var param ?
	// maxWorker := os.Getenv("MAX_WORKERS")

	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs + 1) // numCPUs hot threads + one for async tasks.
	maxWorkers := numCPUs * 4

	jobQueue := make(chan job.Job)

	module := Module{
		Dispatcher: job.NewWorkerPool(maxWorkers),
		JobQueue:   jobQueue,
		Router:     router,
	}

	// A buffered channel that we can send work requests on.
	module.Dispatcher.Run(jobQueue)
	return &module
}

// stops pool and closes JobQueue returns the result of closing both
func (m *Module) Stop() bool {
	log.Printf("%s", "Stopping Module")
	m.Dispatcher.Stop()

	// close the Job queue chan
	close(m.JobQueue)

	// empty queue
	for x := range m.JobQueue {
		_ = x // ignore channel var using blank identifier
	}

	// Make sure that the function does close the channel
	_, ok := <-m.JobQueue

	return ok
}
