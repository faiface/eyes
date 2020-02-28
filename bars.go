package eyes

import (
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type Bars struct {
	BarWidth         float64
	SpaceBetweenBars float64

	imd *imdraw.IMDraw
}

func NewBars(barWidth, spaceBetweenBars float64) *Bars {
	return &Bars{
		BarWidth:         barWidth,
		SpaceBetweenBars: spaceBetweenBars,
		imd:              imdraw.New(nil),
	}
}

func (b *Bars) Draw(marker *Marker, t pixel.Target, area pixel.Rect, clr color.Color) {
	times := marker.TotalTimesOnFixations()

	var maxTime time.Duration
	for _, t := range times {
		if t > maxTime {
			maxTime = t
		}
	}

	b.imd.Clear()
	for i := range times {
		x := area.Min.X + b.SpaceBetweenBars + float64(i)*(b.BarWidth+b.SpaceBetweenBars)
		y := area.Min.Y + b.SpaceBetweenBars
		h := (area.H() - 2*b.SpaceBetweenBars) * float64(times[i]) / float64(maxTime)
		b.imd.Color = clr
		b.imd.Push(pixel.V(x, y), pixel.V(x+b.BarWidth, y+h))
		b.imd.Rectangle(0)
	}
	b.imd.Draw(t)
}
