package main

import (
	"coinsystem/main/tcp"
)

func main() {
	tcp.StartClient("127.0.0.1:1234")
}
