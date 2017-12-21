package amazon

import (
	"errors"
	"github.com/guilhebl/go-worker-pool"
)

// Executable Task implementation for get detail
type GetDetailTask struct{}

func (t *GetDetailTask) Run(payload job.Payload) job.JobResult {
	m := payload.Params
	r := GetOfferDetail(m["id"], m["idType"], m["country"])
	if r == nil {
		return job.NewJobResult(nil, errors.New("error on search"))
	}

	return job.NewJobResult(r, nil)
}

func NewGetDetailTask() GetDetailTask {
	return GetDetailTask{}
}
