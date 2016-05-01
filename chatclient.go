package main

import "net/rpc"
import "fmt"
import "log"
import "os"
import "bufio"
import "time"
// import "container/list"

func checkMessages() {
	for {

	}
}
func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:3410")
	if err != nil {
		log.Fatal("dialing refused")
	}
	var reply string
	username := os.Args[1]
	fmt.Printf("%s connecting to server\n", username)
	client.Call("Server.Register", &username, &reply)
	if err != nil {
		log.Fatal("Register err:", err)
	}
	fmt.Printf("%s", reply)
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for {
			time.Sleep(1000 * time.Millisecond)
			var replies []string
			client.Call("Server.CheckMessages", username, replies)
			for _, v := replies {
				fmt.Println(v)
			}
		}
	}
	var command string
	for scanner.Scan() {
		switch scanner.Text() {
		case "list":
			client.Call("Server.List", &reply)
			fmt.Println(reply)
		case "tell":
			fmt.Println("to whom?")
			target := scanner.Text()
			fmt.Println("what?")
			message := scanner.Text()
			args := &TellArgs{
				User: username,
				Target: target,
				Message: message
			}
			client.Call("Server.Tell", args, &reply)
		case "say":
			fmt.Println("what?")
			message := scanner.Text()
			args := &TellArgs{
				User: username,
				Message: message
			}
			client.Call("Server.Say", args, &reply)
		case "logout":
			client.Call("Server.Logout", username)
			break
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
