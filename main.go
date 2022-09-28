package main

import (
	"fmt"
	"runtime"
	"unsafe"

	"syscall"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

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
var WS_EX_TOOLWINDOW int64 = 128

var hotKeyId int = 0

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	defer glfw.Terminate()

	mode := glfw.GetPrimaryMonitor().GetVideoMode()

	glfw.WindowHint(glfw.Floating, glfw.True)
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.AutoIconify, glfw.False)

	window, err := glfw.CreateWindow(mode.Width, mode.Height, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	monitorX, monitorY := glfw.GetPrimaryMonitor().GetPos()
	window.SetPos(monitorX, monitorY)

	window.SetAttrib(glfw.Resizable, glfw.True)
	window.SetAttrib(glfw.Decorated, glfw.False)
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	hwnd := window.GetWin32Window()
	window.Show()

	syscall.SyscallN(SetWindowPos, uintptr(unsafe.Pointer(hwnd)), uintptr(HWND_TOPMOST), 0, 0, 0, 0, uintptr(TOPMOST_FLAGS))

	syscall.SyscallN(ShowWindow, uintptr(unsafe.Pointer(hwnd)), uintptr(SW_HIDE))
	syscall.SyscallN(SetWindowLongPtrW, uintptr(unsafe.Pointer(hwnd)), uintptr(GWL_EXSTYLE), uintptr(WS_EX_TOOLWINDOW))
	syscall.SyscallN(ShowWindow, uintptr(unsafe.Pointer(hwnd)), uintptr(syscall.SW_SHOW))

	window.SetOpacity(0.8)

	RegisterGlobalHotkey(VK_F4, MOD_ALT|MOD_NOREPEAT, window)
	defer UnregisterGlobalHotkeys()
	// fmt.Scanln()

	window.SetKeyCallback(KeyCallback)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// Do OpenGL stuff.
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

func RegisterGlobalHotkey(key uint, mod uint, window *glfw.Window) {
	hotKeyId++
	_1, _2, err := syscall.SyscallN(RegisterHotkey, uintptr(unsafe.Pointer(nil)), uintptr(hotKeyId), uintptr(mod), uintptr(int16(key)))
	fmt.Println(_1, _2, err)
	fmt.Println(hotKeyId)
}

func UnregisterGlobalHotkeys() {
	for i := 1; i < hotKeyId; i++ {
		syscall.SyscallN(UnregisterHotkey, uintptr(i))
	}
}
