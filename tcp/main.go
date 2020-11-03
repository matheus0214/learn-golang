package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var db = map[string]string{}

type Result struct {
	text []byte
}

func (r *Result) Write(p []byte) (int, error) {
	r.text = p

	return len(p), nil
}

func main() {
	db["d1"] = "Datas 1"
	db["d2"] = "Datas 2"

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn)
	}
}

// Handle connection reciver
// write -> /write/key/datas
// read -> /read/key
func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	scanner.Scan()

	url := scanner.Text()
	options := strings.Split(url, " ")
	splited := strings.Split(options[1], "/")

	conn.Write([]byte("HTTP:1.1 200 OK\r\n"))
	conn.Write([]byte("Content-Type: text/html\r\n"))

	switch splited[1] {
	case "":
                var teste Result

		for key, value := range db {
                        teste.Write([]byte(fmt.Sprintf("%s - %s", key, value)))
                        // fmt.Fprintf(&teste, "%s %s", key, value)
			// resultText += key + " - " + value + "\n"
                }


                special := fmt.Sprint("Content-Length: ", len(teste.text), string("\r\n\r\n"))

		conn.Write([]byte(special))
                fmt.Println(string(teste.text))

		conn.Write(teste.text)
	case "read":
		result := fmt.Sprint(db[splited[2]], "\r\n")
		conn.Write([]byte(fmt.Sprint("Content-Length: ", len(result), "\r\n\r\n")))
		conn.Write([]byte(result))
	case "write":
		db[splited[2]] = splited[3]
		message := "Datas created success"
		conn.Write([]byte(fmt.Sprint("Content-Length: ", len(message), "\r\n\r\n")))
		conn.Write([]byte(message))
	}

	conn.Close()
}
