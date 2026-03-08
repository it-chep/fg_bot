package participant

type Participant struct {
	tgID          int64
	name          string
	username      string
	pingAvailable bool
	reportedToday bool
	reportCount   int
}

func New(options ...Option) *Participant {
	p := &Participant{}
	for _, option := range options {
		option(p)
	}
	return p
}
