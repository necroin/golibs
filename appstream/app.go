//go:build cgo

package main

/*
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation
#include "cpp/window_darwin.h"


*/
import "C"
import (
	"fmt"
	"time"
)

func main() {
	// dysplaysCount := C.GetActiveDisplaysCount()
	// fmt.Println(dysplaysCount)
	// dysplays := C.GetActiveDisplays()
	// fmt.Println(C.IndexArray(dysplays, 0))
	// fmt.Println(C.IndexArray(dysplays, 1))
	// fmt.Println(getCoreGraphicsCoordinateOfDisplay(C.IndexArray(dysplays, 0)))
	// fmt.Println(getCoreGraphicsCoordinateOfDisplay(C.IndexArray(dysplays, 1)))

	windowsInfo := C.GetWindowsInfo()
	windowsCount := C.GetWindowsCount()
	for i := 0; i < int(windowsCount); i++ {
		windows := C.IndexWindowInfo(windowsInfo, C.long(i))
		fmt.Println(windows)

	}
	fmt.Println()
	time.Sleep(time.Second)

}

func getCoreGraphicsCoordinateOfDisplay(id C.CGDirectDisplayID) C.CGRect {
	main := C.CGDisplayBounds(C.CGMainDisplayID())
	r := C.CGDisplayBounds(id)
	return C.CGRectMake(r.origin.x, -r.origin.y-r.size.height+main.size.height,
		r.size.width, r.size.height)
}
