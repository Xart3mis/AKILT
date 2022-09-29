package main

//TODO: REFACTOR EVERYTHING

import (
	_ "embed"
	"log"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/Xart3mis/GoHkar/lib/bundles"
	"github.com/Xart3mis/GoHkar/lib/reg"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/nullboundary/glfont"
	"golang.design/x/hotkey"
)

var u32dll, _ = syscall.LoadLibrary("user32.dll")
var ShowWindow, _ = syscall.GetProcAddress(u32dll, "ShowWindow")
var SetWindowPos, _ = syscall.GetProcAddress(u32dll, "SetWindowPos")
var SetWindowLongPtrW, _ = syscall.GetProcAddress(u32dll, "SetWindowLongPtrW")

var HWND_TOPMOST int = -1
var SWP_NOSIZE int = 0x0001
var SWP_NOMOVE int = 0x0002
var TOPMOST_FLAGS int = SWP_NOMOVE | SWP_NOSIZE

var SW_HIDE int = 0
var GWL_EXSTYLE int = -20
var WS_EX_TOOLWINDOW uint = 0x00000080
var WS_EX_LAYERED uint = 0x80000
var WS_EX_TRANSPARENT uint = 0x20

func init() {
	runtime.LockOSThread()
}

func main() {
	bundles.WriteFiraCodeNerd()
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	defer glfw.Terminate()

	mode := glfw.GetPrimaryMonitor().GetVideoMode()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	glfw.WindowHint(glfw.RedBits, mode.RedBits)
	glfw.WindowHint(glfw.GreenBits, mode.GreenBits)
	glfw.WindowHint(glfw.BlueBits, mode.BlueBits)
	glfw.WindowHint(glfw.RefreshRate, mode.RefreshRate)

	glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True)
	glfw.WindowHint(glfw.AutoIconify, glfw.False)
	glfw.WindowHint(glfw.Floating, glfw.True)
	glfw.WindowHint(glfw.Visible, glfw.False)

	window, err := glfw.CreateWindow(mode.Width, mode.Height, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	monitorX, monitorY := glfw.GetPrimaryMonitor().GetPos()
	window.SetPos(monitorX, monitorY)

	window.SetAttrib(glfw.Resizable, glfw.False)
	window.SetAttrib(glfw.Decorated, glfw.False)

	go func() {
		hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModAlt, hotkey.ModShift}, hotkey.KeyX)
		if err := hk.Register(); err != nil {
			panic("hotkey registration failed")
		}

		for range hk.Keydown() {
			window.SetShouldClose(true)
		}
	}()

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	font, err := glfont.LoadFont("FiraCodeNerd.ttf", int32(52), mode.Width, mode.Height)
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}

	hwnd := window.GetWin32Window()
	window.Show()

	syscall.SyscallN(SetWindowPos, uintptr(unsafe.Pointer(hwnd)), uintptr(HWND_TOPMOST), 0, 0, 0, 0, uintptr(TOPMOST_FLAGS))

	syscall.SyscallN(ShowWindow, uintptr(unsafe.Pointer(hwnd)), uintptr(SW_HIDE))
	syscall.SyscallN(SetWindowLongPtrW, uintptr(unsafe.Pointer(hwnd)), uintptr(GWL_EXSTYLE), uintptr(WS_EX_TOOLWINDOW|WS_EX_TRANSPARENT|WS_EX_LAYERED))
	syscall.SyscallN(ShowWindow, uintptr(unsafe.Pointer(hwnd)), uintptr(syscall.SW_SHOW))

	window.SetOpacity(0.8)
	window.SetCloseCallback(CloseCallback)

	text := ".. .. .. .. .."

	reg.GetUniqueSystemId()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(0, 0, 0, 0)

		font.SetColor(1.0, 1.0, 1.0, 1.0)
		font.Printf(float32(mode.Width)/2-font.Width(1.0, text)/2, float32(mode.Height)/2, 1.0, text)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func CloseCallback(w *glfw.Window) {
	w.SetShouldClose(false)
}
