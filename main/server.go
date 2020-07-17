package main

import (
	"coinsystem/main/database"
	"coinsystem/main/tcp"
)

func main() {
	db := database.DbSetup()
	tcp.StartServer("1234", db)
}
