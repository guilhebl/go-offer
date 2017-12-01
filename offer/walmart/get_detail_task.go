package walmart

import (
	"errors"
	"github.com/guilhebl/go-worker-pool"
	"log"
)

// Executable Task implementation for walmart - get detail
type GetDetailTask struct{}

func (t *GetDetailTask) Run(payload job.Payload) job.JobResult {
	log.Printf("Walmart GetDetail Job: %v", payload.Params)
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
