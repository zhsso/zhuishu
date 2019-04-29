package main

func main() {
	teleBot := &TeleBot{}
	userManager := &UserManager{}
	bookManager := &BookManager{}

	teleBot.run()
	userManager.run()
	bookManager.run()
}
