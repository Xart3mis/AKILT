package main

import (
	"runtime"
	"unsafe"

	"syscall"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	var u32dll, _ = syscall.LoadLibrary("user32.dll")
	var ShowWindow, _ = syscall.GetProcAddress(u32dll, "ShowWindow")
	var SetWindowLongPtrW, _ = syscall.GetProcAddress(u32dll, "SetWindowLongPtrW")

	var SW_HIDE int = 0
	var GWL_EXSTYLE int = -20
	var WS_EX_TOOLWINDOW int64 = 128

	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	defer glfw.Terminate()

	// mode := glfw.GetPrimaryMonitor().GetVideoMode()

	glfw.WindowHint(glfw.Floating, glfw.True)
	glfw.WindowHint(glfw.Visible, glfw.False)
	window, err := glfw.CreateWindow(1366, 768, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	monitorX, monitorY := glfw.GetPrimaryMonitor().GetPos()

	window.SetPos(monitorX, monitorY)

	window.SetAttrib(glfw.Resizable, glfw.True)
	window.SetAttrib(glfw.Decorated, glfw.False)
	window.MakeContextCurrent()

	hwnd := window.GetWin32Window()
	glfw.GetCurrentContext()
	window.Show()

	syscall.SyscallN(uintptr(ShowWindow), uintptr(unsafe.Pointer(hwnd)), uintptr(SW_HIDE))
	syscall.SyscallN(uintptr(SetWindowLongPtrW), uintptr(unsafe.Pointer(hwnd)), uintptr(GWL_EXSTYLE), uintptr(WS_EX_TOOLWINDOW))
	syscall.SyscallN(uintptr(ShowWindow), uintptr(unsafe.Pointer(hwnd)), uintptr(syscall.SW_SHOW))
	window.SetOpacity(0.8)

	for !window.ShouldClose() {

		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
