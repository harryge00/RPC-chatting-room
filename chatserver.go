package main

import "net/rpc"
import "container/list"
import "log"
import "errors"
import "net"
import "net/http"
import "fmt"

type Server struct {
	users *list.List
	name  string
}

func (s *Server) Register(user *string, reply *string) error {
	if user == nil {
		return errors.New("username required")
	}
	log.Printf("%s joined", *user)
	if s.users.Len() == 0 {
		*reply = ""
	}
	for e := s.users.Front(); e != nil; e = e.Next() {
		// fmt.Println(e.Value.(string))
		*reply += e.Value.(string) + "\n"
	}
	s.users.PushBack(*user)
	return nil
}
func (s *Server) List(reply *string) error{
	for e := s.users.Front(); e != nil; e = e.Next() {
			fmt.Println(e.Value.(string))
			*reply += e.Value.(string) + "\n"
	}
	return nil
}
func (s *Server) CheckMessages(username string, reply []string) error{
	fmt.Printf("%s CheckMessages\n", username)
	reply = make([]string, 2)
	reply[0] = "Not"
	reply[1] = "Implemented"
	return nil
}
func (s *Server) Tell(arg *TellArgs, reply *string) error{
	fmt.Printf("%s telling to %s: %s\n", arg.User, arg.Target, arg.Message)
	*reply = "Not Implemented"
	return nil
}
func (s *Server) Say(arg *SayArgs, reply *string) error{
	fmt.Printf("%s says %s\n", arg.User, arg.Message)
	*reply = "Not Implemented"
	return nil
}
func (s *Server) Logout(username string) error{
	fmt.Printf("%s logout\n", username)
	return nil
}
func main() {
	server := new(Server)
	server.users = list.New()
	rpc.Register(server)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":3410")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Println("serving")
	go http.Serve(l, nil)
	for{
		select {

		}
	}
}
