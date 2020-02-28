package main

import (
	"fmt"
	"os"
	"time"

	_ "golang.org/x/image/bmp"

	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"local/eyes"
)

type mode int

const (
	pictureMode mode = iota
	pictureAndMarkerMode
	markerMode
	barsMode
	numberOfModes
)

func run() {
	if len(os.Args) != 3 {
		fmt.Printf("usage: %s <pic_path> <session_path>\n", os.Args[0])
		return
	}

	player, err := eyes.LoadPlayer(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	w, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Eyes Shower",
		Icon:   nil,
		Bounds: player.Picture.Bounds(),
		VSync:  true,
	})
	if err != nil {
		panic(err)
	}

	var (
		paused   = true
		slowdown = 0
		mode     = pictureMode
	)

	marker := eyes.NewMarker(
		int(player.Picture.Bounds().W()/eyes.DefaultMarkerScale),
		int(player.Picture.Bounds().H()/eyes.DefaultMarkerScale),
		eyes.DefaultMarkerRadius/eyes.DefaultMarkerScale,
		eyes.DefaultActivationTime,
	)

	bars := eyes.NewBars(5, 2)

	const step = time.Second / 512
	ticker := &UniformTicker{Step: step, Action: func() {
		if player.Finished() {
			return
		}
		player.Update(step)
		middle := player.LeftEye().Add(player.RightEye()).Scaled(0.5)
		marker.PointAt(middle.Scaled(1.0/eyes.DefaultMarkerScale), step)
	}}

	last := time.Now()
	for !w.Closed() {
		dt := time.Since(last)
		last = time.Now()

		if w.JustPressed(pixelgl.KeySpace) {
			paused = !paused
		}
		if w.JustPressed(pixelgl.KeyEqual) {
			if slowdown > 0 {
				slowdown--
			}
		}
		if w.JustPressed(pixelgl.KeyMinus) {
			slowdown++
		}
		if w.JustPressed(pixelgl.KeyM) {
			mode = (mode + 1) % numberOfModes
		}

		if !paused {
			for i := 0; i < slowdown; i++ {
				dt /= 2
			}
			ticker.Update(dt)
		}

		w.SetTitle(fmt.Sprintf(
			"Eyes Shower | Time: %v | Paused: %v | Slowdown: %d | Fixations: %d",
			player.Session[player.Position].Time.Truncate(time.Second/100),
			paused,
			slowdown,
			marker.NumFixations(),
		))

		w.Clear(colornames.White)
		switch mode {
		case pictureMode:
			player.Draw(w, true)
		case pictureAndMarkerMode:
			player.Draw(w, true)
			marker.Draw(w, player.Picture.Bounds(), colornames.Green)
		case markerMode:
			player.Draw(w, false)
			marker.Draw(w, player.Picture.Bounds(), colornames.Green)
		case barsMode:
			player.Draw(w, false)
			bars.Draw(marker, w, player.Picture.Bounds(), colornames.Blue)
		}
		w.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
