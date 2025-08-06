package services

import (
	"log" // <-- AÑADE ESTE IMPORT para poder imprimir en la consola
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendConfirmationEmail(toEmail, participantName, token string) error {
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
	m.SetHeader("Subject", "¡Bienvenido a la Competencia Ciclista!")

	body := `
		<h1>¡Hola, ` + participantName + `!</h1>
		<p>Tu registro ha sido completado con éxito.</p>
		<p>Este es tu token de confirmación. Guárdalo en un lugar seguro:</p>
		<pre><code>` + token + `</code></pre>
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