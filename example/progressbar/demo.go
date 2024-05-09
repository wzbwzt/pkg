package main

import (
	"io"
	"os"
	"time"

	"github.com/wzbwzt/pkg/pkg/progressbar"
)

func main() {
	progress_bar := progressbar.NewProgressBar(1000, []*progressbar.PhaseParam{
		{
			Phase:    0,
			MultiNum: 5,
			Name:     "Phase 1:",
			Units:    progressbar.UnitsDefault,
		},
		{
			Phase:    1,
			MultiNum: 5,
			Name:     "Phase 2:",
			Units:    progressbar.UnitsBytes,
		},
	})
	phase0_tracker := progress_bar.CreateTracker(0, "000", 1000000)
	for i := 0; i < 1000000; i++ {
		phase0_tracker.Increment(1)
	}
	progress_bar.MarkAsDone(phase0_tracker)

	phase1_tracker := progress_bar.CreateTracker(1, "111", 1000000)
	file, _ := os.CreateTemp("", "*_tmp")
	mutil_w := io.MultiWriter(phase1_tracker, file)
	read_file, _ := os.Open("./bigfile")
	io.Copy(mutil_w, read_file)

	progress_bar.MarkAsDone(phase1_tracker)
	time.Sleep(10 * time.Second)
}
