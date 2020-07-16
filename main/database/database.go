package database

import (
	p "coinsystem/main/proto"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	user        int64
	coins       int64
	lastUpdated string
)

func DbSetup() *sql.DB {
	fmt.Println("Initializing database")
	fmt.Println("Please input mySQL root user password: ")
	password, _ := terminal.ReadPassword(0)

	db, err := sql.Open("mysql", "root:"+
		string(password)+"@tcp(127.0.0.1:3306)/coinsystemdb")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to MySQL")

	return db
}

func Balance(db *sql.DB, request *p.Request) *p.Response {
	userID := request.GetUserID()
	mB := make(map[int64]int64)
	mLU := make(map[int64]string)

	response := &p.Response{
		Query:       p.Query_BALANCE,
		Success:     false,
		Balance:     mB,
		LastUpdated: mLU,
	}

	rows, err := db.Query("SELECT * FROM coinsDB WHERE userID IN " + sqlIntConv(userID))

	if err == nil {
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&user, &coins, &lastUpdated)
			if err != nil {
				fmt.Println(err)
			}
			mB[user] = coins
			mLU[user] = lastUpdated

			fmt.Println(user)
			fmt.Println(coins)
			fmt.Println(lastUpdated)
		}

		err = rows.Err()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}

	response = &p.Response{
		Query:       p.Query_BALANCE,
		Success:     true,
		Balance:     mB,
		LastUpdated: mLU,
	}

	return response
}

func Add(db *sql.DB, request *p.Request) *p.Response {
	user = request.GetUserID()[0]
	before := SingleBalance(db, user)
	after := before + request.GetCoins()

	mB := make(map[int64]int64)
	mLU := make(map[int64]string)

	response := &p.Response{
		Query:       p.Query_ADD,
		Success:     false,
		Balance:     mB,
		LastUpdated: mLU,
	}

	if (before != -1) && (after <= 10000) {
		curTime := time.Now()
		results, err :=
			db.Query("UPDATE coinsDB SET balance = ?, last_updated = ? WHERE userID = ?",
				after, curTime, user)

		if err != nil {
			fmt.Println(err)
		} else {
			mB[user] = after
			mLU[user] = curTime.Format(time.RFC822Z)

			response = &p.Response{
				Query:       p.Query_ADD,
				Success:     true,
				Balance:     mB,
				LastUpdated: mLU,
			}
		}

		defer results.Close()
	}

	return response
}

func Deduct(db *sql.DB, request *p.Request) *p.Response {
	user = request.GetUserID()[0]
	before := SingleBalance(db, user)
	after := before - request.GetCoins()

	mB := make(map[int64]int64)
	mLU := make(map[int64]string)

	response := &p.Response{
		Query:       p.Query_DEDUCT,
		Success:     false,
		Balance:     mB,
		LastUpdated: mLU,
	}

	if (before != -1) && (after >= 0) {
		curTime := time.Now()
		results, err :=
			db.Query("UPDATE coinsDB SET balance = ?, last_updated = ? WHERE userID = ?",
				after, curTime, user)

		if err != nil {
			fmt.Println(err)
		} else {
			mB[user] = after
			mLU[user] = curTime.Format(time.RFC3339)

			response = &p.Response{
				Query:       p.Query_DEDUCT,
				Success:     true,
				Balance:     mB,
				LastUpdated: mLU,
			}
		}

		defer results.Close()
	}

	return response
}

func SingleBalance(db *sql.DB, userID int64) int64 {
	err := db.QueryRow("SELECT balance FROM coinsDB WHERE userID = ?", userID).Scan(&coins)
	if err != nil {
		return -1
	}
	return coins
}

func sqlIntConv(array []int64) string {
	result := strings.Builder{}
	result.WriteString("(")

	if len(array) != 0 {
		for i := 0; i < len(array); i++ {
			result.WriteString("'")
			result.WriteString(strconv.FormatInt(array[i], 10))
			result.WriteString("'")

			if i != len(array)-1 {
				result.WriteString(",")
			}

		}
	}

	result.WriteString(")")
	return result.String()
}

func main() {
	db := DbSetup()

	request := &p.Request{
		Query:  p.Query_BALANCE,
		UserID: []int64{123456, 197639, 678532, 557094},
		Coins:  -1,
	}

	result := Balance(db, request)

	fmt.Println("Result: ", result.GetBalance())
}
