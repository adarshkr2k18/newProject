package main

import (
	"fmt"
	// "os"
	// "io/ioutil"
	// "log"
	// "strings"
	"path/filepath"
	"encoding/json"
	// "runtime"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"

	// "golang.org/x/image/colornames"

	"fyne.io/fyne/v2/container"
	// "github.com/vova616/screenshot"
	// "fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
)

type Config struct {
	myContainer   *fyne.Container
}
var myFile Config
// myLeftStack := stack.New()
// myRightStack := stack.New()

var myApp fyne.App
var myWindow fyne.Window
func openingFileManager(myApp1 fyne.App, myWindow1 fyne.Window,currentClient Client) {
	myApp=myApp1
	myWindow=myWindow1
	data1 := []byte("Open_File_manager")
	_, _ = currentClient.clientId.Write(data1) 
    // createFileDialog(myWindow)
}
func openFileManager(currentClient Client,dirPath string){
	fmt.Println("dfdffd")
	if(dirPath==""){
		fmt.Println("fff")
	getWindowsDrives(currentClient)
	}else{
		fmt.Println(dirPath)
		setContainor(dirPath,currentClient)
	}
}

func createFileDialog(){
	fmt.Println("dfd")
	d := dialog.NewCustom(
		"File manageer", 
		"Cancel", 
		container.NewScroll(myFile.myContainer),
		myWindow)
	d.Resize(fyne.Size{Width: 600,Height:600})

	d.Show()
}

func setHeader(dirPath string,currentClient Client) *fyne.Container{
		b1:=widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
			fmt.Println("Back")
			if(len(dirPath)==3){
				getWindowsDrives(currentClient)
			}else{
				setContainor(filepath.Dir(dirPath),currentClient)
			}
		})
		if(dirPath==""){
			b1.Disable()
		}
		b2:=widget.NewButtonWithIcon("", theme.DownloadIcon(), func() {
			fmt.Println("2")
			menuItem1 := fyne.NewMenuItem("A", nil)
		menuItem2 := fyne.NewMenuItem("B", nil)
		menuItem3 := fyne.NewMenuItem("C", nil)
		menu := fyne.NewMenu("File", menuItem1, menuItem2, menuItem3)

		popUpMenu := widget.NewPopUpMenu(menu, myWindow.Canvas())

		// popUpMenu.ShowAtPosition(*Expect mouse position here*)
		popUpMenu.Show()
		})
		b3:=widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
			fmt.Println("3")
		})
		b4:=widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {
			fmt.Println("4")
			getWindowsDrives(currentClient)
		})

		// Label on the right side
		// inputString := binding.NewString()
		// largeText :=widget.NewEntryWithData(inputString)
		// largeText.SetText("dirPath")
		largeText:=widget.NewLabel(dirPath)
		largeText.Resize(fyne.Size{Width: 200})
	header := container.NewHBox(
		b1,b2,b3,b4,largeText,
	)
	return header;
}
func setContainor( dirPath string,currentClient Client){
	container := container.NewVBox(setHeader(dirPath,currentClient)) 
	// data1 := []byte("Open_File_manager"+dirPath)
	// _, _ = currentClient.clientId.Write(data1)

	buffer := make([]byte, 1024)
	n, err := currentClient.clientId.Read(buffer)
	if err != nil {
		fmt.Println("Error reading data:", err)
	}
	var entries []string
	err = json.Unmarshal(buffer[:n], &entries)   // Unmarshal the JSON data
	if err != nil {
		fmt.Println("Error decoding JSON data:", err)
	}
	fmt.Println("entries",entries)

	// for _, entry := range entries {    // Iterate through the entries
	// 	// fmt.Println(entry.Name())
	// 	if strings.EqualFold(entry.Name(), "System Volume Information") {
	// 		continue
	// 	}
	// 	var fullPath=filepath.Join(dirPath, entry.Name())
	// 	if entry.IsDir() {
	// 		fmt.Println(entry.Name()+"Type: Directory")
	// 		w:=widget.NewButtonWithIcon(entry.Name(), theme.FolderIcon(), func() {
	// 			setContainor(fullPath,currentClient)	
	// 				fmt.Println(string(fullPath))
	// 		})	
	// 		w.Alignment = widget.ButtonAlignLeading
	// 		container.Add(w)
	// 	} else {
	// 		fmt.Println(entry.Name()+"Type: File")
	// 		w:=widget.NewButtonWithIcon(entry.Name(), theme.FileIcon(), func() {
	// 				fmt.Println(string(fullPath))
	// 		})	
	// 		w.Alignment = widget.ButtonAlignLeading
	// 		container.Add(w)
	// 	}
		
		
	// }
	
		container.Add(layout.NewSpacer())
		myFile.myContainer= container
		createFileDialog()

}

func getWindowsDrives(currentClient Client) {

	buffer := make([]byte, 1024)
	n, err := currentClient.clientId.Read(buffer)
	if err != nil {
		fmt.Println("Error reading data:", err)
	}
	var drives []string
	err = json.Unmarshal(buffer[:n], &drives)   // Unmarshal the JSON data
	if err != nil {
		fmt.Println("Error decoding JSON data:", err)
	}
	fmt.Println("drives",drives)
	container := container.NewVBox(setHeader("",currentClient)) 
		fmt.Println("Drives:")
		for _, drive := range drives {
			//  fmt.Printf("%T\n", drive) 
			var dirPath=drive
			w:=widget.NewButtonWithIcon("Local Disk (" + string(drive[0]) + ":)", theme.DocumentPrintIcon(), func() {
				// setContainor(d,currentClient)	
				data1 := []byte("Open_File_manager"+dirPath)
				_, _ = currentClient.clientId.Write(data1)
					fmt.Println(string(dirPath))
			})				
			w.Alignment = widget.ButtonAlignLeading
			container.Add(w)
		}
		container.Add(layout.NewSpacer())
		myFile.myContainer= container
	createFileDialog()
}