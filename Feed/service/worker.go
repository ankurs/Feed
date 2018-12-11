package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/carousell/Orion/utils/worker"
)

const (
	workerAddFollowerFeed = "AddFollowerFeed"
)

type workFollowerTaks struct {
	FollowerId string `json:"follower_id"`
	FeedId     string `json:"feed_id"`
	Ts         int64  `json:"ts"`
}

func submitFollowingFeedItem(s *svc, ctx context.Context, followerId, feedId string, ts time.Time) {
	w := workFollowerTaks{
		FollowerId: followerId,
		FeedId:     feedId,
		Ts:         ts.Unix(),
	}
	data, _ := json.Marshal(w)
	s.worker.Schedule(ctx, workerAddFollowerFeed, string(data))
}

func processFollowingFeedItem(s *svc) worker.Work {
	return func(ctx context.Context, payload string) error {
		w := workFollowerTaks{}
		json.Unmarshal([]byte(payload), &w)
		return s.storage.AddFollowingFeedItem(ctx, w.FollowerId, w.FeedId, time.Unix(w.Ts, 0))
	}
}

func registerWorkerTask(s *svc) {
	s.worker.RegisterTask(workerAddFollowerFeed, processFollowingFeedItem(s))
}
