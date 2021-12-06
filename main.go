package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"test-client/proofofwork"
	"time"
)

var timeEOF int64 = 0

func main() {
	for {
		run()
	}
}

func run() {
	conn, err := net.Dial("tcp", os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"))
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected (" + os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT") + ")")

	decoder := json.NewDecoder(conn)
	newChallengeResponse(conn)
	timeEOF = time.Now().Unix()

	for {
		receive := []string{}
		conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		err := decoder.Decode(&receive)

		if err == io.EOF {
			if checkDisconnect() {
				fmt.Println("Disconnected")
				break
			}
			continue
		}
		if err != nil {
			fmt.Println("Error json: ", err)
			newChallengeResponse(conn)
			continue
		}
		if len(receive) == 0 {
			fmt.Println("message not found")
			continue
		}

		switch receive[0] {
		case "SOLVE":
			solve(conn, receive[2])
		case "GRANT":
			grant(conn, receive)
		default:
			fmt.Println("type message not found")
			continue
		}
	}
}

func send(conn net.Conn, message interface{}) {
	fmt.Println(">> ", message)
	encoder := json.NewEncoder(conn)
	err := encoder.Encode(&message)
	if err != nil {
		fmt.Println("Encode error = ", err)
		return
	}
}

func newChallengeResponse(conn net.Conn) {
	t := time.NewTimer(3 * time.Second)
	<-t.C
	fmt.Println("===========  Started new Challenge Response  ===========")
	choose(conn)
}

func choose(conn net.Conn) {
	message := [...]string{
		"CHOOSE",
		os.Getenv("NAME"),
	}
	send(conn, message)
}

func solve(conn net.Conn, hashCash string) {
	cash := strings.Split(hashCash, ":")
	if len(cash) < 2 {
		fmt.Println("bits not found")
		return
	}
	targetBits, err := strconv.Atoi(cash[1])
	if err != nil {
		fmt.Println("bits not found")
		return
	}

	pow := &proofofwork.ProofOfWork{
		HashCash:   hashCash,
		TargetBits: targetBits,
	}
	nonce, sha1Hash := pow.Search()

	message := [...]string{
		"VERIFY",
		os.Getenv("NAME"),
		strconv.Itoa(nonce),
		base64.StdEncoding.EncodeToString(sha1Hash),
	}
	send(conn, message)
}

func grant(conn net.Conn, receive []string) {
	fmt.Println("Word of Wisdom: \n", receive[2])
	fmt.Println("======================================================")
	newChallengeResponse(conn)
}

func checkDisconnect() bool {
	currentTime := time.Now().Unix()
	if currentTime-timeEOF < int64(time.Second) {
		return true
	}
	timeEOF = currentTime
	return false
}
