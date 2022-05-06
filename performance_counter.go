package main

import (
	"fmt"
	"time"

	"github.com/inkyblackness/imgui-go"
)

type PerformanceCounter struct {
	lastAddTime   time.Time
	lastFPS       float64
	lastFrameTime float64

	frameTimes     [200]float32
	frameTimeIndex int

	fpsSum, ftSum float64
	count         int

	doingVsync bool
	SlowDown   bool

	frame uint

	gofast bool
}

func (pc *PerformanceCounter) DoCount() {
	pc.frame++

	elapsed := time.Since(pc.lastAddTime)

	pc.lastFrameTime = float64(elapsed.Microseconds()) / 1000
	pc.lastFPS = 1000.0 / (pc.lastFrameTime)

	pc.PushFT(pc.lastFrameTime)

	pc.fpsSum += pc.lastFPS
	pc.ftSum += pc.lastFrameTime
	pc.count++

	if pc.count >= 5000 {
		pc.count = 0
		pc.fpsSum = 0
		pc.ftSum = 0
	}

	pc.lastAddTime = time.Now()
}

func (pc *PerformanceCounter) PushFT(ft float64) {
	pc.frameTimes[pc.frameTimeIndex] = float32(ft)
	pc.frameTimeIndex++
	pc.frameTimeIndex = pc.frameTimeIndex % 200
}

func (pc *PerformanceCounter) DrawUI() {
	imgui.BeginV("Performance", &PerformanceShown, 0)
	imgui.Text(fmt.Sprintf("Frame: %v", pc.frame))
	imgui.Text(fmt.Sprintf("Frame Time (ms): %.2f", pc.ftSum/float64(pc.count)))
	imgui.Text(fmt.Sprintf("FPS avg: %.2f", pc.fpsSum/float64(pc.count)))

	imgui.PlotHistogramV("## Frametimes", pc.frameTimes[:], pc.frameTimeIndex, "Frame Time", 0, 30, imgui.Vec2{X: 0, Y: 100})

	if imgui.Checkbox("VSync", &pc.doingVsync) {
		Game.win.SetVSync(pc.doingVsync)
		pc.count = 0
		pc.fpsSum = 0
		pc.ftSum = 0
	}

	imgui.Checkbox("Slow down", &pc.SlowDown)

	imgui.Checkbox("Turbo", &pc.gofast)
	if pc.gofast {
		UpdatesPerFrame = 20
	} else {
		UpdatesPerFrame = 1
	}

	imgui.End()

	if pc.SlowDown {
		time.Sleep(200 * time.Millisecond)
	}
}
