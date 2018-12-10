package service

import (
	"strings"

	proto "github.com/ankurs/Feed/Feed/Feed_proto"
)

type feedInfo struct {
	proto.FeedItem
}

func (f *feedInfo) GetVerb() string {
	return strings.ToLower(f.FeedItem.GetVerb().String())
}

func (f *feedInfo) GetCVerb() string {
	return strings.ToLower(f.FeedItem.GetCompatibilityVerb().String())
}
