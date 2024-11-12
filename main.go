package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Bouncing Square", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer func(window *sdl.Window) {
		err := window.Destroy()
		if err != nil {
		}
	}(window)

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	// Initial position, size, and velocity of the square
	square := sdl.Rect{X: 300, Y: 200, W: 50, H: 50}
	velocityX, velocityY := int32(5), int32(5)

	running := true
	for running {
		// Handle events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			}
		}

		// Update square position
		square.X += velocityX
		square.Y += velocityY

		// Check for collisions with window edges
		if square.X <= 0 || square.X+square.W >= 800 {
			velocityX = -velocityX // Reverse horizontal direction
		}
		if square.Y <= 0 || square.Y+square.H >= 600 {
			velocityY = -velocityY // Reverse vertical direction
		}

		// Clear the surface by filling it with black
		err = surface.FillRect(nil, 0)
		if err != nil {
			panic(err)
		}

		// Draw the square with a color (purple)
		color := sdl.Color{R: 255, G: 0, B: 255, A: 255} // purple
		pixel := sdl.MapRGBA(surface.Format, color.R, color.G, color.B, color.A)
		err = surface.FillRect(&square, pixel)
		if err != nil {
			panic(err)
		}

		// Update the window surface to show changes
		err = window.UpdateSurface()
		if err != nil {
			panic(err)
		}

		// Delay to control frame rate
		framerate := (1.0 / 60.0) * 1000.0
		println(uint32(framerate))
		sdl.Delay(uint32(framerate)) // roughly 60 FPS
	}
}
