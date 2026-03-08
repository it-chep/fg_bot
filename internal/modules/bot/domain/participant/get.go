package participant

func (p *Participant) GetTgID() int64 {
	return p.tgID
}

func (p *Participant) GetName() string {
	return p.name
}

func (p *Participant) GetUsername() string {
	return p.username
}

func (p *Participant) GetPingAvailable() bool {
	return p.pingAvailable
}

func (p *Participant) GetReportedToday() bool {
	return p.reportedToday
}

func (p *Participant) GetReportCount() int {
	return p.reportCount
}
