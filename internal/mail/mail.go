package mail

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/sabrodigan/bookings-app/internal/models"
)

// SendBookingConfirmation emails the guest a booking confirmation.
// If SMTP is not configured, the message is logged and nil is returned so booking still succeeds.
func SendBookingConfirmation(res models.Reservation) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	if host == "" || port == "" {
		fmt.Printf("[mail] SMTP not configured — confirmation for %s (%s to %s, room %s)\n",
			res.Email, res.StartDate, res.EndDate, res.RoomName)
		return nil
	}

	if from == "" {
		from = user
	}
	if from == "" {
		from = "bookings@localhost"
	}

	room := res.RoomName
	if room == "" {
		room = "your selected room"
	}

	subject := "Booking Confirmation"
	body := fmt.Sprintf(`Hello,

Your booking is confirmed.

Room: %s
Arrival: %s
Departure: %s
Email: %s

Thank you for booking with us.
`, room, res.StartDate, res.EndDate, res.Email)

	msg := strings.Join([]string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", res.Email),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	addr := host + ":" + port
	var auth smtp.Auth
	if user != "" && pass != "" {
		auth = smtp.PlainAuth("", user, pass, host)
	}

	return smtp.SendMail(addr, auth, from, []string{res.Email}, []byte(msg))
}
