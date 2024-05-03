package winapi

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	moduser32                                  = syscall.NewLazyDLL("user32.dll")
	procGetDesktopWindow                       = moduser32.NewProc("GetDesktopWindow")
	procGetWindowRect                          = moduser32.NewProc("GetWindowRect")
	procGetClientRect                          = moduser32.NewProc("GetClientRect")
	procScreenToClient                         = moduser32.NewProc("ScreenToClient")
	procGetWindowDC                            = moduser32.NewProc("GetWindowDC")
	procReleaseDC                              = moduser32.NewProc("ReleaseDC")
	procPhysicalToLogicalPointForPerMonitorDPI = moduser32.NewProc("PhysicalToLogicalPointForPerMonitorDPI")
)

// Retrieves a handle to the desktop window.
// The desktop window covers the entire screen.
// The desktop window is the area on top of which other windows are painted.
func GetDesktopWindow() windows.HWND {
	result, _, _ := procGetDesktopWindow.Call()
	return windows.HWND(result)
}

// Retrieves the dimensions of the bounding rectangle of the specified window. The dimensions are given in screen coordinates that are relative to the upper-left corner of the screen.
func GetWindowRect(handleWindow windows.HWND) (windows.Rect, error) {
	result := windows.Rect{}
	ret, _, _ := procGetWindowRect.Call(uintptr(handleWindow), uintptr(unsafe.Pointer(&result)))
	if ret == 0 {
		return result, fmt.Errorf("[GetWindowRect] failed: %s", windows.GetLastError())
	}
	return result, nil
}

// Retrieves the coordinates of a window's client area. The client coordinates specify the upper-left and lower-right corners of the client area.
// Because client coordinates are relative to the upper-left corner of a window's client area, the coordinates of the upper-left corner are (0,0).
func GetClientRect(handleWindow windows.HWND) (windows.Rect, error) {
	result := windows.Rect{}
	ret, _, _ := procGetClientRect.Call(uintptr(handleWindow), uintptr(unsafe.Pointer(&result)))
	if ret == 0 {
		return result, fmt.Errorf("[GetClientRect] failed: %s", windows.GetLastError())
	}
	return result, nil
}

// The ScreenToClient function converts the screen coordinates of a specified point on the screen to client-area coordinates.
func ScreenToClient(handleWindow windows.HWND, point POINT) (POINT, error) {
	ret, _, err := procScreenToClient.Call(uintptr(handleWindow), uintptr(unsafe.Pointer(&point)))
	if ret == 0 {
		return point, fmt.Errorf("[ScreenToClient] failed: %s", err)
	}
	return point, nil
}

// The GetWindowDC function retrieves the device context (DC) for the entire window, including title bar, menus, and scroll bars.
// A window device context permits painting anywhere in a window, because the origin of the device context is the upper-left corner of the window instead of the client area.
// GetWindowDC assigns default attributes to the window device context each time it retrieves the device context. Previous attributes are lost.
func GetWindowDC(handleWindow windows.HWND) (HDC, error) {
	result, _, err := procGetWindowDC.Call(uintptr(handleWindow))
	if result == 0 {
		return HDC(result), fmt.Errorf("[GetWindowDC] failed: %s", err)
	}
	return HDC(result), nil
}

// The ReleaseDC function releases a device context (DC), freeing it for use by other applications.
// The effect of the ReleaseDC function depends on the type of DC.
// It frees only common and window DCs. It has no effect on class or private DCs.
func ReleaseDC(handleWindow windows.HWND, handleDeviceContext HDC) error {
	ret, _, err := procReleaseDC.Call(uintptr(handleWindow), uintptr(handleDeviceContext))
	if ret == 0 {
		return fmt.Errorf("[ReleaseDC] failed: %s", err)
	}
	return nil
}

func PhysicalToLogicalPointForPerMonitorDPI(handleWindow windows.HWND, point POINT) (POINT, error) {
	ret, _, err := procPhysicalToLogicalPointForPerMonitorDPI.Call(uintptr(handleWindow), uintptr(unsafe.Pointer(&point)))
	if ret == 0 {
		return point, fmt.Errorf("[PhysicalToLogicalPointForPerMonitorDPI] failed: %s", err)
	}
	return point, nil
}
