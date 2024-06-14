package cpuload

import (
	"context"
	"math"
	"time"

	"github.com/shirou/gopsutil/cpu"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

type CPULoadScheduler struct {
	handle framework.Handle
}

var _ = framework.ScorePlugin(&CPULoadScheduler{})

const Name = "CPULoadScheduler"

func (ts *CPULoadScheduler) Name() string {
	return Name
}

func (t *CPULoadScheduler) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	currentTime := time.Now().Hour()
	if currentTime >= 20 && currentTime <= 21 {
		cpuUsage, err := cpu.Percent(0, true)
		if err != nil {
			return 0, nil
		}
		if cpuUsage[0] > 20 {
			return 100, nil
		} else {
			return -100, nil
		}
	} else {
		return -100, nil
	}
}

func (ts *CPULoadScheduler) ScoreExtensions() framework.ScoreExtensions {
	return ts
}

func (ts *CPULoadScheduler) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	var highest int64 = -math.MaxInt64
	var lowest int64 = math.MaxInt64
	for _, nodeScore := range scores {
		if nodeScore.Score > highest {
			highest = nodeScore.Score
		}
		if nodeScore.Score < lowest {
			lowest = nodeScore.Score
		}
	}

	// Преобразуем диапазон баллов в диапазон, соответствующий требованиям фреймворка
	oldRange := highest - lowest
	newRange := framework.MaxNodeScore - framework.MinNodeScore
	for i, nodeScore := range scores {
		if oldRange == 0 {
			scores[i].Score = framework.MinNodeScore
		} else {
			scores[i].Score = ((nodeScore.Score - lowest) * newRange / oldRange) + framework.MinNodeScore
		}
	}

	return nil
}

func New(_ context.Context, _ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	return &CPULoadScheduler{handle: h}, nil
}
