package service

import "github.com/ankurs/Feed/Feed/service/store"

type Config struct {
	Store  store.Config
	Worker WorkerConfig
}

type WorkerConfig struct {
	Host     string
	Queue    string
	Username string
	Password string
}
