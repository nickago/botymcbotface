package main

import (
	"fmt"
  "io"
	"net"
  "os"
  "strings"
)

const USERNAME string = "attila_the_bot"

func main() {
  var connData [1024]byte

	addr, err := net.ResolveTCPAddr("tcp", "chat.freenode.net:6667")
  checkError(err)

  conn, err := net.DialTCP("tcp", nil, addr)
  checkError(err)

  conn.Write([]byte(fmt.Sprintf("NICK %s\r\n", USERNAME)))
	conn.Write([]byte(fmt.Sprintf("USER %s 0 * :(Here lies %s)\r\n", USERNAME, USERNAME)))

  for {
    n, err := conn.Read(connData[0:])
    if err == io.EOF {
      println("Connection Terminated")
      os.Exit(1)
    }
    dataString := string(connData[:n])

    for _, line := range strings.Split(dataString, "\n") {
      go handleConnection(conn, line)
    }

    println(dataString)
    println("----------------------------------------------")
  }
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
  }

func handleConnection(conn *net.TCPConn, dataString string) {
  data := strings.Split(dataString, " ")

  if data[0] == "PING" {
    conn.Write([]byte(fmt.Sprintf("PONG %s\r\n", data[1])))
  }

  if len(data) > 1 && (data[1] == "376" || data[1] == "422") {
    conn.Write([]byte(fmt.Sprintf(":%s JOIN ##whitehat\r\n", USERNAME)))
    conn.Write([]byte(fmt.Sprintf(":%s PRIVMSG ##whitehat :I'm Alive\r\n", USERNAME)))
  }

  if len(data) > 1 && data[1] == "QUIT" {
    os.Exit(1)
  }

  if len(data) > 1 && data[1] == "PRIVMSG" && data[2] == USERNAME {
    println("Stuff")
    sender := strings.Split(data[0], "!")[0]
    conn.Write([]byte(fmt.Sprintf(":%s PRIVMSG %s :This isn't the turing test.", USERNAME, sender)))
  }
}
