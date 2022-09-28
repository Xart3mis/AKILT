package main

//TODO: REFACTOR EVERYTHING

import (
	"fmt"
	"log"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-vgo/robotgo"
	"github.com/nullboundary/glfont"
)

type MousePos struct{ xpos, ypos float64 }

var u32dll, _ = syscall.LoadLibrary("user32.dll")
var ShowWindow, _ = syscall.GetProcAddress(u32dll, "ShowWindow")
var SetWindowPos, _ = syscall.GetProcAddress(u32dll, "SetWindowPos")
var RegisterHotkey, _ = syscall.GetProcAddress(u32dll, "RegisterHotKey")
var UnregisterHotkey, _ = syscall.GetProcAddress(u32dll, "UnregisterHotkey")
var SetWindowLongPtrW, _ = syscall.GetProcAddress(u32dll, "SetWindowLongPtrW")

var MOD_ALT uint = 0x0001
var MOD_CONTROL uint = 0x0002
var MOD_NOREPEAT uint = 0x4000
var MOD_SHIFT uint = 0x0004
var MOD_WIN uint = 0x0008

var VK_F4 uint = 0x73

var HWND_TOPMOST int = -1
var SWP_NOSIZE int = 0x0001
var SWP_NOMOVE int = 0x0002
var TOPMOST_FLAGS int = SWP_NOMOVE | SWP_NOSIZE

var SW_HIDE int = 0
var GWL_EXSTYLE int = -20
var WS_EX_TOOLWINDOW uint = 0x00000080
var WS_EX_LAYERED uint = 0x00080000
var WS_EX_TRANSPARENT uint = 0x00000020

var hotKeyId int = 0

var CurrentMousePosition MousePos = MousePos{0, 0}

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
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
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	font, err := glfont.LoadFont("./MedusaGothic.ttf", int32(52), mode.Width, mode.Height)
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	hwnd := window.GetWin32Window()
	window.Show()

	syscall.SyscallN(SetWindowPos, uintptr(unsafe.Pointer(hwnd)), uintptr(HWND_TOPMOST), 0, 0, 0, 0, uintptr(TOPMOST_FLAGS))

	syscall.SyscallN(ShowWindow, uintptr(unsafe.Pointer(hwnd)), uintptr(SW_HIDE))
	syscall.SyscallN(SetWindowLongPtrW, uintptr(unsafe.Pointer(hwnd)), uintptr(GWL_EXSTYLE), uintptr(WS_EX_TOOLWINDOW|WS_EX_TRANSPARENT|WS_EX_LAYERED))
	syscall.SyscallN(ShowWindow, uintptr(unsafe.Pointer(hwnd)), uintptr(syscall.SW_SHOW))

	window.SetOpacity(0.8)

	RegisterGlobalHotkey(VK_F4, MOD_ALT|MOD_NOREPEAT, window)
	defer UnregisterGlobalHotkeys()

	window.SetMouseButtonCallback(MouseButtonCallback)
	window.SetCursorPosCallback(CursorPosCallback)
	window.SetKeyCallback(KeyCallback)
	window.SetScrollCallback(ScrollCallback)

	text := "You've Been Pwned"

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(0, 0, 0, 0)

		font.SetColor(1.0, 1.0, 1.0, 1.0)
		font.Printf(float32(mode.Width)/2-font.Width(1.0, text)/2, float32(mode.Height)/2, 1.0, text)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press || action == glfw.Repeat {
		fmt.Println(glfw.GetKeyScancode(key))
		if key == glfw.KeyQ && mods == glfw.ModControl+glfw.ModAlt+glfw.ModShift {
			w.SetShouldClose(true)
		}
	}
}

func ScrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	robotgo.Scroll(int(xoff), int(yoff))
	fmt.Println(xoff, yoff)
}

func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == 0 && button == 0 {
		// robotgo.Toggle("left", "up")
		fmt.Println("Left Click Release")
	}
	if action == 1 && button == 1 {
		// robotgo.Toggle("right", "down")
		fmt.Println("Right Click Press")
	}
	if action == 0 && button == 1 {
		// robotgo.Toggle("right", "up")
		fmt.Println("Right Click Release")
	}
	if action == 1 && button == 0 {
		// robotgo.Toggle("left", "down")
		fmt.Println("Left Click Press")
	}
}

func CursorPosCallback(w *glfw.Window, xpos float64, ypos float64) {
	CurrentMousePosition.xpos = xpos
	CurrentMousePosition.ypos = ypos
}

func RegisterGlobalHotkey(key uint, mod uint, window *glfw.Window) {
	hotKeyId++
	_, _, err := syscall.SyscallN(RegisterHotkey, uintptr(unsafe.Pointer(nil)), uintptr(hotKeyId), uintptr(mod), uintptr(int16(key)))
	fmt.Println(err)
}

func UnregisterGlobalHotkeys() {
	for i := 1; i < hotKeyId; i++ {
		syscall.SyscallN(UnregisterHotkey, uintptr(i))
	}
}

func PollHotkey() {

}
