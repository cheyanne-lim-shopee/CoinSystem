package main

import (
	"coinsystem/main/database"
	p "coinsystem/main/proto"
	"coinsystem/main/tcp"
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
)

func main() {
	db := database.DbSetup()
	c := tcp.ServerSetup("1234")

	for {
		data := tcp.ReadFromClient(c)

		if string(data) == "exit" {
			fmt.Println("TCP server exiting...")
			return
		}

		request := &p.Request{}
		err := proto.Unmarshal(data, request)
		if err != nil {
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

			tcp.WriteToClient(response, c)
		}
	}
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
