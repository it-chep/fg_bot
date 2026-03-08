package fg

import "time"

type Option func(f *FG)

func WithID(id int64) Option {
	return func(f *FG) {
		f.id = id
	}
}

func WithName(name string) Option {
	return func(f *FG) {
		f.name = name
	}
}

func WithChatID(chatID int64) Option {
	return func(f *FG) {
		f.chatID = chatID
	}
}

func WithAdminTgID(adminTgID int64) Option {
	return func(f *FG) {
		f.adminTgID = adminTgID
	}
}

func WithCreatedAt(t time.Time) Option {
	return func(f *FG) {
		f.createdAt = t
	}
}
