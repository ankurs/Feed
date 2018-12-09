package service

import (
	"context"
	"testing"

	proto "github.com/ankurs/Feed/Feed/Feed_proto"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	s := getTestService()
	req := new(proto.RegisterRequest)
	req.UserName = "test1"
	req.Password = "password"
	req.Email = "xyz@abc.com"
	req.FirstName = "ABC"
	req.LastName = "XYZ"
	resp, err := s.Register(context.Background(), req)
	assert.NotNil(t, resp, "register should return response")
	assert.NoError(t, err, "register should not return error")
}

func TestLogin(t *testing.T) {

}

func TestFeed(t *testing.T) {

}
