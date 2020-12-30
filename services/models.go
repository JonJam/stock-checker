package services

type SmsClient interface {
	Send(body string) error
}
