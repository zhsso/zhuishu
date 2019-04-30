package main

import "github.com/go-redis/redis"

var teleBot *TeleBot
var userManager *UserManager
var bookManager *BookManager
var redisClient *redis.Client

func main() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	teleBot = &TeleBot{}
	userManager = &UserManager{}
	bookManager = &BookManager{}

	go teleBot.run()
	bookManager.run()
}
