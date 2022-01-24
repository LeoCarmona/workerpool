package workerpool

import (
	"fmt"
	"time"
)

type Settings struct {
	MinWorkers   int
	MaxWorkers   int
	IdleTimeout  time.Duration
	UpScaling    int
	DownScaling  int
	Queue        int
	PanicHandler func(panicErr interface{})
}

func (s *Settings) String() string {
	return fmt.Sprintf(settingsStringFormat, s.MinWorkers, s.MaxWorkers, s.IdleTimeout.String(), s.UpScaling, s.DownScaling, s.Queue)
}

const settingsStringFormat = "Settings [minWorkers=%d, maxWorkers=%d, idleTimeout=%s, upScaling=%d, downScaling=%d, queue=%v]"
