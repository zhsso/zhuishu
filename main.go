package main
import "os"
import "fmt"

import "github.com/go-redis/redis"

var teleBot *TeleBot
var userManager *UserManager
var bookManager *BookManager
var redisClient *redis.Client

func main() {
	redisAddr :=os.Getenv("REDIS_PORT_6379_TCP_ADDR")
	redisPort := os.Getenv("REDIS_PORT_6379_TCP_PORT")
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisAddr, redisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	teleBot = &TeleBot{}
	userManager = newUserManager()
	bookManager = newBookManager()

	go teleBot.run()

	bookManager.loadAll()
	userManager.loadAll()

	bookManager.run()
	stop := make(chan bool)
	<-stop
}
