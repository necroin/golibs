//go:build windows

package winapi

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modgdi32                   = syscall.NewLazyDLL("gdi32.dll")
	procCreateCompatibleDC     = modgdi32.NewProc("CreateCompatibleDC")
	procCreateCompatibleBitmap = modgdi32.NewProc("CreateCompatibleBitmap")
	procSelectObject           = modgdi32.NewProc("SelectObject")
	procBitBlt                 = modgdi32.NewProc("BitBlt")
	procGetDIBits              = modgdi32.NewProc("GetDIBits")
	procDeleteDC               = modgdi32.NewProc("DeleteDC")
	procDeleteObject           = modgdi32.NewProc("DeleteObject")
)

// The CreateCompatibleDC function creates a memory device context (DC) compatible with the specified device.
func CreateCompatibleDC(handleDeviceContext HDC) (HDC, error) {
	result, _, err := procCreateCompatibleDC.Call(uintptr(handleDeviceContext))
	if result == 0 {
		return HDC(result), fmt.Errorf("[CreateCompatibleDC] failed: %s", err)
	}
	return HDC(result), nil
}

// The CreateCompatibleBitmap function creates a bitmap compatible with the device that is associated with the specified device context.
func CreateCompatibleBitmap(handleDeviceContext HDC, bitmapWidth int32, bitmapHeight int32) (HBITMAP, error) {
	result, _, err := procCreateCompatibleBitmap.Call(uintptr(handleDeviceContext), uintptr(bitmapWidth), uintptr(bitmapHeight))
	if result == 0 {
		return HBITMAP(result), fmt.Errorf("[CreateCompatibleBitmap] failed: %s", err)
	}
	return HBITMAP(result), nil
}

// The SelectObject function selects an object into the specified device context (DC). The new object replaces the previous object of the same type.
func SelectObject(handleDeviceContext HDC, handleObject HGDIOBJ) (HGDIOBJ, error) {
	result, _, err := procSelectObject.Call(uintptr(handleDeviceContext), uintptr(handleObject))
	if result == 0 {
		return HGDIOBJ(result), fmt.Errorf("[SelectObject] failed: %s", err)
	}
	return HGDIOBJ(result), nil
}

// The BitBlt function performs a bit-block transfer of the color data corresponding to a rectangle of pixels from the specified source device context into a destination device context.
func BitBlt(
	destinationHandleDeviceContext HDC,
	destinationX, destinationY int32,
	width, height int32,
	sourceHandleDeviceContext HDC,
	sourceX, sourceY int32,
	rasterOperationCode uint,
) error {
	ret, _, _ := procBitBlt.Call(
		uintptr(destinationHandleDeviceContext),
		uintptr(destinationX),
		uintptr(destinationY),
		uintptr(width),
		uintptr(height),
		uintptr(sourceHandleDeviceContext),
		uintptr(sourceX),
		uintptr(sourceY),
		uintptr(rasterOperationCode),
	)
	if ret == 0 {
		return fmt.Errorf("[BitBlt] failed: %s", windows.GetLastError())
	}
	return nil
}

// The GetDIBits function retrieves the bits of the specified compatible bitmap and copies them into a buffer as a DIB using the specified format.
func GetDIBits(handleDeviceContext HDC, handleBitmap HBITMAP, firstScanLine uint32, scanLinesCount uint32, receiveBuffer *byte, bitmapInfo *BITMAPINFO, colorsFormat uint32) error {
	ret, _, _ := procGetDIBits.Call(
		uintptr(handleDeviceContext),
		uintptr(handleBitmap),
		uintptr(firstScanLine),
		uintptr(scanLinesCount),
		uintptr(unsafe.Pointer(receiveBuffer)),
		uintptr(unsafe.Pointer(bitmapInfo)),
		uintptr(colorsFormat),
	)

	if ret == 0 {
		return fmt.Errorf("[GetDIBits] failed: %s", windows.GetLastError())
	}
	return nil
}

// The DeleteDC function deletes the specified device context (DC).
func DeleteDC(handleDeviceContext HDC) error {
	ret, _, err := procDeleteDC.Call(uintptr(handleDeviceContext))
	if ret == 0 {
		return fmt.Errorf("[DeleteDC] failed: %s", err)
	}
	return nil
}

// The DeleteObject function deletes a logical pen, brush, font, bitmap, region, or palette, freeing all system resources associated with the object. After the object is deleted, the specified handle is no longer valid.
func DeleteObject(handleObject HGDIOBJ) error {
	ret, _, err := procDeleteObject.Call(uintptr(handleObject))
	if ret == 0 {
		return fmt.Errorf("[DeleteObject] failed: %s", err)
	}
	return nil
}
