package tcp

import (
	"bufio"
	p "coinsystem/main/proto"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func StartClient(host string) {
	CONNECT := host
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		log.Fatal(err)
		return
	}

	if c != nil {
		for {
			request := generateRequest()

			if request == nil {
				WriteToServer([]byte("exit"), c)
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
					WriteToServer(data, c)
				}

				lineBreak()
				message := ReadFromServer(c)
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

func WriteToServer(data []byte, c net.Conn) []byte {
	if c != nil {
		data = append(data, []byte("\n")[0])
		c.Write(data)
		return data
	}
	fmt.Println("tcpC: failed to send")
	return nil
}

func ReadFromServer(c net.Conn) []byte {
	if c != nil {
		text, _ := bufio.NewReader(c).ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		return []byte(text)
	}
	fmt.Println("tspC: failed to read")
	return nil
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
