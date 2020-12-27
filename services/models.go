package services

type SmsClient interface {
	send(body string) error
}
