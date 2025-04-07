//go:build windows

package winapi

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modkernel32      = syscall.NewLazyDLL("kernel32.dll")
	procGlobalAlloc  = modkernel32.NewProc("GlobalAlloc")
	procGlobalFree   = modkernel32.NewProc("GlobalFree")
	procGlobalLock   = modkernel32.NewProc("GlobalLock")
	procGlobalUnlock = modkernel32.NewProc("GlobalUnlock")
)

// Allocates the specified number of bytes from the heap.
func GlobalAlloc(flags uint32, bytes uintptr) (HGLOBAL, error) {
	result, _, _ := procGlobalAlloc.Call(uintptr(flags), bytes)
	if result == 0 {
		return HGLOBAL(result), fmt.Errorf("[GlobalAlloc] failed: %s", windows.GetLastError())
	}
	return HGLOBAL(result), nil
}

// Frees the specified global memory object and invalidates its handle.
func GlobalFree(handleMemory HGLOBAL) (HGLOBAL, error) {
	result, _, _ := procGlobalFree.Call(uintptr(handleMemory))
	if HGLOBAL(result) == handleMemory {
		return HGLOBAL(result), fmt.Errorf("[GlobalFree] failed: %s", windows.GetLastError())
	}
	return HGLOBAL(result), nil
}

// Locks a global memory object and returns a pointer to the first byte of the object's memory block.
func GlobalLock(handleMemory HGLOBAL) (unsafe.Pointer, error) {
	result, _, _ := procGlobalLock.Call(uintptr(handleMemory))
	if result == 0 {
		return unsafe.Pointer(result), fmt.Errorf("[GlobalLock] failed: %s", windows.GetLastError())
	}
	return unsafe.Pointer(result), nil
}

// Decrements the lock count associated with a memory object that was allocated with GMEM_MOVEABLE. This function has no effect on memory objects allocated with GMEM_FIXED.
func GlobalUnlock(handleMemory HGLOBAL) bool {
	ret, _, _ := procGlobalUnlock.Call(uintptr(handleMemory))
	return ret != 0
}
