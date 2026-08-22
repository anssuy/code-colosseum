package judge

import "context"

type Job struct {
	Ctx       context.Context
	Language  string
	Code      string
	TestCases []TestCase
	ResultCh  chan<- Result
}

type Pool struct {
	jobs chan Job
}

func NewPool(workers int) *Pool {
	p := &Pool{
		jobs: make(chan Job, workers*2),
	}
	for range workers {
		go p.worker()
	}
	return p
}

func (p *Pool) worker() {
	for job := range p.jobs {
		result := Run(job.Ctx, job.Language, job.Code, job.TestCases)
		job.ResultCh <- result
	}
}

func (p *Pool) Submit(job Job) {
	p.jobs <- job
}
