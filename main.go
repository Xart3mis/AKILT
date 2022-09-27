package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

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

	glfw.WindowHint(glfw.RedBits, mode.RedBits)
	glfw.WindowHint(glfw.BlueBits, mode.BlueBits)
	glfw.WindowHint(glfw.GreenBits, mode.GreenBits)
	glfw.WindowHint(glfw.RefreshRate, mode.RefreshRate)
	glfw.WindowHint(glfw.Decorated, glfw.False)
	glfw.WindowHint(glfw.Floating, glfw.True)


	window, err := glfw.CreateWindow(mode.Width, mode.Height, "Testing", glfw.GetPrimaryMonitor(), nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.SetOpacity(0.5)

	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	for !window.ShouldClose() {

		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
