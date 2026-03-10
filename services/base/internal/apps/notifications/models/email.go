package models

type EmailNotification struct {
	Template  string
	FromEmail string
	ToEmails  []string
	Data      map[string]string
	Subject   string
}
