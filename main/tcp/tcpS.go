package tcp

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func ServerSetup(port string) net.Conn {
	PORT := ":" + port
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println("Server is ready")
	return c
}

func WriteToClient(data []byte, c net.Conn) []byte {
	if c != nil {
		data = append(data, []byte("\n")[0])
		c.Write(data)
		return data
	}
	fmt.Println("tcpS: failed to send")
	return nil
}

func ReadFromClient(c net.Conn) []byte {
	if c != nil {
		text, _ := bufio.NewReader(c).ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		return []byte(text)
	}
	fmt.Println("tspS: failed to read")
	return nil
}
