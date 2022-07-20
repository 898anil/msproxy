package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

const version = "5.7.1"
const connectionId = 12

func PutLengthEncodedInt(n uint64) []byte {
	switch {
	case n <= 250:
		return []byte{byte(n)}

	case n <= 0xffff:
		return []byte{0xfc, byte(n), byte(n >> 8)}

	case n <= 0xffffff:
		return []byte{0xfd, byte(n), byte(n >> 8), byte(n >> 16)}

	case n <= 0xffffffffffffffff:
		return []byte{0xfe, byte(n), byte(n >> 8), byte(n >> 16), byte(n >> 24),
			byte(n >> 32), byte(n >> 40), byte(n >> 48), byte(n >> 56)}
	}
	return nil
}

func main() {
	// salt := []byte("dshjsdhfksjfdhsdsd")
	capabilities := 1
	SERVER_STATUS_AUTOCOMMIT := 0x0002
	var OK_HEADER byte = 0x00
	listener, err := net.Listen("tcp", ":7878")
	checkErr(err)
	running := true
	for running {
		conn, err := listener.Accept()
		fmt.Println("came")
		checkErr(err)
		data := make([]byte, 4, 128)
		data = append(data, 10)
		data = append(data, version...)
		data = append(data, byte(connectionId), byte(connectionId>>8), byte(connectionId>>16), byte(connectionId>>24))
		// data = append(data, salt[0:8]...)
		data = append(data, 0)
		data = append(data, byte(capabilities), byte(capabilities>>8))
		data = append(data, uint8(33))
		data = append(data, byte(SERVER_STATUS_AUTOCOMMIT), byte(SERVER_STATUS_AUTOCOMMIT>>8))
		data = append(data, byte(capabilities>>16), byte(capabilities>>24))
		data = append(data, 0x15)
		data = append(data, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
		// data = append(data, salt[8:]...)
		data = append(data, 0)
		length := len(data) - 4
		data[0] = byte(length)
		data[1] = byte(length >> 8)
		data[2] = byte(length >> 16)
		data[3] = 0
		conn.Write(data)

		reader := bufio.NewReaderSize(conn, 1024)
		header := make([]byte, 4)
		io.ReadFull(reader, header)
		data = make([]byte, header[0])
		io.ReadFull(reader, data)
		fmt.Println(string(data))

		okData := make([]byte, 4, 32)
		okData = append(okData, OK_HEADER)
		okData = append(okData, PutLengthEncodedInt(0)...)
		okData = append(okData, PutLengthEncodedInt(0)...)
		okData = append(okData, byte(SERVER_STATUS_AUTOCOMMIT), byte(SERVER_STATUS_AUTOCOMMIT>>8))
		okData = append(okData, 0, 0)
		length = len(okData) - 4
		okData[0] = byte(length)
		okData[1] = byte(length >> 8)
		okData[2] = byte(length >> 16)
		okData[3] = 3
		fmt.Printf("okData: %v\n", okData)
		conn.Write(okData)

		reader = bufio.NewReaderSize(conn, 1024)
		header = make([]byte, 4)
		io.ReadFull(reader, header)
		data = make([]byte, header[0])
		io.ReadFull(reader, data)
		fmt.Println(string(data))
		// conn.Close()
	}
}
