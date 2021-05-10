// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

//!+broadcaster
type client chan<- string // an outgoing message channel

type userStruct struct {
	username string
	msg      string
}

type clientStruct struct {
	client     client
	connection net.Conn
}

var (
	entering    = make(chan clientStruct)
	leaving     = make(chan client)
	messages    = make(chan string) // all incoming client messages
	serverName  = "irc-server > "
	admin       client
	adminUser   string
	direct      = make(chan userStruct)
	kick        = make(chan string)
	users       map[string]client
	connections map[string]net.Conn
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	users = make(map[string]client)
	connections = make(map[string]net.Conn)

	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}

		case usersStr := <-direct:
			users[usersStr.username] <- usersStr.msg
			fmt.Println(serverName + "/msg" + usersStr.username + "message to " + usersStr.msg)

		case clientStr := <-entering:
			if len(clients) == 0 {
				clientStr.client <- serverName + "Congrats, you were the first user."
				clientStr.client <- serverName + "You're the new IRC Server ADMIN"
				fmt.Printf("irc-server >[%s] was promoted as the channel ADMIN\n", adminUser)
				admin = clientStr.client
			}

			clients[clientStr.client] = true
			users[adminUser] = clientStr.client
			connections[adminUser] = clientStr.connection

		case cli := <-leaving:
			if admin == cli {

				for newAdmin := range clients {
					admin = newAdmin
					newAdmin <- serverName + "You're the new admin!"
					continue
				}
			}

			delete(clients, cli)
			close(cli)

		case user := <-kick:
			connection := connections[user]
			client := users[user]
			delete(clients, client)
			delete(connections, user)
			delete(users, user)
			connection.Close()
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {

	var buf = make([]byte, 1024)
	conn.Read(buf)

	myUser := string(bytes.Trim(buf, "\x00"))
	adminUser = string(bytes.Trim(buf, "\x00"))

	ch := make(chan string) // outgoing client messages

	fmt.Printf("%s New connected user [%s] \n", serverName, myUser)

	go clientWriter(conn, ch)

	ch <- serverName + "Welcome to the Simple IRC Server"
	ch <- serverName + "Your user [" + myUser + "] is successfully logged"

	messages <- serverName + "New connected user [" + myUser + "]"
	entering <- clientStruct{ch, conn}

	input := bufio.NewScanner(conn)

	//will read inputs from clients
	for input.Scan() {
		if len(input.Text()) > 0 && string(input.Text()[0]) == "/" {
			slice := strings.Split(input.Text(), " ")
			command := slice[0]
			switch command {

			//get list of users
			case "/users":
				str := ""
				for listUsers := range users {
					str += listUsers + ", "
				}
				ch <- serverName + str
			//Send message
			case "/msg":
				if len(slice) < 3 {
					ch <- "Command usage: /msg <user> <message>"
					continue
				}
				targetUser := slice[1]

				if _, ok := users[targetUser]; ok {
					directMessage := input.Text()[strings.Index(input.Text(), targetUser)+len(targetUser)+1:]
					direct <- userStruct{targetUser, myUser + " > " + directMessage}
				} else {
					ch <- "User: " + targetUser + " not found in the user list"
				}

			case "/time":

				timezone := "America/Mexico_City"
				loc, _ := time.LoadLocation(timezone)
				theTime := time.Now().In(loc).Format("15:04\n")
				ch <- "Local Time: " + timezone + " " + strings.Trim(theTime, "\n")

			case "/user":
				if len(slice) < 2 {
					ch <- "Command usage: /user <user>"
					continue
				}
				user := slice[1]
				if _, ok := users[user]; ok {
					ip := connections[user].RemoteAddr().String()
					ch <- "Username: " + user + " IP: " + ip

				} else {
					ch <- "User: " + user + " doesn't exist"
				}

			case "/kick":
				if len(slice) < 2 {
					ch <- "Command usage: /kick <user>"
					continue
				}
				if ch == admin {
					user := slice[1]

					if _, ok := users[user]; ok {
						messages <- serverName + "[" + user + "] was kicked from channel"
						fmt.Println(serverName + "[" + user + "] was kicked from channel")
						kick <- user
					} else {
						ch <- "User: " + user + " doesn't exist"
					}
				} else {
					ch <- "You don't have permission to use the application"
				}
			default:
				ch <- "Sub-commands: /users /user /msg /kick /time"

			}

		} else {
			messages <- myUser + " > " + input.Text()
		}
	}

	leaving <- ch
	messages <- serverName + "[" + myUser + "] left channel"
	fmt.Println(serverName + "[" + myUser + "] left channel")
	delete(users, myUser)
	delete(connections, myUser)
	conn.Close()
}

func clientWriter(conn net.Conn, channel <-chan string) {
	for msg := range channel {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	if len(os.Args) != 5 {
		fmt.Println("Expected:")
		fmt.Println("./server -host localhost -port <port>")
	}

	host := flag.String("host", "localhost", "localhost")
	port := flag.String("port", "8000", "port")
	flag.Parse()

	var address = "" + *host + ":" + *port

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(serverName + "Simple IRC Server started at " + address)
	if strings.Compare(address, listener.Addr().String()) != 0 {
		fmt.Println(serverName + "or at " + listener.Addr().String())
	}
	fmt.Println(serverName + "Ready for receiving new clients")

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main
