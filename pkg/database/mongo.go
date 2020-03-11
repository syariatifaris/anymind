package database

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const(
	dsnFmt = "mongodb://%s:%d"
)

func NewSimpleMongoClient(ctx context.Context, host string, port int) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(fmt.Sprintf(dsnFmt, host, port))
	clt, err := mongo.NewClient(opts)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create mongo client")
	}
	err = clt.Connect(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to mongodb")
	}
	return clt, nil
}