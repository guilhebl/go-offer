package offer

import (
	"log"
	"github.com/guilhebl/go-worker-pool"
	"errors"
)

const (
	JobTypeSearch = "SEARCH"
	JobTypeGet    = "GET"
)

// Executable Job implementation for Offer-Go, represents a model to be executed in parallel (if multiple CPUs available)
type JobRunner struct{}

func (e *JobRunner) Run(payload worker.Payload) (interface{}, error) {
	log.Printf("Run %s", payload.JobType)

	switch payload.JobType {
	case JobTypeSearch:
		{
			return SearchOffers(payload.Params), nil
		}
	case JobTypeGet:
		{
			return GetOfferDetail(payload.Params), nil
		}

	default:
		return nil, errors.New("invalid Job Type")
	}
}

func NewJobRunner() JobRunner {
	return JobRunner {}
}
