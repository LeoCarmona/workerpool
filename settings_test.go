package workerpool

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSettings_String(t *testing.T) {
	settings := &Settings{}
	assert.Equal(t, "Settings [minWorkers=0, maxWorkers=0, idleTimeout=0s, upScaling=0, downScaling=0, queue=0]", settings.String())

	settings = &Settings{
		MinWorkers:  1,
		MaxWorkers:  2,
		IdleTimeout: 3 * time.Second,
		UpScaling:   4,
		DownScaling: 5,
		Queue:       6,
	}
	assert.Equal(t, "Settings [minWorkers=1, maxWorkers=2, idleTimeout=3s, upScaling=4, downScaling=5, queue=6]", settings.String())
}
