// -----------------------------------------------------
// REDIS-X
// © Utsav Virani
// Written by: (Utsav Virani)
// -----------------------------------------------------

package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Listen for TCP connections on port 6379
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Accept incoming connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}
		if len(value.array) == 0 {
			fmt.Println("Invalid request, empty array")
			continue
		}

		cmd := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writter := NewWritter(conn)
		handler, ok := Handlers[cmd]
		if !ok {
			fmt.Printf("Unknown command: %v\n", cmd)
			writter.Write(Value{typ: "string", str: ""})
			continue
		}
		result := handler(args)
		writter.Write(result)
	}
}
