package store

import (
	"context"
	"errors"
	"fmt"
	"i2/pkg/models"
	"time"

	"github.com/charmbracelet/log"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Store struct {
	NatsConn *nats.Conn
	Replicas int
	Ctx      context.Context
	Bucket   string
	Timeout  time.Duration
}

func NewStore(ctx context.Context, nconf *models.Nats) (*Store, error) {
	conn, err := Connect(nconf)
	if err != nil {
		log.Errorf("error connecting to nats: %s, %v", nconf.URL, err)
		return nil, err
	}
	timeout := time.Duration(nconf.Timeout) * time.Second
	if nconf.Timeout == 0 {
		timeout = 2 * time.Second
	}
	return &Store{
		NatsConn: conn,
		Replicas: nconf.Replicas,
		Ctx:      ctx,
		Bucket:   nconf.Bucket,
		Timeout:  timeout,
	}, nil

}

func (s *Store) Close() {
	if s.NatsConn != nil {
		s.NatsConn.Close()
	}
}

func Connect(nconf *models.Nats) (*nats.Conn, error) {
	if nconf == nil {
		return nil, errors.New("NATS config is nil")
	}
	log.Info("Connect")
	nc, err := nats.Connect(
		nconf.URL,
		nats.Name("Ivan Homelab Grid"),
		nats.UserInfo(nconf.User, nconf.Password),
		nats.Timeout(1*time.Second),
	)
	if err != nil || nc == nil {
		fmt.Println("Error connecting to NATS Cluster:", err, nconf.URL, nconf.User, nconf.Password)
		return nil, err
	}
	return nc, nil
}

func SetKV(ctx context.Context, key, bucket string, value []byte, nc *nats.Conn) error {
	if nc == nil {
		return fmt.Errorf("nats connection is nil")
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return err
	}
	metadataKVStore, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket: bucket,
		TTL:    time.Minute * 30,
	})

	if err != nil {
		return err
	}

	_, err = metadataKVStore.Put(ctx, key, value)

	return err
}

func GetKV(ctx context.Context, key, bucket string, nc *nats.Conn) ([]byte, error) {
	if nc == nil {
		return nil, fmt.Errorf("nats connection is nil")
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	metadataKVStore, err := js.KeyValue(ctx, bucket)
	if err != nil {
		return nil, err
	}

	metadataEntry, err := metadataKVStore.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return metadataEntry.Value(), nil
}

func GetKeys(ctx context.Context, bucket string, nc *nats.Conn) ([]string, error) {
	if nc == nil {
		return nil, fmt.Errorf("nats connection is nil")
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	metadataKVStore, err := js.KeyValue(ctx, bucket)
	if err != nil {
		return nil, err
	}
	keys, err := metadataKVStore.Keys(ctx)
	if err != nil {
		return nil, err
	}

	return keys, nil
}
