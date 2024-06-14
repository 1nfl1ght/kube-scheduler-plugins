package timebased

import (
	"context"
	"math"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

type TimeScheduler struct {
	handle framework.Handle
}

var _ = framework.ScorePlugin(&TimeScheduler{})

const Name = "TimeScheduler"

func (ts *TimeScheduler) Name() string {
	return Name
}

func (ts *TimeScheduler) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	currentTime := time.Now()
	if currentTime.Hour() == 19 && currentTime.Minute() == 0 {
		return 100, nil
	} else if currentTime.Hour() == 20 && currentTime.Minute() == 0 {
		return -100, nil
	}

	return 0, nil
}

func (ts *TimeScheduler) ScoreExtensions() framework.ScoreExtensions {
	return ts
}

func (ts *TimeScheduler) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
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
	return &TimeScheduler{handle: h}, nil
}
