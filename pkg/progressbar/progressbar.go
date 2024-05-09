package progressbar

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Phase int
type Units progress.Units

var (
	UnitsBytes   = Units(progress.UnitsBytes)
	UnitsDefault = Units(progress.UnitsDefault)
)

type PhaseParam struct {
	Phase    Phase
	MultiNum int
	Name     string
	Units    Units

	done chan *Tracker
}
type ProgressBar struct {
	progress.Writer
	total     int
	handler   int64
	phases    []*PhaseParam
	st        time.Time
	lastPhase Phase
}
type Tracker struct {
	*progress.Tracker
	phase Phase
}

func (t *Tracker) Write(p []byte) (n int, err error) {
	n = len(p)
	t.Increment(int64(n))
	return
}

func (p *ProgressBar) MarkAsDone(t *Tracker) {
	phase := p.phases[t.phase]
	t.MarkAsDone()
	t.Reset()
	t.UpdateMessage(fmt.Sprintf("%s Done", phase.Name))
	t.UpdateTotal(0)
	if p.lastPhase == t.phase {
		atomic.AddInt64(&p.handler, 1)

		p.SetPinnedMessages(
			fmt.Sprintf(">> Total  :%d", p.total),
			fmt.Sprintf(">> Handler:%d", p.handler),
			fmt.Sprintf(">> Percent:%4.1f%%", float64(p.handler)/float64(p.total)*100),
			fmt.Sprintf(">> Duration:%s", time.Since(p.st).String()),
		)
	}
	phase.done <- t
}

func NewProgressBar(total int, phases []*PhaseParam) *ProgressBar {
	out := ProgressBar{
		total:     total,
		st:        time.Now(),
		lastPhase: Phase(len(phases) - 1),
	}

	totalMultiNum := 0
	trackers := [][]*progress.Tracker{}
	for _, v := range phases {
		totalMultiNum += v.MultiNum
		v.done = make(chan *Tracker, v.MultiNum)

		phase_trackers := []*progress.Tracker{}
		for i := 0; i < v.MultiNum; i++ {
			t := &progress.Tracker{Message: fmt.Sprintf("%s Readying", v.Name), Total: 0, Units: progress.Units(v.Units)}
			phase_trackers = append(phase_trackers, t)
			v.done <- &Tracker{Tracker: t, phase: v.Phase}
		}
		trackers = append(trackers, phase_trackers)
	}

	out.phases = phases

	pw := progress.NewWriter()
	pw.SetAutoStop(false)
	pw.SetMessageLength(36)
	pw.SetNumTrackersExpected(totalMultiNum)
	pw.SetSortBy(progress.SortByNone)
	pw.SetStyle(progress.StyleDefault)
	pw.SetTrackerLength(25)
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetUpdateFrequency(time.Millisecond * 100)
	pw.Style().Colors = progress.StyleColorsExample
	pw.Style().Options.PercentFormat = "%4.1f%%"
	pw.Style().Visibility.ETA = true
	pw.Style().Visibility.ETAOverall = true
	pw.Style().Visibility.Percentage = true
	pw.Style().Visibility.Speed = true
	pw.Style().Visibility.SpeedOverall = false
	pw.Style().Visibility.Time = true
	pw.Style().Visibility.TrackerOverall = false
	pw.Style().Visibility.Value = false
	pw.Style().Visibility.Pinned = true

	go pw.Render()

	out.Writer = pw

	for _, t := range trackers {
		pw.AppendTrackers(t)
	}

	return &out
}

func (pb *ProgressBar) CreateTracker(phase Phase, msg string, total int64) *Tracker {
	ph := pb.phases[phase]
	tracker := <-ph.done
	tracker.UpdateMessage(text.FgWhite.Sprintf("%s %s", ph.Name, msg))
	tracker.UpdateTotal(total)

	return tracker
}
