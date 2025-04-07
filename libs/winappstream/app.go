//go:build windows

package winappstream

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"unsafe"

	"github.com/necroin/golibs/utils/finalizer"
	"github.com/necroin/golibs/utils/promise"
	"github.com/necroin/golibs/utils/winapi"
	"github.com/necroin/golibs/utils/winutils"
	"golang.org/x/sys/windows"
)

const (
	encodedDataChannelLen = 1
)

type Cache struct {
	captureRect  windows.Rect
	bitmap       winapi.HBITMAP
	bitmapHeader winapi.BITMAPINFOHEADER
	memptr       unsafe.Pointer
	finalizer    *finalizer.Finalizer
}

func NewCache(desktopHDC winapi.HDC, desktopCompatibleHDC winapi.HDC, captureRect windows.Rect) (*Cache, error) {
	finalizer := finalizer.NewFinalizer()

	imageWidth := winutils.RectWidth(captureRect)
	imageHeight := winutils.RectHeight(captureRect)

	bitmap, err := winapi.CreateCompatibleBitmap(desktopHDC, imageWidth, imageHeight)
	if err != nil {
		return nil, fmt.Errorf("[NewCache] failed create compatible bitmap: %s", err)
	}
	finalizer.AddFunc(func() { winapi.DeleteObject(winapi.HGDIOBJ(bitmap)) })

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

	bitmapDataSize := uint32(((int64(imageWidth)*int64(bitmapHeader.BiBitCount) + 31) / 32) * 4 * int64(imageHeight))
	memptr, err := windows.LocalAlloc(windows.LMEM_FIXED, bitmapDataSize)
	if err != nil {
		return nil, fmt.Errorf("[NewCache] failed LocalAlloc: %s", err)
	}
	finalizer.AddFunc(func() { windows.LocalFree(windows.Handle(memptr)) })

	return &Cache{
		captureRect:  captureRect,
		bitmap:       bitmap,
		bitmapHeader: bitmapHeader,
		memptr:       unsafe.Pointer(memptr),
		finalizer:    finalizer,
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
	ctx                  context.Context
	cancel               context.CancelFunc
	finalizer            *finalizer.Finalizer
}

func NewApp(pid winapi.ProcessId) (*App, error) {
	finalizer := finalizer.NewFinalizer()

	windowHandles := winutils.GetWindowHandlesByProcessId(pid)
	if len(windowHandles) == 0 {
		return nil, errors.New("[NewApp] process has 0 window handles")
	}

	desktopHWND := winapi.GetDesktopWindow()
	desktopHDC, err := winapi.GetWindowDC(desktopHWND)
	if err != nil {
		return nil, fmt.Errorf("[NewApp] failed get desktop device context: %s", err)
	}
	finalizer.AddFunc(func() { winapi.ReleaseDC(desktopHWND, desktopHDC) })

	desktopCompatibleHDC, err := winapi.CreateCompatibleDC(desktopHDC)
	if err != nil {
		return nil, fmt.Errorf("[NewApp] failed create compatible device context: %s", err)
	}
	finalizer.AddFunc(func() { winapi.DeleteDC(desktopCompatibleHDC) })

	captureRect, err := winutils.GetCaptureRectByHandles(windowHandles)
	if err != nil {
		return nil, fmt.Errorf("[NewApp] failed get capture rect: %s", err)
	}

	cache, err := NewCache(desktopHDC, desktopCompatibleHDC, captureRect)
	if err != nil {
		return nil, fmt.Errorf("[NewApp] failed create cache: %s", err)
	}
	finalizer.AddFunc(func() { cache.Destroy() })

	encodedData := make(chan *promise.Promise[image.Image, []byte], encodedDataChannelLen)
	finalizer.AddFunc(func() { close(encodedData) })

	ctx, cancel := context.WithCancel(context.Background())
	finalizer.AddFunc(func() { cancel() })

	return &App{
		pid:                  pid,
		windowHandles:        windowHandles,
		desktopHWND:          desktopHWND,
		desktopHDC:           desktopHDC,
		desktopCompatibleHDC: desktopCompatibleHDC,
		cache:                cache,
		encodedData:          encodedData,
		ctx:                  ctx,
		cancel:               cancel,
		finalizer:            finalizer,
	}, nil
}

func (app *App) Destroy() {
	if app.finalizer != nil {
		app.finalizer.Execute()
		app.finalizer = nil
	}
}

func (cache *Cache) Destroy() {
	if cache.finalizer != nil {
		cache.finalizer.Execute()
		cache.finalizer = nil
	}
}

func (app *App) UpdateHandles() {
	app.windowHandles = winutils.GetWindowHandlesByProcessId(app.pid)
}

func (app *App) CaptureImageScreenVersion() (image.Image, error) {
	app.UpdateHandles()
	captureRect, err := winutils.GetCaptureRectByHandles(app.windowHandles)
	if err != nil {
		return nil, fmt.Errorf("[CaptureImageScreenVersion] failed get capture rect: %s", err)
	}

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
		return nil, fmt.Errorf("[CaptureImageScreenVersion] failed bit blt: %s", err)
	}

	if err := winapi.GetDIBits(app.desktopCompatibleHDC, app.cache.bitmap, 0, uint32(imageHeight), (*uint8)(app.cache.memptr), (*winapi.BITMAPINFO)(unsafe.Pointer(&app.cache.bitmapHeader)), winapi.DIB_RGB_COLORS); err != nil {
		return nil, fmt.Errorf("[CaptureImageScreenVersion] failed GetDIBits: %s", err)
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
			select {
			case <-app.ctx.Done():
				return
			default:
				if len(app.encodedData) == encodedDataChannelLen {
					continue
				}

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
		}
	}()
}
