package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"encode", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			HlsVideoConversion(string(d.Body))
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

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
