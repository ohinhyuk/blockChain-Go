package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// Node 구조체는 네트워크 내의 다른 노드를 나타냅니다.
type Node struct {
	IP   string // 노드의 IP 주소
	Port string // 노드의 포트 번호
}

// NewNode는 새로운 노드를 생성합니다.
func NewNode(ip, port string) *Node {
	return &Node{IP: ip, Port: port}
}

// Blockchain은 블록들의 연결 리스트와 네트워크 노드 목록을 유지합니다.
// type Blockchain struct {
// 	Blocks []*Block
// 	Nodes  []*Node
// }

// AddNode는 새로운 노드를 블록체인 네트워크에 추가합니다.
// AddNode는 새로운 노드를 블록체인 네트워크에 추가합니다.
// 이미 존재하는 노드는 추가하지 않습니다.
func (bc *Blockchain) AddNode(node *Node) {
	if !bc.isNodeExists(node) {
		bc.Nodes = append(bc.Nodes, node)
		fmt.Printf("New node added: %s:%s\n", node.IP, node.Port)
	} else {
		fmt.Printf("Node already exists: %s:%s\n", node.IP, node.Port)
	}
}

// isNodeExists 함수는 주어진 노드가 이미 노드 리스트에 존재하는지 확인합니다.
func (bc *Blockchain) isNodeExists(newNode *Node) bool {
	for _, node := range bc.Nodes {
		if node.IP == newNode.IP && node.Port == newNode.Port {
			return true
		}
	}
	return false
}
// BroadcastBlock은 네트워크에 새로운 블록을 전파합니다.
func (bc *Blockchain) BroadcastBlock(block *Block) {
	for _, node := range bc.Nodes {
		sendBlock(node, block)
		fmt.Printf("Broadcasting block to %s:%s\n", node.IP, node.Port)
	}
}

// sendBlock은 특정 노드에게 블록을 전송합니다.
func sendBlock(node *Node, block *Block) {
	address := fmt.Sprintf("%s:%s", node.IP, node.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Error connecting to node %s: %s\n", address, err)
		return
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(block)
	if err != nil {
		fmt.Printf("Error encoding block: %s\n", err)
		return
	}
}

func startMining(bc *Blockchain, stopCh chan struct{}) {
	for {
		select {
		case <-stopCh:
			
			time.Sleep(5 * time.Second) // 채굴 중단 후 10초 대기
			continue
		default:
			fmt.Println("채굴 시작!!")
			newBlock := NewBlock("New Block Data", bc.Blocks[len(bc.Blocks)-1].Hash)
			pow := NewProofOfWork(newBlock)
			nonce, _ := pow.Run()

			if nonce != -1 {
				fmt.Println("채굴 성공!! 리더가 되었습니다!!!")
				bc.Blocks = append(bc.Blocks, newBlock)
				bc.BroadcastBlock(newBlock)
				time.Sleep(5 * time.Second)

			} else {
				fmt.Println("채굴 실패. 다시 시도.")
			}
			// time.Sleep(10 * time.Second) // 다음 채굴 시도 전 10초 대기
		}
	}
}



func handleClient(conn net.Conn, bc *Blockchain, stopCh chan struct{}) {
	defer conn.Close()

	decoder := gob.NewDecoder(conn)
	var block Block
	err := decoder.Decode(&block)
	if err != nil {
		fmt.Println("Error decoding block:", err)
		return
	}

	if isBlockPresent(bc, &block) {
		// fmt.Println("중복된 블록이 들어왔습니다. 블록을 무시합니다.")
		return
	}

	if isValidBlock(&block, bc) {
		fmt.Println("채굴 실패.. 다른 사람이 리더가 되었습니다..")
		bc.Blocks = append(bc.Blocks, &block)
		 bc.BroadcastBlock(&block)
		
		time.Sleep(5 * time.Second)
		// stopCh <- struct{}{} // 채굴 중단 신호 전송
		
		 // 채굴 중단 후 대기

		
	} else {
		fmt.Println("수신된 블록이 유효하지 않음.")
	}
}



// isBlockPresent 함수는 주어진 블록이 블록체인에 이미 존재하는지 확인합니다.
func isBlockPresent(bc *Blockchain, block *Block) bool {
	for _, blk := range bc.Blocks {
		if bytes.Equal(blk.Hash, block.Hash) {
			return true
		}
	}
	return false
}

// isValidBlock 함수는 주어진 블록의 유효성을 검증합니다.
func isValidBlock(block *Block, bc *Blockchain) bool {
	// 여기에 블록의 유효성 검증 로직을 추가합니다.
	return true
}


func startServer(port string, bc *Blockchain, stopCh chan struct{}) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}

		go handleClient(conn, bc , stopCh)
	}
}

// 사용자 입력을 처리하고 블록을 생성하여 네트워크에 전파하는 함수
func createAndSendBlock(conn net.Conn, bc *Blockchain) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message to create a block: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)

		if msg == "exit" {
			fmt.Println("Exiting block creation loop.")
			return
		}

		// 여기에서 새로운 블록을 생성하고 네트워크에 전파합니다.
		// newBlock := NewBlock(msg, []byte{}) // 간단한 예시, 실제로는 이전 블록 해시 등 필요
		newBlock := bc.AddBlock(msg)
		bc.BroadcastBlock(newBlock)
	}
}

// connectToPeer 함수는 특정 노드에 연결합니다.
func connectToPeer(peerIP string, peerPort string, bc *Blockchain, stopCh chan struct{}) {
	conn, err := net.Dial("tcp", peerIP+":"+peerPort)

	newNode := NewNode(peerIP, peerPort)
	bc.AddNode(newNode)


	if err != nil {
		fmt.Println("Connection Failed:", err.Error())
		return
	}
	defer conn.Close()

	fmt.Println("Connected to peer. Type 'exit' to stop.")
	// go createAndSendBlock(conn, bc)

	// 클라이언트로부터 받은 메시지를 계속 읽고 출력합니다.
	// for {
	// 	buf := make([]byte, 1024)
	// 	len, err := conn.Read(buf)
	// 	if err != nil {
	// 		fmt.Println("Error reading:", err.Error())
	// 		break
	// 	}
	// 	fmt.Println("Received:", string(buf[:len]))
	// }
}

func createAndBroadcastBlock(bc *Blockchain) {
	for {
		
		time.Sleep(10 * time.Second) // 채굴 간격을 10초로 조정
		
		fmt.Println("채굴 시작~~")
		// 채굴 과정 (블록 생성 및 PoW 실행)
		newBlock := bc.AddBlock("New Block Data")
		pow := NewProofOfWork(newBlock)
		nonce, _ := pow.Run()

		if nonce != -1 {
			fmt.Println("채굴 성공!! 리더가 되었습니다!!!")
			bc.Blocks = append(bc.Blocks, newBlock)
			bc.BroadcastBlock(newBlock)
		} else {
			fmt.Println("리더에 실패하였습니다..! 다시 채굴 시작!")
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

	bc := NewBlockchain()
	stopCh := make(chan struct{})

	go startServer(serverPort, bc, stopCh)
	go connectToPeer(peerIP, peerPort, bc, stopCh)

	go startMining(bc, stopCh)

	select {}
}
