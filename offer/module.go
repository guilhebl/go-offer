package offer

import (
	"github.com/gorilla/mux"
	"github.com/guilhebl/go-offer/common/cache"
	"github.com/guilhebl/go-offer/common/config"
	"github.com/guilhebl/go-worker-pool"
	"log"
	"runtime"
	"sync"
	"github.com/guilhebl/go-offer/common/db"
	"net/http"
)

// centralized module manager which holds references to JobQueue and other global app scoped objects
// Singleton enforcing the module will be initialized at max. once per app.
type Module struct {
	JobQueue   chan job.Job
	Dispatcher *job.WorkerPool
	Router     *mux.Router
	RedisCache *cache.RedisCache
	CassandraClient *db.CassandraClient
}

var instance *Module
var once sync.Once

func BuildInstance(mode string) *Module {
	once.Do(func() {
		instance = newModule(mode)
	})
	return instance
}

func GetInstance() *Module {
	return instance
}

// Builds a new module which is a container for the running app instance
// router - the router configuration with URL routes and mapped action handlers
// mode - test or production modes, which will make the app read from either test or prod config properties.
func newModule(mode string) *Module {
	log.Printf("New Module, mode: %s", mode)

	// init config
	config.BuildInstance(mode)

	// init mux
	router := NewRouter()
	// init static folder
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	// init Db
	clusterConfig := db.BuildInstance(config.GetProperty("cassandraHost"),
		config.GetProperty("cassandraUser"),
		config.GetProperty("cassandraPassword"),
		config.GetProperty("cassandraKeyspace"),
		config.GetIntProperty("cassandraPort"),)

	// fetch ENV var param ?
	// maxWorker := os.Getenv("MAX_WORKERS")
	numCPUs := runtime.NumCPU()
	maxWorkers := numCPUs
	workerPool := job.NewWorkerPool(maxWorkers)
	jobQueue := make(chan job.Job)

	module := Module{
		Dispatcher: &workerPool,
		JobQueue:   jobQueue,
		Router:     router,
		CassandraClient: clusterConfig,
	}

	// A buffered channel that we can send work requests on.
	module.Dispatcher.Run(jobQueue)

	// init cache
	if config.GetBoolProperty("cacheEnabled") {
		host := config.GetProperty("cacheHost")
		port := config.GetProperty("cachePort")
		cacheDefaultExpiration := config.GetIntProperty("cacheExpirationSeconds")
		module.RedisCache = cache.BuildInstance(host, port, cacheDefaultExpiration)
	}

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
