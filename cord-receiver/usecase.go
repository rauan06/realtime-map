package usecase

import (
	"github.com/redis/go-redis/v9"
)

type IRecieverUseCase interface {
	ProcessOBUData()
}

type RecieverUseCase struct {
	cache redis.Client
}
