package services

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"os"
)

type EmailService struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

func NewEmailService() *EmailService {
	return &EmailService{
		Host: os.Getenv("SMTP_HOST"),
		Port: os.Getenv("SMTP_PORT"),
		User: os.Getenv("SMTP_USER"),
		Pass: os.Getenv("SMTP_PASS"),
		From: os.Getenv("SMTP_FROM"),
	}
}

func (s *EmailService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.User, s.Pass, s.Host)

	fromHeader := s.From
	envelopeFrom := s.User
	if parsed, err := mail.ParseAddress(s.From); err == nil {
		fromHeader = parsed.String()
		envelopeFrom = parsed.Address
	} else if s.From == "" {
		fromHeader = s.User
	}

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", fromHeader, to, subject, body)

	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)

	return smtp.SendMail(addr, auth, envelopeFrom, []string{to}, []byte(msg))
}

func (s *EmailService) SendVerificationCode(to, code string) error {
	subject := "ANDARIA - Código de Verificación"
	body := fmt.Sprintf(`Hola,

Tu código de verificación para ANDARIA es:

%s

Este código expirará en 10 minutos.

Si no solicitaste este código, ignora este mensaje.

Saludos,
Equipo ANDARIA`, code)

	return s.SendEmail(to, subject, body)
}

func (s *EmailService) SendPasswordResetCode(to, code string) error {
	subject := "ANDARIA - Recuperación de Contraseña"
	body := fmt.Sprintf(`Hola,

Tu código para restablecer tu contraseña en ANDARIA es:

%s

Este código expirará en 20 minutos.

Si no solicitaste restablecer tu contraseña, ignora este mensaje.

Saludos,
Equipo ANDARIA`, code)

	return s.SendEmail(to, subject, body)
}

func (s *EmailService) SendAgencyManagerWelcome(to, code, agencyName string) error {
	subject := "ANDARIA - Bienvenido como Encargado de Agencia"
	body := fmt.Sprintf(`Hola,

Has sido registrado como Encargado de Agencia en ANDARIA para la agencia: %s

Para activar tu cuenta, usa el siguiente código de verificación:

%s

Después de verificar tu correo, podrás establecer tu contraseña y acceder al sistema.

Este código expirará en 10 minutos.

Saludos,
Equipo ANDARIA`, agencyName, code)

	return s.SendEmail(to, subject, body)
}
