package emailutils

import (
  "net/smtp"
)

type EmailSender struct {
  auth smtp.Auth
  addr string
  from string
}

func (sender *EmailSender) Send(message string, recipient string) error {
  err := smtp.SendMail(sender.addr, sender.auth, sender.from,
    []string{recipient}, []byte(message))
  return err
}

type SenderConfig struct {
  ServerAddr string
  SenderName string
  Identity string
  Password string
}

func (config *SenderConfig) MakeSender() *EmailSender {
  newSender := new(EmailSender)
  newSender.from = config.SenderName
  newSender.addr = config.ServerAddr + ":smtp"
  newSender.auth = smtp.PlainAuth("", config.Identity, config.Password, config.ServerAddr)

  return newSender
}
