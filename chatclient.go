package main

import "net/rpc"
import "fmt"
import "log"
import "os"
import "bufio"
import "time"
import "common"
import "strings"
// import "container/list"

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:3410")
	if err != nil {
		log.Fatal("dialing refused")
	}
	var reply string
	username := os.Args[1]
	fmt.Printf("%s connecting to server\n", username)
	client.Call("Server.Register", username, &reply)
	if err != nil {
		log.Fatal("Register err:", err)
	}
	fmt.Printf("%s", reply)
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for {
			time.Sleep(1000 * time.Millisecond)
			var replies []string
			client.Call("Server.CheckMessages", username, &replies)
			for _, v := range replies {
				fmt.Println(v)
			}
		}
	} ()
	var arg string
	loop:
	for scanner.Scan() {
		line := scanner.Text()
		params := strings.Fields(line)
		switch params[0] {
		case "list":
			client.Call("Server.List", arg, &reply)
			fmt.Println(reply)
		case "tell":
			target := params[1]
			message := strings.Join(params[2:], " ")
			args := &common.TellArgs{
				User: username,
				Target: target,
				Message: message,
			}
			client.Call("Server.Tell", args, &reply)
		case "say":
			message := strings.Join(params[1:], " ")
			args := &common.SayArgs{
				User: username,
				Message: message,
			}
			client.Call("Server.Say", args, &reply)
		case "logout":
			client.Call("Server.Logout", username, &reply)
			break loop
		case "shutdown":
			client.Call("Server.Shutdown", username, &reply)
			break loop
		default:
			fmt.Println("help?")
		}
	}
	fmt.Println("client exiting")
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
