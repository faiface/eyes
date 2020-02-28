package main

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/eyes"
)

func computeTimesForSession(session eyes.Session) []time.Duration {
	marker := eyes.NewMarker(2048, 2048, eyes.DefaultMarkerRadius, eyes.DefaultActivationTime)

	for i := 1; i < len(session); i++ {
		dt := session[i].Time - session[i-1].Time
		middle := session[i].LeftPos.Add(session[i].RightPos).Scaled(0.5)
		marker.PointAt(middle, dt)
	}

	return marker.TotalTimesOnFixations()
}

func main() {
	numSessions := len(os.Args) - 1
	var totalTimes []time.Duration

	for i, path := range os.Args[1:] {
		fmt.Fprintf(os.Stderr, "Processing %s, file %d out of %d.\n", path, i+1, len(os.Args)-1)
		session, err := eyes.LoadSession(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		times := computeTimesForSession(session)

		for len(totalTimes) < len(times) {
			totalTimes = append(totalTimes, 0)
		}

		for i := range times {
			totalTimes[i] += times[i]
		}
	}

	averageTimes := make([]time.Duration, len(totalTimes))
	for i := range averageTimes {
		averageTimes[i] = totalTimes[i] / time.Duration(numSessions)
	}

	for i := range averageTimes {
		fmt.Println(averageTimes[i].Seconds())
	}
}
