package eyes

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/faiface/pixel"
)

type Session []SessionDataPoint

type SessionDataPoint struct {
	Time     time.Duration
	LeftPos  pixel.Vec
	RightPos pixel.Vec
}

func DecodeSession(r io.Reader) (Session, error) {
	cr := csv.NewReader(r)
	_, err := cr.Read() // first line is field names
	if err != nil {
		return nil, err
	}
	var session Session
	lineNum := 0
	for {
		line, err := cr.Read()
		lineNum++
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", lineNum, err)
		}
		var (
			timeData   int64
			leftXData  float64
			leftYData  float64
			rightXData float64
			rightYData float64
		)
		_, _ = fmt.Sscanf(line[0], "%d", &timeData)
		_, _ = fmt.Sscanf(line[1], "%f", &leftXData)
		_, _ = fmt.Sscanf(line[2], "%f", &leftYData)
		_, _ = fmt.Sscanf(line[3], "%f", &rightXData)
		_, _ = fmt.Sscanf(line[4], "%f", &rightYData)
		session = append(session, SessionDataPoint{
			Time:     time.Second / 1000 * time.Duration(timeData),
			LeftPos:  pixel.V(leftXData, leftYData),
			RightPos: pixel.V(rightXData, rightYData),
		})
	}
	// normalize times
	if len(session) > 0 {
		firstTime := session[0].Time
		for i := range session {
			session[i].Time -= firstTime
		}
	}
	return session, nil
}
