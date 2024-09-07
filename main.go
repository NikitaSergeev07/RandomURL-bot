package main

import (
	tgClient "RandomURL/clients/telegram"
	event_consumer "RandomURL/consumer/event-consumer"
	"RandomURL/events/telegram"
	"RandomURL/storage/redis"
	"flag"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
	redisAddr = "localhost:6379"
	batchSize = 100
)


func main() {
	eventsProcessor := telegram.New(tgClient.New(
		tgBotHost, mustToken()),
		redis.NewRedisStorage(redisAddr, "", 0),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()
	if *token == "" {
		log.Fatal("token is required")
	}
	return *token
}
