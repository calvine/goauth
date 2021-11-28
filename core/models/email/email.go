package email

// const (
// 	messageTemplate = "To:{{ .To -}}"
// )

// TODO: add HTML support
type EmailMessage struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

// func (em EmailMessage) GetMessageData() string {

// }
