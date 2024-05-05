package winutils

import (
	"errors"
	"fmt"

	"github.com/necroin/golibs/utils/winapi"
	"golang.org/x/sys/windows"
)

func RectWidth(value windows.Rect) int32 {
	return value.Right - value.Left
}

func RectHeight(value windows.Rect) int32 {
	return value.Bottom - value.Top
}

func RectSize(value windows.Rect) int32 {
	return RectWidth(value) * RectHeight(value)
}

func RectEqual(value1 windows.Rect, value2 windows.Rect) bool {
	return value1.Left == value2.Left && value1.Top == value2.Top && value1.Right == value2.Right && value1.Bottom == value2.Bottom
}

func IsValidRect(rect windows.Rect) bool {
	return !(RectWidth(rect) == 0 || RectHeight(rect) == 0)
}

func FindLargestRect(values []windows.Rect) (int, error) {
	if len(values) == 0 {
		return 0, errors.New("no rects")
	}
	result := 0
	for i := 1; i < len(values); i++ {
		if RectSize(values[i]) > RectSize(values[result]) {
			result = i
		}
	}
	return result, nil
}

func ScreenToClientRect(handleWindow windows.HWND, value windows.Rect) windows.Rect {
	MinPoint := winapi.POINT{X: value.Left, Y: value.Top}
	MaxPoint := winapi.POINT{X: value.Right, Y: value.Bottom}

	MinPoint, _ = winapi.ScreenToClient(handleWindow, MinPoint)
	MaxPoint, _ = winapi.ScreenToClient(handleWindow, MaxPoint)

	value.Left = MinPoint.X
	value.Top = MinPoint.Y
	value.Right = MaxPoint.X
	value.Bottom = MaxPoint.Y

	return value
}

func GetWindowHandlesRects(handles []windows.HWND) []windows.Rect {
	result := []windows.Rect{}
	for _, handle := range handles {
		rect, _ := winapi.GetWindowRect(handle)
		result = append(result, rect)
	}
	return result
}

func GetWindowHandlesClientRects(handles []windows.HWND) []windows.Rect {
	result := []windows.Rect{}
	for _, handle := range handles {
		rect, _ := winapi.GetClientRect(handle)
		result = append(result, rect)
	}
	return result
}

func GetCaptureRect(windowRect windows.Rect, clientRect windows.Rect) windows.Rect {
	diffX := (RectWidth(windowRect) - RectWidth(clientRect)) / 2
	diffY := (RectHeight(windowRect) - RectHeight(clientRect)) / 2

	captureRect := windowRect
	captureRect.Left += diffX
	captureRect.Right -= diffX
	captureRect.Top += diffY / 2
	captureRect.Bottom -= diffY / 2

	return captureRect
}

func GetCaptureRectByHandles(windowHandles []windows.HWND) (windows.Rect, error) {
	windowRects := GetWindowHandlesRects(windowHandles)
	clientRects := GetWindowHandlesClientRects(windowHandles)

	largestRectIndex, err := FindLargestRect(windowRects)
	if err != nil {
		return windows.Rect{}, fmt.Errorf("[GetCaptureRect] failed find largest rect index:%s", err)
	}

	largestWindowRect := windowRects[largestRectIndex]
	largestClientRect := clientRects[largestRectIndex]
	return GetCaptureRect(largestWindowRect, largestClientRect), nil
}
