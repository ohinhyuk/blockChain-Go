package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func startServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Server is listening on port " + port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		go handleClient(conn)
	}
}



// 사용자 입력을 처리하고 서버로 메시지를 전송하는 함수
func sendUserInput(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)

		if msg == "exit" {
			fmt.Println("Exiting user input loop.")
			return
		}

		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error writing to server:", err.Error())
			return
		}
	}
}





// 기존 connectToPeer 함수에 사용자 입력 처리 부분 추가
func connectToPeer(peerIP string, peerPort string) {
	conn, err := net.Dial("tcp", peerIP+":"+peerPort)
	if err != nil {
		fmt.Println("Connection Failed:", err.Error())
		return
	}
	defer conn.Close()

	fmt.Println("Connected to peer. Type 'exit' to stop.")
	go sendUserInput(conn)

	// 여기서는 클라이언트로부터 받은 메시지를 계속 읽고 출력합니다.
	handleClient(conn)
}



func handleClient(conn net.Conn) {
    defer conn.Close()

    // 버퍼 크기 설정
    buf := make([]byte, 1024)

	fmt.Println("debug")

    for {
		fmt.Println("debug2")
        // 클라이언트로부터 메시지 읽기
        len, err := conn.Read(buf)
		
		fmt.Println("len: ", len)

        if err != nil {
            fmt.Println("Error reading:", err.Error())
            return
        }
        recvMsg := string(buf[:len])
        fmt.Println("Received:", recvMsg)

        // 메시지 응답
        var response string
        fmt.Print("Enter response: ")
        fmt.Scanln(&response)
        _, err = conn.Write([]byte(response))
        if err != nil {
            fmt.Println("Error writing:", err.Error())
            return
        }
    }
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage:", os.Args[0], "<server_port> <peer_ip> <peer_port>")
		os.Exit(1)
	}

	serverPort := os.Args[1]
	peerIP := os.Args[2]
	peerPort := os.Args[3]

	go startServer(serverPort)
	go connectToPeer(peerIP, peerPort)

	// Wait indefinitely
	select {}
}
