package fg

import "time"

type FG struct {
	id        int64
	name      string
	chatID    int64
	adminTgID int64
	createdAt time.Time
}

func New(options ...Option) *FG {
	f := &FG{}
	for _, option := range options {
		option(f)
	}
	return f
}
