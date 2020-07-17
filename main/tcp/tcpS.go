package tcp

import (
	"bufio"
	"coinsystem/main/database"
	p "coinsystem/main/proto"
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"strings"
)

func StartServer(port string, db *sql.DB) {
	PORT := ":" + port
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal("tcp server listener error:", err)
		return
	}

	for {
		// accept new connection
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("tcp server accept error", err)
		}

		// spawn off goroutine to able to accept new connections
		go handleConnection(conn, db)
	}
}

func handleConnection(conn net.Conn, db *sql.DB) {
	// read buffer from client after enter is hit
	data, err := ReadFromClient(conn)

	if err != nil || string(data) == "exit" {
		log.Println("client left..")
		conn.Close()
		return
	}

	request := &p.Request{}
	err2 := proto.Unmarshal(data, request)
	if err2 != nil {
		log.Fatal("Unmarshaling error: ", err)
	}

	result := execute(request, db)
	if result == nil {
		fmt.Println("-> Error in processing request")
	} else {
		fmt.Println(result.GetBalance())

		response, err := proto.Marshal(result)
		if err != nil {
			log.Fatal("Marshaling error: ", err)
		}

		WriteToClient(response, conn)
	}

	// recursive func to handle io.EOF for random disconnects
	handleConnection(conn, db)
}

func execute(request *p.Request, db *sql.DB) *p.Response {
	query := request.GetQuery()

	if query == p.Query_ADD {
		return database.Add(db, request)
	} else if query == p.Query_DEDUCT {
		return database.Deduct(db, request)
	} else if query == p.Query_BALANCE {
		return database.Balance(db, request)
	} else {
		return nil
	}
}

func WriteToClient(data []byte, c net.Conn) []byte {
	if c != nil {
		data = append(data, []byte("\n")[0])
		c.Write(data)
		return data
	}
	fmt.Println("tcpS: failed to send")
	return nil
}

func ReadFromClient(c net.Conn) ([]byte, error) {
	if c != nil {
		text, err := bufio.NewReader(c).ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		return []byte(text), err
	}
	fmt.Println("tspS: failed to read")
	return nil, nil
}
