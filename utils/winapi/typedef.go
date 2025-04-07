//go:build windows

package winapi

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type (
	ProcessId uint32
	HANDLE    windows.Handle
	HDC       HANDLE
	HBITMAP   HANDLE
	HGDIOBJ   HANDLE
	HGLOBAL   HANDLE
	HRESULT   HANDLE
	LPVOID    unsafe.Pointer
)

type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

type RGBQUAD struct {
	RgbBlue     byte
	RgbGreen    byte
	RgbRed      byte
	RgbReserved byte
}

type BITMAPINFO struct {
	BmiHeader BITMAPINFOHEADER
	BmiColors *RGBQUAD
}

type POINT struct {
	X int32
	Y int32
}
