package main

import (
	controller "DevKit-Neuro-server/proto"
	"bufio"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"os"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Please provide host:port to connect to")
		os.Exit(1)
	}

	// Resolve the string address to a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", os.Args[1])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Dial to the address with UDP
	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data := &controller.ChannelsDataSet{
		Channel1: 1.1,
		Channel2: 1.2,
		Channel3: 1.3,
		Channel4: 1.4,
		Channel5: 1.5,
		Channel6: 1.6,
		Channel7: 1.7,
		Channel8: 1.8,
		Id:       1,
	}

	//metrics := &controller.EmgMetric{
	//	Metrics: data,
	//}

	msg, err := proto.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}

	// Send a message to the server
	fmt.Println(msg)
	_, err = conn.Write(msg)
	fmt.Println("send...")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Ниже как пример для сервера
	msg_in := &controller.RawDataPack{}

	err = proto.Unmarshal(msg, msg_in)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	fmt.Println(msg_in)

	// Read from the connection untill a new line is send
	responce, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the data read from the connection to the terminal
	fmt.Print("> ", string(responce))
}
