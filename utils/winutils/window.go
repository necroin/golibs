//go:build windows

package winutils

import (
	"sync"
	"syscall"
	"unsafe"

	"github.com/necroin/golibs/utils/winapi"
	"golang.org/x/sys/windows"
)

type GetWindowHandlesByProcessIdCallbackWrapper struct {
	pid      winapi.ProcessId
	result   []windows.HWND
	callback uintptr
	mutex    sync.Mutex
}

var (
	GetWindowHandlesByProcessIdCallback = GetWindowHandlesByProcessIdCallbackWrapper{
		result: []windows.HWND{},
		mutex:  sync.Mutex{},
	}
)

func init() {
	GetWindowHandlesByProcessIdCallback.callback = syscall.NewCallback(func(hwnd windows.HWND, lParam *windows.HWND) uintptr {
		processId := uint32(0)
		windows.GetWindowThreadProcessId(hwnd, &processId)
		if GetWindowHandlesByProcessIdCallback.pid == winapi.ProcessId(processId) {
			*lParam = hwnd
			GetWindowHandlesByProcessIdCallback.result = append(GetWindowHandlesByProcessIdCallback.result, hwnd)
		}
		return 1
	})
}

func GetWindowHandlesByProcessId(pid winapi.ProcessId) []windows.HWND {
	var hwnd windows.HWND
	GetWindowHandlesByProcessIdCallback.mutex.Lock()
	defer GetWindowHandlesByProcessIdCallback.mutex.Unlock()
	GetWindowHandlesByProcessIdCallback.pid = pid
	GetWindowHandlesByProcessIdCallback.result = []windows.HWND{}
	windows.EnumWindows(GetWindowHandlesByProcessIdCallback.callback, unsafe.Pointer(&hwnd))
	return GetWindowHandlesByProcessIdCallback.result
}
