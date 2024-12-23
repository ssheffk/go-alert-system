package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

func EmailClient(emailAppPassword, yourMail, recipient, hostAddress, hostPort, mailSubject, mailBody string) {

	fullServerAddress := hostAddress + ":" + hostPort
	// Create email headers
	headerMap := make(map[string]string)
	headerMap["From"] = yourMail
	headerMap["To"] = recipient
	headerMap["Subject"] = mailSubject
	mailMessage := ""
	for k, v := range headerMap {
		mailMessage += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	// Add a blank line to separate headers from the body
	mailMessage += "\r\n" + mailBody

	// TLS Configuration
	tlsConfigurations := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         hostAddress,
	}

	// Establish a TLS connection
	conn, err := tls.Dial("tcp", fullServerAddress, tlsConfigurations)
	if err != nil {
		log.Panic("Error establishing TLS connection:", err)
	}

	// Create a new SMTP client
	newClient, err := smtp.NewClient(conn, hostAddress)
	if err != nil {
		log.Panic("Error creating SMTP client:", err)
	}

	// Authenticate
	authenticate := smtp.PlainAuth("", yourMail, emailAppPassword, hostAddress)
	if err = newClient.Auth(authenticate); err != nil {
		log.Panic("Authentication error:", err)
	}

	// Set the sender and recipient
	if err = newClient.Mail(yourMail); err != nil {
		log.Panic("Error setting sender:", err)
	}
	if err = newClient.Rcpt(recipient); err != nil {
		log.Panic("Error setting recipient:", err)
	}

	// Send the email data
	writer, err := newClient.Data()
	if err != nil {
		log.Panic("Error obtaining writer:", err)
	}

	_, err = writer.Write([]byte(mailMessage))
	if err != nil {
		log.Panic("Error writing message:", err)
	}

	err = writer.Close()
	if err != nil {
		log.Panic("Error closing writer:", err)
	}

	// Close the connection
	err = newClient.Quit()
	if err != nil {
		log.Println("Error closing SMTP client:", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
