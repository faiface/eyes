package eyes

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type Marker struct {
	Width, Height  int
	Radius         float64
	ActivationTime time.Duration

	activeSpots  map[image.Point]bool
	spots        []spot
	numFixations int
	inFixation   bool

	imd *imdraw.IMDraw
}

type spot struct {
	activation   time.Duration
	timeSpent    time.Duration
	orderVisited int
}

func NewMarker(width, height int, radius float64, activationTime time.Duration) *Marker {
	return &Marker{
		Width:          width,
		Height:         height,
		Radius:         radius,
		ActivationTime: activationTime,
		activeSpots:    make(map[image.Point]bool),
		spots:          make([]spot, width*height),
		imd:            imdraw.New(nil),
	}
}

func (m *Marker) at(x, y int) *spot {
	return &m.spots[y*m.Width+x]
}

func (m *Marker) TimeSpentAt(x, y int) time.Duration {
	return m.at(x, y).timeSpent
}

func (m *Marker) OrderVisited(x, y int) int {
	return m.at(x, y).orderVisited
}

func (m *Marker) NumFixations() int {
	return m.numFixations
}

func (m *Marker) PointAt(target pixel.Vec, delta time.Duration) {
	var notActive []image.Point
	for pt := range m.activeSpots {
		s := m.at(pt.X, pt.Y)
		s.activation -= delta
		if s.activation < 0 {
			s.activation = 0
			notActive = append(notActive, pt)
		}
	}
	for _, pt := range notActive {
		delete(m.activeSpots, pt)
	}

	for y := target.Y - m.Radius; y <= target.Y+m.Radius; y++ {
		for x := target.X - m.Radius; x <= target.X+m.Radius; x++ {
			if math.Hypot(x-target.X, y-target.Y) > m.Radius {
				continue
			}
			ix, iy := int(x), int(y)
			if ix < 0 || ix >= m.Width || iy < 0 || iy >= m.Height {
				continue
			}
			s := m.at(ix, iy)
			s.activation += 2 * delta
			if s.activation > m.ActivationTime {
				s.activation = m.ActivationTime
			}
			m.activeSpots[image.Pt(ix, iy)] = true
		}
	}

	var increased []image.Point
	for pt := range m.activeSpots {
		s := m.at(pt.X, pt.Y)
		if s.activation >= m.ActivationTime {
			s.timeSpent += delta
			increased = append(increased, pt)
		}
	}

	if len(increased) == 0 && m.inFixation {
		m.inFixation = false
	} else if len(increased) > 0 && !m.inFixation {
		m.inFixation = true
		m.numFixations++
	}

	for _, pt := range increased {
		if m.at(pt.X, pt.Y).orderVisited == 0 {
			m.at(pt.X, pt.Y).orderVisited = m.numFixations
		}
	}
}

func (m *Marker) Draw(t pixel.Target, area pixel.Rect, clr color.Color) {
	var maxTimeSpent time.Duration
	for _, s := range m.spots {
		if s.timeSpent > maxTimeSpent {
			maxTimeSpent = s.timeSpent
		}
	}

	spotW, spotH := area.W()/float64(m.Width), area.H()/float64(m.Height)
	rgba := pixel.ToRGBA(clr)
	m.imd.Clear()

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if m.at(x, y).timeSpent == 0 {
				continue
			}
			brightness := float64(m.at(x, y).timeSpent) / float64(maxTimeSpent)
			m.imd.Color = rgba.Scaled(brightness)
			m.imd.Push(pixel.V(float64(x)*spotW, float64(y)*spotH))
			m.imd.Push(pixel.V(float64(x+1)*spotW, float64(y+1)*spotH))
			m.imd.Rectangle(0)
		}
	}

	m.imd.Draw(t)
}

func (m *Marker) TotalTimesOnFixations() []time.Duration {
	times := make([]time.Duration, m.NumFixations())
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if m.OrderVisited(x, y) == 0 {
				continue
			}
			times[m.OrderVisited(x, y)-1] += m.TimeSpentAt(x, y)
		}
	}
	return times
}
