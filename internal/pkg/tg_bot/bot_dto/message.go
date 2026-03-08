package bot_dto

type Button struct {
	Text string
}

type Message struct {
	Chat             int64
	Text             string
	ReplyToMessageID int
	Buttons          []Button
}
