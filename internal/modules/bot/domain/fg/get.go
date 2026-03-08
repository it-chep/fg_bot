package fg

import "time"

func (f *FG) GetID() int64 {
	return f.id
}

func (f *FG) GetName() string {
	return f.name
}

func (f *FG) GetChatID() int64 {
	return f.chatID
}

func (f *FG) GetAdminTgID() int64 {
	return f.adminTgID
}

func (f *FG) GetCreatedAt() time.Time {
	return f.createdAt
}
