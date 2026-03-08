package report

import "time"

func (r *Report) GetReportMessageLink() string {
	return r.reportMessageLink
}

func (r *Report) GetReportName() string {
	return r.reportName
}

func (r *Report) GetCreatedAt() time.Time {
	return r.createdAt
}
