package tcp

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func ClientSetup(host string) net.Conn {
	CONNECT := host
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return c
}

func WriteToServer(data []byte, c net.Conn) []byte {
	if c != nil {
		data = append(data, []byte("\n")[0])
		c.Write(data)
		return data
	}
	fmt.Println("tcpC: failed to send")
	return nil
}

func ReadFromServer(c net.Conn) []byte {
	if c != nil {
		text, _ := bufio.NewReader(c).ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		return []byte(text)
	}
	fmt.Println("tspC: failed to read")
	return nil
}
