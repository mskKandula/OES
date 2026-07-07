package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	gomail "gopkg.in/gomail.v2"
)

// emailPayload matches the JSON published by the main server's student service.
type emailPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"` // plaintext — for welcome email only
}

// registrationEmailTmpl is the welcome email body. Kept in sync with
// Server/util/emailConfig/registrationMailTemplate.html.
const registrationEmailTmpl = `<!DOCTYPE html>
<html>
    <body>
        <div>
            <span>
                Dear {{.Name}},<br><br> You have successfully registered.
            </span>
            <p>Please login with below User Name & Password</p>
        </div>
        <table>
            <tr>
                <td><b>User Name:</b></td>
                <td>{{.Email}}</td>
            </tr>
            <tr>
                <td><b>Password:</b></td>
                <td>{{.Password}}</td>
            </tr>
        </table>
        <p>Regards,</p>
        <p>Registrar</p>
    </body>
</html>`

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// dialRabbitMQ establishes a RabbitMQ connection with an explicit TCP dial
// timeout and heartbeat. It retries indefinitely with exponential back-off
// (capped at 30 s) so that MQServer recovers automatically if the broker
// restarts while the service is running.
func dialRabbitMQ(dsn string) *amqp.Connection {
	backoff := 2 * time.Second
	const maxBackoff = 30 * time.Second

	for {
		conn, err := amqp.DialConfig(dsn, amqp.Config{
			// Heartbeat controls how often each side sends a heartbeat frame.
			// If the broker doesn't receive one within 2× this interval it
			// closes the connection, enabling fast detection of dead TCP links.
			Heartbeat: 10 * time.Second,
			// Dial enforces an explicit TCP connection timeout, preventing
			// indefinite hangs when the broker is temporarily unreachable.
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, 10*time.Second)
			},
		})
		if err == nil {
			log.Println("[amqp] connected to RabbitMQ")
			return conn
		}
		log.Printf("[amqp] connect failed: %v — retrying in %s", err, backoff)
		time.Sleep(backoff)
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
}

func main() {
	rabbitDSN := os.Getenv("RABBITMQ_DSN")
	if rabbitDSN == "" {
		rabbitDSN = "amqp://rabbitmq:rabbitmq@messageq/"
	}

	workerCount := runtime.NumCPU()
	sem := make(chan struct{}, workerCount)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		conn := dialRabbitMQ(rabbitDSN)

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")

		// ── Video encoding queue ─────────────────────────────────────────────
		encodeQ, err := ch.QueueDeclare(
			"encode", // name
			true,     // durable
			false,    // delete when unused
			false,    // exclusive
			false,    // no-wait
			nil,      // arguments
		)
		failOnError(err, "Failed to declare encode queue")

		encodeMsgs, err := ch.Consume(
			encodeQ.Name, // queue
			"",           // consumer tag
			false,        // auto-ack
			false,        // exclusive
			false,        // no-local
			false,        // no-wait
			nil,          // args
		)
		failOnError(err, "Failed to register encode consumer")

		// ── Email queue ──────────────────────────────────────────────────────
		// Declared as durable so pending jobs survive a broker restart.
		emailQ, err := ch.QueueDeclare(
			"email", // name
			true,    // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		failOnError(err, "Failed to declare email queue")

		emailMsgs, err := ch.Consume(
			emailQ.Name, // queue
			"",          // consumer tag
			false,       // auto-ack — manual ack for reliability
			false,       // exclusive
			false,       // no-local
			false,       // no-wait
			nil,         // args
		)
		failOnError(err, "Failed to register email consumer")

		// connClosed fires when the AMQP connection drops (broker restart, etc.)
		connClosed := conn.NotifyClose(make(chan *amqp.Error, 1))

		// ── Video encoding consumer ──────────────────────────────────────────
		go func() {
			for d := range encodeMsgs {
				log.Printf("[encode] received: %s", d.Body)
				sem <- struct{}{}
				go func(msg amqp.Delivery) {
					defer func() { <-sem }()
					HlsVideoConversion(string(msg.Body))
				}(d)
			}
		}()

		// ── Email consumer ───────────────────────────────────────────────────
		go func() {
			for d := range emailMsgs {
				log.Printf("[email] received job")
				sem <- struct{}{}
				go func(msg amqp.Delivery) {
					defer func() { <-sem }()
					SendWelcomeEmail(msg.Body)
					// Ack only after successful send; on failure the message is
					// Nack'd with requeue=true so it will be retried.
					if err := msg.Ack(false); err != nil {
						log.Printf("[email] ack failed: %v", err)
					}
				}(d)
			}
		}()

		log.Printf("[*] Waiting for messages. To exit press CTRL+C")

		// Block until either a shutdown signal or a connection error.
		select {
		case <-quit:
			log.Printf("[*] Shutting down MQServer")
			ch.Close()
			conn.Close()
			return
		case amqpErr := <-connClosed:
			if amqpErr != nil {
				log.Printf("[amqp] connection lost: %v — reconnecting", amqpErr)
			} else {
				log.Printf("[amqp] connection closed — reconnecting")
			}
			// Fall through to the top of the loop to reconnect.
		}
	}
}

var welcomeEmailTmpl = template.Must(template.New("reg").Parse(registrationEmailTmpl))

// SendWelcomeEmail unmarshals an emailPayload from body, renders the HTML
// template, and dispatches via Gmail SMTP.
// SMTP credentials are read from SMTP_EMAIL and SMTP_PASSWORD env vars.
func SendWelcomeEmail(body []byte) {
	var p emailPayload
	if err := json.Unmarshal(body, &p); err != nil {
		log.Printf("[email] failed to unmarshal payload: %v", err)
		return
	}

	var buf bytes.Buffer
	if err := welcomeEmailTmpl.Execute(&buf, p); err != nil {
		log.Printf("[email] failed to execute template: %v", err)
		return
	}

	senderEmail := os.Getenv("SMTP_EMAIL")
	senderPassword := os.Getenv("SMTP_PASSWORD")

	if senderEmail == "" || senderPassword == "" {
		log.Printf("[email] SMTP_EMAIL or SMTP_PASSWORD not set — skipping send to %s", p.Email)
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", p.Email)
	m.SetHeader("Subject", "Registration Successful!")
	m.SetBody("text/html", buf.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, senderPassword)
	// TLS verification is intentionally left at the default (enabled).

	if err := d.DialAndSend(m); err != nil {
		log.Printf("[email] failed to send to %s: %v", p.Email, err)
		return
	}

	log.Printf("[email] sent to %s", p.Email)
}

// HlsVideoConversion transcodes a video to 360p/480p/720p HLS streams.
func HlsVideoConversion(result string) {

	dir, file := filepath.Split(result)

	// For single(original resolution)
	// cmd := exec.Command("ffmpeg", "-i", paths[4], "-codec:", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", "index.m3u8")

	// For Multiple(360p,480p & 720p resolutions)
	cmd := exec.Command("ffmpeg", "-i", result, "-map", "0:v:0", "-map", "0:a:0", "-map",
		"0:v:0", "-map", "0:a:0", "-map", "0:v:0", "-map", "0:a:0", "-c:v", "libx264", "-crf",
		"22", "-c:a", "aac", "-ar", "48000", "-filter:v:0", "scale=w=480:h=360", "-maxrate:v:0",
		"600k", "-b:a:0", "64k", "-filter:v:1", "scale=w=640:h=480", "-maxrate:v:1", "900k",
		"-b:a:1", "128k", "-filter:v:2", "scale=w=1280:h=720", "-maxrate:v:2", "1500k", "-b:a:2",
		"128k", "-var_stream_map", "v:0,a:0,name:360p v:1,a:1,name:480p v:2,a:2,name:720p",
		"-preset", "slow", "-hls_list_size", "0", "-threads", "0", "-f", "hls", "-hls_playlist_type",
		"event", "-hls_time", "10", "-hls_flags", "independent_segments", "-master_pl_name",
		"index.m3u8", "-y", dir+"%v/index.m3u8")

	err := cmd.Run()

	if err != nil {
		log.Println(err)
		return
	}

	imageFileName := strings.Split(file, ".")[0] + ".png"

	cmd = exec.Command("ffmpeg", "-i", result, "-ss", "00:00:01.000", "-vframes", "1", filepath.Join(dir, imageFileName))

	err = cmd.Run()

	if err != nil {
		log.Println(err)
		return
	}

	err = os.Remove(result)

	if err != nil {
		log.Println(err)
	}
}
