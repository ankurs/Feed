package service

import proto "github.com/ankurs/Feed/Feed/Feed_proto"

type FeedService interface {
	proto.FeedServer
	proto.AccountServer
	proto.FollowServer
}
