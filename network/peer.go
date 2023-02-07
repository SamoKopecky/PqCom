package network

import (
	"bufio"
	"fmt"
	"os"
)

func Start(laddr, daddr string, lport, dport int) {
	go listen(laddr, lport)
	for {
		fmt.Println("Press key to start sending")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		send(daddr, dport)
	}

}
