package main

import "net/rpc"
import "log"
import "errors"
import "net"
import "net/http"
import "fmt"
import "common"

type Server struct {
	users map[string]chan string
	name  string
	shutChan chan bool
}

func (s *Server) Register(username string, reply *string) error {
	if username == "" {
		return errors.New("username required")
	}
	if _, exists := s.users[username]; exists {
		return errors.New("username already registered")
	}
	msg := username +" joined"
	for user, msgQueue := range s.users {
		*reply += user + "\n"
		msgQueue <- msg
	}
	s.users[username] = make(chan string, 100)
	return nil
}
func (s *Server) List(arg string, reply *string) error{
	for user, _ := range s.users {
		*reply += user + "\n"
	}
	return nil
}
func (s *Server) CheckMessages(username string, reply *[]string) error{
	fmt.Printf("%s CheckMessages\n", username)
	*reply = make([]string, len(s.users[username]))
	
	for i:=0;len(s.users[username]) > 0;i++ {
		(*reply)[i] = <- s.users[username]
	}
	return nil
}
func (s *Server) Tell(arg *common.TellArgs, reply *string) error{
	fmt.Printf("%s telling to %s: %s\n", arg.User, arg.Target, arg.Message)
	msg := arg.User + " tells you " + arg.Message
	s.users[arg.Target] <- msg
	return nil
}
func (s *Server) Say(arg *common.SayArgs, reply *string) error{
	fmt.Printf("%s says %s\n", arg.User, arg.Message)
	msg := arg.User + " says " + arg.Message
	for _, msgQueue := range s.users {
		msgQueue <- msg
	}
	return nil
}
func (s *Server) Logout(username string, reply *string) error{
	fmt.Printf("%s logout\n", username)
	delete(s.users, username)
	return nil
}
func (s *Server) Shutdown(username string, reply *string) error {
	fmt.Printf("%s request shutdown\n", username)
	s.shutChan <- true
	return nil
}
func main() {
	server := new(Server)
	server.users = make(map[string] chan string)
	server.shutChan = make(chan bool)
	rpc.Register(server)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":3410")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Println("serving")
	go http.Serve(l, nil)
	<- server.shutChan
	fmt.Println("shutting")
}
