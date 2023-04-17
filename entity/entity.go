package entity

import (
	"context"

	"github.com/rafaeldajuda/database"
	v1 "github.com/rafaeldajuda/handler/v1"
)

type API struct {
	Name        string
	Host        string
	Port        string
	Ctx         context.Context
	MongoConfig database.MongoConfig
	v1.HandlerV1
}
