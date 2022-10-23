package webcam

import (
	"gocv.io/x/gocv"
)

// func CaptureWebcam() []byte {
// 	bundles.WriteDSGrab()
// 	f, err := os.MkdirTemp("", "Pic_dir-")

// 	if err != nil {
// 		log.Println(err)
// 	}

// 	path := f + "\\Temp_Pic.jpg"

// 	cmd := exec.Command("./DSGrab.exe", "-r", "80x50", path)
// 	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
// 	cmd.Run()

// 	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	var b []byte = make([]byte, 6000000)
// 	file.Read(b)

// 	return b
// }

func CaptureWebcam() ([]byte, error) {
	webcam, _ := gocv.VideoCaptureDevice(0)
	img := gocv.NewMat()

	defer webcam.Close()
	defer img.Close()

	for i := 0; i < 10; i++ {
		webcam.Read(&img)
	}
	webcam.Read(&img)

	buf, err := gocv.IMEncode(gocv.PNGFileExt, img)

	if err != nil {
		return nil, err
	}

	return buf.GetBytes(), nil
}
