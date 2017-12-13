package bestbuy

import (
	"errors"
	"github.com/guilhebl/go-worker-pool"
)

// Executable Task implementation for search
type SearchTask struct{}

func (t *SearchTask) Run(payload job.Payload) job.JobResult {
	r := search(payload.Params)
	if r == nil {
		return job.NewJobResult(nil, errors.New("error on search"))
	}

	return job.NewJobResult(r, nil)
}

func NewSearchTask() SearchTask {
	return SearchTask{}
}
