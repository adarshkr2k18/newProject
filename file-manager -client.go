package main

import (
	"fmt"
	"os"
	"io/ioutil"
	// "log"
	// "strings"
	// "path/filepath"
	"net"
	"encoding/json"
	// "runtime"
	"fyne.io/fyne/v2"
	// "fyne.io/fyne/v2/widget"
	// "fyne.io/fyne/v2/layout"

	// "golang.org/x/image/colornames"

	// "fyne.io/fyne/v2/container"
	// "github.com/vova616/screenshot"
	// "fyne.io/fyne/v2/data/binding"
	// "fyne.io/fyne/v2/dialog"
	// "fyne.io/fyne/v2/theme"


	
)

type Config struct {
	myContainer   *fyne.Container
}
var myFile Config
// myLeftStack := stack.New()
// myRightStack := stack.New()

func openFileManager(serverId net.Conn,dirPath string) {
	if(dirPath==""){
		fmt.Println("fff")
		sendWindowsDrives(serverId)
	}else{
		fmt.Println(dirPath)
		sendFileAndFolder(serverId,dirPath)
	}
    // createFileDialog(myWindow)
	fmt.Println("dfdd")
}

// func createFileDialog(myWindow fyne.Window){
// 	d := dialog.NewCustom(
// 		"File manageer", 
// 		"Cancel", 
// 		container.NewScroll(myFile.myContainer),
// 		myWindow)
// 	d.Resize(fyne.Size{Width: 600,Height:600})

// 	d.Show()
// }

// func setHeader(myWindow fyne.Window,dirPath string) *fyne.Container{
// 		b1:=widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
// 			fmt.Println("Back")
// 			if(len(dirPath)==3){
// 				sendWindowsDrives(myWindow)
// 			}else{
// 				setContainor(myWindow,filepath.Dir(dirPath))
// 			}
// 		})
// 		if(dirPath==""){
// 			b1.Disable()
// 		}
// 		b2:=widget.NewButtonWithIcon("", theme.DownloadIcon(), func() {
// 			fmt.Println("2")
// 			menuItem1 := fyne.NewMenuItem("A", nil)
// 		menuItem2 := fyne.NewMenuItem("B", nil)
// 		menuItem3 := fyne.NewMenuItem("C", nil)
// 		menu := fyne.NewMenu("File", menuItem1, menuItem2, menuItem3)

// 		popUpMenu := widget.NewPopUpMenu(menu, myWindow.Canvas())

// 		// popUpMenu.ShowAtPosition(*Expect mouse position here*)
// 		popUpMenu.Show()
// 		})
// 		b3:=widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
// 			fmt.Println("3")
// 		})
// 		b4:=widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {
// 			fmt.Println("4")
// 			sendWindowsDrives(myWindow)
// 		})

// 		// Label on the right side
// 		// inputString := binding.NewString()
// 		// largeText :=widget.NewEntryWithData(inputString)
// 		// largeText.SetText("dirPath")
// 		largeText:=widget.NewLabel(dirPath)
// 		largeText.Resize(fyne.Size{Width: 200})
// 	header := container.NewHBox(
// 		b1,b2,b3,b4,largeText,
// 	)
// 	return header;
// }
// func setContainor(myWindow fyne.Window, dirPath string){
// 	container := container.NewVBox(setHeader(myWindow,dirPath)) 
// 	entries, err := ioutil.ReadDir(dirPath)  // Get a list of files and folders in the specified directory
// 	if err != nil {
// 		log.Fatal(err)
// 	}	
// 	for _, entry := range entries {    // Iterate through the entries
// 		// fmt.Println(entry.Name())
// 		if strings.EqualFold(entry.Name(), "System Volume Information") {
// 			continue
// 		}
// 		var fullPath=filepath.Join(dirPath, entry.Name())
// 		if entry.IsDir() {
// 			fmt.Println(entry.Name()+"Type: Directory")
// 			w:=widget.NewButtonWithIcon(entry.Name(), theme.FolderIcon(), func() {
// 				setContainor(myWindow,fullPath)	
// 					fmt.Println(string(fullPath))
// 			})	
// 			w.Alignment = widget.ButtonAlignLeading
// 			container.Add(w)
// 		} else {
// 			fmt.Println(entry.Name()+"Type: File")
// 			w:=widget.NewButtonWithIcon(entry.Name(), theme.FileIcon(), func() {
// 					fmt.Println(string(fullPath))
// 			})	
// 			w.Alignment = widget.ButtonAlignLeading
// 			container.Add(w)
// 		}
		
		
// 	}
	
// 		container.Add(layout.NewSpacer())
// 		myFile.myContainer= container
// 		createFileDialog(myWindow)

// }

func sendWindowsDrives(serverId net.Conn) {
	var drives []string
	for _, driveLetter := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		drivePath := string(driveLetter) + ":\\"
		_, err := os.Stat(drivePath)
		if err == nil {
			// fmt.Println(drives, drivePath)
			drives = append(drives, drivePath)
		}
	}
	fmt.Println(drives)
	
	data1 := []byte("Open_File_manager")
	_, _ = serverId.Write(data1)

	jsonData, err := json.Marshal(drives)
	if err != nil {
		fmt.Println(err)
	}

	// Send the JSON data to the server
	_, err = serverId.Write(jsonData)
	if err != nil {
		fmt.Println(err)
	}

}

func sendFileAndFolder(serverId net.Conn,dirPath string){
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(entries,entries[0].Name())
	fmt.Printf("%T",entries[0])
	jsonData, err := json.Marshal(entries)
	if err != nil {
		fmt.Println(err)
	}

	data1 := []byte("Open_File_manager"+dirPath)
	_, _ = serverId.Write(data1)

	// Send the JSON data to the server
	_, err = serverId.Write(jsonData)
	if err != nil {
		fmt.Println(err)
	}

}