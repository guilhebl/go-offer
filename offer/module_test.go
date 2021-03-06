package offer

import (
	"testing"
)

// tests if app module is built correctly setting up worker pool and other global scoped objects
func TestGetInstance(t *testing.T) {

	module := BuildInstance("test")

	if module == nil {
		t.Error("Error while creating Module")
	}

	d := module.Dispatcher

	if &d == nil || len(d.Workers) == 0 || d.WorkerPool == nil {
		t.Error("Error while creating Module Dispatcher")
	}

	j := module.JobQueue
	if &j == nil {
		t.Error("Error while creating Module Dispatcher JobQueue")
	}
}
