package storage

import (
	"context"
	"fmt"

	mem "github.com/rs/rest-layer-mem"
	mongo "github.com/rs/rest-layer-mongo"
	"github.com/rs/rest-layer/resource"
	"gopkg.in/mgo.v2"

	"github.com/ultimateboy/discapi/config"
)

// NewHandler returns a new resource.Storer handler for in-memory or mongodb
func NewHandler(ctx context.Context, session *mgo.Session, cfg *config.Config, entity string) resource.Storer {
	if cfg.InMemoryStorage {
		return mem.NewHandler()
	}

	return mongo.NewHandler(session, cfg.MongoDB, entity)
}

// NewSession returns nil for in-memory, or a MongoDB Session
func NewSession(ctx context.Context, cfg *config.Config) (*mgo.Session, error) {
	if cfg.InMemoryStorage {
		return nil, nil
	}

	return mgo.Dial(fmt.Sprintf("%s:%d", cfg.MongoHost, cfg.MongoPort))
}
