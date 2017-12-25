package offer

import (
	"github.com/gorilla/mux"
	"github.com/guilhebl/go-offer/common/config"
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

func BuildInstance(router *mux.Router, mode string) *Module {
	once.Do(func() {
		instance = newModule(router, mode)
	})
	return instance
}

func GetInstance() *Module {
	return instance
}

// Builds a new module which is a container for the running app instance
// router - the router configuration with URL routes and mapped action handlers
// mode - test or production modes, which will make the app read from either test or prod config properties.
func newModule(router *mux.Router, mode string) *Module {
	log.Printf("New Module, mode: %s", mode)

	// init config
	config.BuildInstance(mode)

	// fetch ENV var param ?
	// maxWorker := os.Getenv("MAX_WORKERS")
	numCPUs := runtime.NumCPU()
	maxWorkers := numCPUs

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
