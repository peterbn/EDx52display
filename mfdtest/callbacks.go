package main

import "fmt"

func onEnumerate(hdevice uintptr, context uintptr) uintptr {
	fmt.Println("onEnumerate")
	fmt.Println("hdevice", hdevice)
	fmt.Println("context", context)

	device = hdevice

	initPages()
	return S_OK
}

func onDeviceChanged(hdevice uintptr, added bool, context uintptr) uintptr {
	fmt.Println("onDeviceChanged")
	fmt.Println("hdevice", hdevice)
	fmt.Println("added", added)
	fmt.Println("context", context)

	if added {
		device = hdevice
		initPages()
	} else {
		device = 0
	}

	return S_OK
}

func onPageChange(hdevice uintptr, page uint32, setActive bool, context uintptr) uintptr {
	fmt.Println("onPageChange")
	fmt.Println("hdevice", hdevice)
	fmt.Println("page", page)
	fmt.Println("setActive", setActive)
	fmt.Println("context", context)

	if setActive {
		showPage(page)
	}

	return S_OK
}

func onSoftButton(hdevice uintptr, buttons uint32, context uintptr) uintptr {
	fmt.Println("onSoftButton")
	fmt.Println("hdevice", hdevice)
	fmt.Println("buttons", buttons)
	fmt.Println("context", context)

	return S_OK
}
