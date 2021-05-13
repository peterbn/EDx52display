package main

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	S_OK = 0x00000000
)

const (
	pluginName = "TestPlugin"
)

var directOutput *syscall.LazyDLL

var context uintptr = 0xCAFEBABE

var device uintptr = 0

var pages = 0

func main() {

	directOutput = syscall.NewLazyDLL("./DepInclude/DirectOutput.dll")

	initialize(pluginName)
	defer deinitialize()

	registerDeviceCallback()
	enumerate()

	reader := bufio.NewReader(os.Stdin)
	reader.ReadLine()

}

func initPages() {
	if pages == 0 {
		fmt.Println("initPages", device)
		registerPageCallback(device)
		registerSoftButtonCallback(device)

		addPage(0, true)
		showPage(0)
		pages = pages + 1
	}
}

func showPage(pagenum uint32) {
	if pagenum != 0 {
		return
	}
	setString(0, 1, "Hello, World!")
}

func initialize(pluginName string) {
	proc := directOutput.NewProc("DirectOutput_Initialize")
	pluginNamePtr, err := syscall.UTF16PtrFromString(pluginName)
	if err != nil {
		panic(err)
	}
	hresult, _, err := proc.Call(uintptr(unsafe.Pointer(pluginNamePtr)))

	if hresult != S_OK {
		panic(err)
	}
}

func deinitialize() {
	proc := directOutput.NewProc("DirectOutput_Deinitialize")
	hresult, _, err := proc.Call()
	if hresult != S_OK {
		panic(err)
	}
}

func registerDeviceCallback() {
	proc := directOutput.NewProc("DirectOutput_RegisterDeviceCallback")
	callback := syscall.NewCallback(onDeviceChanged)

	hresult, _, err := proc.Call(callback, context)
	if hresult != S_OK {
		panic(err)
	}
}

func enumerate() {
	proc := directOutput.NewProc("DirectOutput_Enumerate")
	callback := syscall.NewCallback(onEnumerate)

	hresult, _, err := proc.Call(callback, context)
	if hresult != S_OK {
		panic(err)
	}
}

func registerPageCallback(device uintptr) {
	proc := directOutput.NewProc("DirectOutput_RegisterPageCallback")

	callback := syscall.NewCallback(onPageChange)
	hresult, _, err := proc.Call(device, callback, context)
	if hresult != S_OK {
		panic(err)
	}
}

func registerSoftButtonCallback(device uintptr) {
	proc := directOutput.NewProc("DirectOutput_RegisterSoftButtonCallback")

	callback := syscall.NewCallback(onSoftButton)
	hresult, _, err := proc.Call(device, callback, context)
	if hresult != S_OK {
		panic(err)
	}
}

func addPage(pageNumber int, active bool) {
	proc := directOutput.NewProc("DirectOutput_AddPage")
	var flag uintptr = 0
	if active {
		flag = 1
	}
	hresult, _, err := proc.Call(device, uintptr(pageNumber), flag)
	if hresult != S_OK {
		panic(err)
	}
}

func setString(page, index uint32, line string) {
	fmt.Println("setString", page, index, line)
	proc := directOutput.NewProc("DirectOutput_SetString")
	linePtr, _ := syscall.UTF16PtrFromString(line)

	lineLen := uintptr(len(line))

	hresult, _, err := proc.Call(device, uintptr(page), uintptr(index), lineLen, uintptr(unsafe.Pointer(linePtr)))
	if hresult != S_OK {
		fmt.Println("hresult", hresult)
		panic(err)
	}
}
