package main

import (
	"log"
	"net"
	"os"
)

type processDefinition struct {
	Name      string   `json:"name"`      // "name": "my-api",
	Cmd       string   `json:"cmd"`       // "cmd": "./my-api",
	Args      []string `json:"args"`      // "args": ["--port", "8080"],
	Cwd       string   `json:"cwd"`       // "cwd": "/home/user/projects/my-api",
	Instances int      `json:"instances"` // "instances": 2,
}

func handleConnection(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		data := buf[0:n]
		dataAsString := string(data) + " from server"

		_, err = conn.Write([]byte(dataAsString))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	socketPath := "/tmp/echo.sock"
	os.Remove(socketPath)

	l, err := net.Listen("unix", "/tmp/echo.sock")
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}

		go handleConnection(conn)
	}
}
