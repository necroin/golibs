package winappstream

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"unsafe"

	"github.com/necroin/golibs/utils/promise"
	"github.com/necroin/golibs/utils/winapi"
	"github.com/necroin/golibs/utils/winutils"
	"golang.org/x/sys/windows"
)

type Cache struct {
	captureRect  windows.Rect
	bitmap       winapi.HBITMAP
	bitmapHeader winapi.BITMAPINFOHEADER
	hmem         winapi.HGLOBAL
	memptr       unsafe.Pointer
}

func NewCache(desktopHDC winapi.HDC, desktopCompatibleHDC winapi.HDC, captureRect windows.Rect) (*Cache, error) {
	imageWidth := winutils.RectWidth(captureRect)
	imageHeight := winutils.RectHeight(captureRect)

	bitmap, err := winapi.CreateCompatibleBitmap(desktopHDC, imageWidth, imageHeight)
	if err != nil {
		return nil, fmt.Errorf("[NewCache] failed create compatible bitmap: %s", err)
	}

	if _, err := winapi.SelectObject(desktopCompatibleHDC, winapi.HGDIOBJ(bitmap)); err != nil {
		return nil, fmt.Errorf("[NewCache] failed select bitmap: %s", err)
	}

	bitmapHeader := winapi.BITMAPINFOHEADER{}
	bitmapHeader.BiSize = uint32(unsafe.Sizeof(bitmapHeader))
	bitmapHeader.BiPlanes = 1
	bitmapHeader.BiBitCount = 32
	bitmapHeader.BiWidth = imageWidth
	bitmapHeader.BiHeight = -imageHeight
	bitmapHeader.BiCompression = winapi.BI_RGB
	bitmapHeader.BiSizeImage = 0

	bitmapDataSize := uintptr(((int64(imageWidth)*int64(bitmapHeader.BiBitCount) + 31) / 32) * 4 * int64(imageHeight))
	hmem, err := winapi.GlobalAlloc(winapi.GMEM_MOVEABLE, bitmapDataSize)
	if err != nil {
		return nil, fmt.Errorf("[NewCache] failed GlobalAlloc: %s", err)
	}

	memptr, err := winapi.GlobalLock(hmem)
	if err != nil {
		return nil, fmt.Errorf("[NewCache] failed GlobalLock: %s", err)
	}

	return &Cache{
		captureRect:  captureRect,
		bitmap:       bitmap,
		bitmapHeader: bitmapHeader,
		hmem:         hmem,
		memptr:       memptr,
	}, nil
}

type App struct {
	pid                  winapi.ProcessId
	windowHandles        []windows.HWND
	desktopHWND          windows.HWND
	desktopHDC           winapi.HDC
	desktopCompatibleHDC winapi.HDC
	cache                *Cache
	encodedData          chan *promise.Promise[image.Image, []byte]
}

func NewApp(pid winapi.ProcessId) (*App, error) {
	windowHandles := winutils.GetWindowHandlesByProcessId(pid)

	desktopHWND := winapi.GetDesktopWindow()
	desktopHDC, err := winapi.GetWindowDC(desktopHWND)
	if err != nil {
		return nil, fmt.Errorf("[NewApp] failed get desktop device context: %s", err)
	}

	desktopCompatibleHDC, err := winapi.CreateCompatibleDC(desktopHDC)
	if err != nil {
		return nil, fmt.Errorf("[NewApp] failed create compatible device context: %s", err)
	}

	captureRect := winutils.GetCaptureRect(windowHandles)
	cache, err := NewCache(desktopHDC, desktopCompatibleHDC, captureRect)
	if err != nil {
		return nil, fmt.Errorf("[NewApp] failed create cache: %s", err)
	}

	return &App{
		pid:                  pid,
		windowHandles:        windowHandles,
		desktopHWND:          desktopHWND,
		desktopHDC:           desktopHDC,
		desktopCompatibleHDC: desktopCompatibleHDC,
		cache:                cache,
		encodedData:          make(chan *promise.Promise[image.Image, []byte], 10),
	}, nil
}

func (app *App) Destroy() {
	defer winapi.ReleaseDC(app.desktopHWND, app.desktopHDC)
	defer winapi.DeleteDC(app.desktopCompatibleHDC)
	defer app.cache.Destroy()
}

func (cache *Cache) Destroy() {
	defer winapi.DeleteObject(winapi.HGDIOBJ(cache.bitmap))
	defer winapi.GlobalUnlock(cache.hmem)
	defer winapi.GlobalFree(cache.hmem)
}

func (app *App) CaptureImageScreenVersion() (image.Image, error) {
	captureRect := winutils.GetCaptureRect(app.windowHandles)
	if !winutils.RectEqual(captureRect, app.cache.captureRect) {
		newCache, err := NewCache(app.desktopHDC, app.desktopCompatibleHDC, captureRect)
		if err != nil {
			return nil, fmt.Errorf("[CaptureImageScreenVersion] failed update cache: %s", err)
		}
		app.cache.Destroy()
		app.cache = newCache
	}
	imageWidth := winutils.RectWidth(captureRect)
	imageHeight := winutils.RectHeight(captureRect)

	if err := winapi.BitBlt(app.desktopCompatibleHDC, 0, 0, imageWidth, imageHeight, app.desktopHDC, captureRect.Left, captureRect.Top, winapi.SRCCOPY|winapi.CAPTUREBLT); err != nil {
		return nil, fmt.Errorf("failed bit blt: %s", err)
	}

	if err := winapi.GetDIBits(app.desktopCompatibleHDC, app.cache.bitmap, 0, uint32(imageHeight), (*uint8)(app.cache.memptr), (*winapi.BITMAPINFO)(unsafe.Pointer(&app.cache.bitmapHeader)), winapi.DIB_RGB_COLORS); err != nil {
		return nil, fmt.Errorf("failed GetDIBits: %s", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, int(imageWidth), int(imageHeight)))

	i := 0
	src := uintptr(app.cache.memptr)
	for y := int32(0); y < imageHeight; y++ {
		for x := int32(0); x < imageWidth; x++ {
			B := *(*uint8)(unsafe.Pointer(src))
			G := *(*uint8)(unsafe.Pointer(src + 1))
			R := *(*uint8)(unsafe.Pointer(src + 2))

			img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = R, G, B, 255

			i += 4
			src += 4
		}
	}

	return img, nil
}

func (app *App) HttpImageCaptureHandler() HttpImageCaptureHandler {
	return NewHttpImageCaptureHandler(app)
}

func (app *App) LaunchStream() {
	go func() {
		for {
			img, err := app.CaptureImageScreenVersion()
			if img == nil || err != nil {
				continue
			}
			app.encodedData <- promise.NewPromise[image.Image, []byte](img, func(img image.Image) ([]byte, error) {
				buf := &bytes.Buffer{}
				err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 75})
				return buf.Bytes(), err
			})
		}

	}()
}

// func (app *App) CaptureImage() (*image.RGBA, error) {
// 	largestWindowRect := app.captureRect
// 	largestWidth := largestWindowRect.Right - largestWindowRect.Left
// 	largestHeight := largestWindowRect.Bottom - largestWindowRect.Top

// 	resultImage := image.NewRGBA(image.Rect(0, 0, int(largestWidth), int(largestHeight)))

// 	for handleIndex, handle := range app.windowHandles {
// 		fmt.Printf("Hanle index: %d\n", handleIndex)

// 		handleRect := app.handlesRects[handleIndex]
// 		handleClientRect := app.handlesClientRects[handleIndex]
// 		fmt.Printf("Hanle Rect: %v\n", handleRect)
// 		fmt.Printf("Hanle Client Rect: %v\n", handleClientRect)

// 		handleWidth := handleRect.Right - handleRect.Left
// 		handleHeight := handleRect.Bottom - handleRect.Top

// 		handleClientWidth := handleClientRect.Right - handleClientRect.Left
// 		handleClientHeight := handleClientRect.Bottom - handleClientRect.Top

// 		diffX := handleWidth - handleClientWidth
// 		diffY := handleHeight - handleClientHeight
// 		fmt.Printf("Diff X: %d\n", diffX)
// 		fmt.Printf("Diff Y: %d\n", diffY)

// 		handleImage, err := CaptureWindowImage(handle, handleRect.Left, handleRect.Top, largestWidth, largestHeight, diffX, diffY, largestWidth, largestHeight)

// 		if err == nil {
// 			outputFile, _ := os.Create(fmt.Sprintf("test%d.png", handleIndex))
// 			png.Encode(outputFile, handleImage)
// 			outputFile.Close()
// 			for x := handleImage.Bounds().Min.X; x < handleImage.Bounds().Max.X; x++ {
// 				for y := handleImage.Bounds().Min.Y; y < handleImage.Bounds().Max.Y; y++ {
// 					handlePixelColor := handleImage.RGBAAt(x, y)
// 					if handlePixelColor.R == 0 && handlePixelColor.G == 0 && handlePixelColor.B == 0 {
// 						continue
// 					}
// 					resultImageColor := resultImage.RGBAAt(x, y)
// 					if resultImageColor.R == 0 && resultImageColor.G == 0 && resultImageColor.B == 0 {
// 						resultImage.SetRGBA(x, y, handlePixelColor)
// 					}
// 				}
// 			}
// 		}

// 	}

// 	outputFile, _ := os.Create("test.png")
// 	png.Encode(outputFile, resultImage)
// 	outputFile.Close()
// 	return resultImage, nil
// }
