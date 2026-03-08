package participant

type Option func(p *Participant)

func WithTgID(tgID int64) Option {
	return func(p *Participant) {
		p.tgID = tgID
	}
}

func WithName(name string) Option {
	return func(p *Participant) {
		p.name = name
	}
}

func WithUsername(username string) Option {
	return func(p *Participant) {
		p.username = username
	}
}

func WithPingAvailable(enabled bool) Option {
	return func(p *Participant) {
		p.pingAvailable = enabled
	}
}

func WithReportedToday(reported bool) Option {
	return func(p *Participant) {
		p.reportedToday = reported
	}
}

func WithReportCount(count int) Option {
	return func(p *Participant) {
		p.reportCount = count
	}
}
