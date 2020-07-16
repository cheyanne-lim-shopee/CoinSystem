package main

import (
	"bufio"
	p "coinsystem/main/proto"
	"coinsystem/main/tcp"
	"fmt"
	"github.com/golang/protobuf/proto"
	"os"
	"strconv"
	"strings"
)

func main() {
	c := tcp.ClientSetup("127.0.0.1:1234")

	if c != nil {
		for {
			request := generateRequest()

			if request == nil {
				tcp.WriteToServer([]byte("exit"), c)
				fmt.Println("TCP client exiting...")
				lineBreak()
				return
			} else {
				data, err := proto.Marshal(request)

				if err != nil {
					fmt.Println("Marshaling error: ", err)
					lineBreak()
					continue
				} else {
					tcp.WriteToServer(data, c)
				}

				lineBreak()
				message := tcp.ReadFromServer(c)
				response := &p.Response{}
				err = proto.Unmarshal(message, response)

				if err != nil {
					fmt.Println("Unmarshaling error: ", err)
					lineBreak()
					continue
				} else {
					if response.GetSuccess() {
						generateResponse(response)
					} else {
						fmt.Print("->: Error in processing request \n")
						lineBreak()
					}
				}
			}
		}
	}
}

func generateResponse(response *p.Response) {
	if response.GetQuery() == p.Query_BALANCE {
		if len(response.GetBalance()) == 0 {
			fmt.Println("No valid user ids given")
		} else {
			fmt.Println("Showing balances of valid user ids")
		}
	}

	for i := range response.GetBalance() {
		fmt.Println("User", i, ": Balance", response.GetBalance()[i],
			"coins, Last Updated", response.GetLastUpdated()[i])
	}

	lineBreak()
}

func generateRequest() *p.Request {
	query := queryPrompt()

	if query == p.Query_QUIT {
		return nil
	}

	userID := userIDPrompt()
	users := append(make([]int64, 0), userID)
	coins := int64(-1)

	if query == p.Query_BALANCE {
		for userID != -1 {
			fmt.Println("Add more users (Input 'STOP' to end)")
			userID = userIDPrompt()
			users = append(users, userID)
		}
	} else {
		coins = coinsPrompt()
	}

	request := &p.Request{
		Query:  query,
		UserID: users,
		Coins:  coins,
	}

	return request
}

func queryPrompt() p.Query {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What do you want to do?")

	for {
		fmt.Println("1) Add coins to user")
		fmt.Println("2) Deduct coins from user")
		fmt.Println("3) Check coin balance of user")
		fmt.Println("4) Quit")

		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		value, err := strconv.Atoi(text)

		if err == nil && 1 <= value && value <= 4 {
			if value == 1 {
				return p.Query_ADD
			} else if value == 2 {
				return p.Query_DEDUCT
			} else if value == 3 {
				return p.Query_BALANCE
			} else {
				return p.Query_QUIT
			}
		}
		fmt.Println("Please input an integer between 1 and 4.")
	}
}

func userIDPrompt() int64 {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input requested userID")

	for {
		text, _ := reader.ReadString('\n')
		if text == "STOP\n" {
			return -1
		} else {
			text = strings.TrimSuffix(text, "\n")
			value, err := strconv.ParseInt(text, 10, 64)
			if err == nil && 1 <= value {
				return value
			}
		}
		fmt.Println("Please input a valid input")
	}
}

func coinsPrompt() int64 {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input number of coins")

	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		value, err := strconv.ParseInt(text, 10, 64)
		if err == nil && 1 <= value {
			return value
		}
		fmt.Println("Please input a positive integer")
	}
}

func lineBreak() {
	fmt.Println("----------------------------------------------------")
}

//func main() {
//	c := tcp.ClientSetup("127.0.0.1:1234")
//
//	if c != nil {
//		for {
//			reader := bufio.NewReader(os.Stdin)
//			fmt.Print(">> ")
//			text, _ := reader.ReadString('\n')
//			tcp.WriteToServer([]byte(text), c)
//
//			message := tcp.ReadFromServer(c)
//			fmt.Print("->: " + message)
//
//			if strings.TrimSpace(text) == "STOP" {
//				fmt.Println("TCP client exiting...")
//				return
//			}
//		}
//	}
//}

//func main() {
//	c := tcp.ClientSetup("127.0.0.1:1234")
//
//	if c != nil {
//		for {
//			request := generateRequest()
//
//			if request == nil {
//				fmt.Println("TCP client exiting...")
//				return
//			} else {
//				data, err := proto.Marshal(request)
//
//				if err != nil {
//					fmt.Println("marshaling error: ", err)
//					continue
//				} else {
//					check := tcp.WriteToServer(data, c)
//
//					// For testing only
//					newRequest := &p.Request{}
//					err = proto.Unmarshal(data, newRequest)
//					if err != nil {
//						log.Fatal("unmarshaling error: ", err)
//					}
//
//					fmt.Println("check: ", check)
//					fmt.Println("q: ", newRequest.GetQuery())
//					fmt.Println("id: ", newRequest.GetUserID())
//					fmt.Println("c: ", newRequest.GetCoins())
//				}
//
//				message := tcp.ReadFromServer(c)
//				fmt.Print("->: " + message)
//			}
//		}
//	}
//}
