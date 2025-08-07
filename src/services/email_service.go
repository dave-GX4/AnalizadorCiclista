package services

import (
	"log" // <-- AÑADE ESTE IMPORT para poder imprimir en la consola
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendConfirmationEmail(toEmail, participantName, participantCode string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpPortStr == "" || smtpUser == "" || smtpPassword == "" {
		log.Println("ADVERTENCIA: La configuración SMTP no está completa en el .env. No se puede enviar el correo.")
		return nil
	}

	smtpPort, _ := strconv.Atoi(smtpPortStr)

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "¡Registro Confirmado! Tu Código de Participante")

	body := `
		<h1>¡Hola, ` + participantName + `!</h1>
		<p>Tu registro en la Competencia Ciclista ha sido completado con éxito.</p>
		<p>Tu código de participante oficial es:</p>
		<div style="background-color: #f0f0f0; border: 1px solid #ccc; padding: 10px 20px; font-size: 24px; font-weight: bold; text-align: center; margin: 20px 0;">
			` + participantCode + `
		</div>
		<p>Este es el código que necesitarás el día del evento. ¡Guárdalo en un lugar seguro!</p>
		<p>¡Nos vemos en la carrera!</p>
	`
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		// CAMBIO 1: Imprimimos el error REAL en la consola del backend para que lo veas.
		log.Printf("ERROR CRÍTICO AL ENVIAR CORREO: %v", err)

		// CAMBIO 2: Devolvemos el error para que el handler sepa que algo falló.
		return err
	}

	log.Printf("Correo enviado exitosamente a: %s", toEmail)
	return nil
}