package directoutput

import (
	"syscall"
	"unsafe"
)

type dllFunc func(a ...uintptr) error

var mod *syscall.DLL

var procInitialize, procDeinitialize dllFunc

// init will load the DLL required for directoutput to work
func setupBindings() {
	mod = syscall.MustLoadDLL("DirectOutput.dll")
	procInitialize = callableFromProc("DirectOutput_Initialize")
	procDeinitialize = callableFromProc("DirectOutput_Deinitialize")
}

func teardownBindings() {
	mod.Release()
}

// Initialize will initialize the DirectOutput module for use
func Initialize(appName string) error {
	setupBindings()

	pApppName, err := uintptrStringArg(appName)
	err = procInitialize(pApppName)
	return err
}

// Deinitialize cleans up the library
func Deinitialize() error {
	defer teardownBindings()
	return procDeinitialize()
}

func callableFromProc(procname string) dllFunc {
	p := mod.MustFindProc("DirectOutput_Initialize")
	return func(a ...uintptr) error {
		r1, _, err := p.Call(a...)
		if r1 < 0 {
			return err
		}
		return nil
	}
}

func uintptrStringArg(arg string) (uintptr, error) {
	ptr, err := syscall.UTF16PtrFromString(arg)
	return uintptr(unsafe.Pointer(ptr)), err
}
