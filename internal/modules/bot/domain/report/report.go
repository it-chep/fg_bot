package report

import "time"

type Report struct {
	reportMessageLink string
	reportName        string
	createdAt         time.Time
}

func New(options ...Option) *Report {
	r := &Report{}
	for _, option := range options {
		option(r)
	}
	return r
}
