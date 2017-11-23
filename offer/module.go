package offer

import (
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
}

var instance *Module
var once sync.Once

func GetInstance() *Module {
	once.Do(func() {
		instance = newModule()
	})
	return instance
}

func newModule() *Module {
	log.Printf("%s", "New Module")

	// fetch ENV var param ?
	// maxWorker := os.Getenv("MAX_WORKERS")

	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs + 1) // numCPUs hot threads + one for async tasks.
	maxWorkers := numCPUs * 4

	jobQueue := make(chan job.Job)

	module := Module{
		Dispatcher: job.NewWorkerPool(maxWorkers),
		JobQueue:   jobQueue,
	}

	// A buffered channel that we can send work requests on.
	module.Dispatcher.Run(jobQueue)
	return &module
}
