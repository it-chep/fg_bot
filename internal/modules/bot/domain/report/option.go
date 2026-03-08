package report

import "time"

type Option func(r *Report)

func WithReportMessageLink(link string) Option {
	return func(r *Report) {
		r.reportMessageLink = link
	}
}

func WithReportName(name string) Option {
	return func(r *Report) {
		r.reportName = name
	}
}

func WithCreatedAt(t time.Time) Option {
	return func(r *Report) {
		r.createdAt = t
	}
}
