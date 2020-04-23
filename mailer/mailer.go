package mailer

type Mailer interface {
	Send(input Input) error
}

type Input struct {
	From    string
	To      string
	Subject string
	Text    string
	Html    string
}
