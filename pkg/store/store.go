package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsConf struct {
	User     string
	Password string
	Url      string
	Replicas int
	Stream   string
	Bucket   string
}

type Store struct {
	NatsConn *nats.Conn
	Replicas int
	Ctx      context.Context
	Bucket   string
}

func NewStore(ctx context.Context, nconf *NatsConf) (*Store, error) {
	conn, err := Connect(nconf)
	if err != nil {
		log.Errorf("error connecting to nats: %s", err)
		return nil, err
	}
	return &Store{
		NatsConn: conn,
		Replicas: nconf.Replicas,
		Ctx:      ctx,
		Bucket:   nconf.Bucket,
	}, nil

}

func (s *Store) Close() {
	s.NatsConn.Close()
}

func Connect(nconf *NatsConf) (*nats.Conn, error) {
	if nconf == nil {
		return nil, errors.New("NATS config is nil")
	}
	// nc, err := nats.Connect("connect.ngs.global", nats.UserCredentials("nats/creds/nats.creds"), nats.Name("Ivan Public Grid"))
	nc, err := nats.Connect(nconf.Url, nats.Name("Ivan Homelab Grid"), nats.UserInfo(nconf.User, nconf.Password))
	if err != nil || nc == nil {
		fmt.Println("Error connecting to NATS Cluster:", err, nconf.Url, nconf.User, nconf.Password)
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
		TTL:    time.Minute * 10,
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
