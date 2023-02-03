package storage

import (
	"context"
	"time"

	pkgerrors "github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"

	loggerutils "github.com/jkrus/master_api/pkg/etcd-configer/internal/logger-utils"
	"github.com/jkrus/master_api/pkg/etcd-configer/loading"
)

// OptionStorage ...
type OptionStorage struct {
	etcdClient *clientv3.Client
}

// NewOptionStorage ...
func NewOptionStorage(endpoint string, dialTimeout time.Duration) (OptionStorage, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: dialTimeout,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
		// add a minimal logger configuration for suppressing logs from the etcd client
		LogConfig: loggerutils.NewMinimalLoggerConfig(),
	})
	if err != nil {
		return OptionStorage{}, pkgerrors.WithMessage(err, "unable to create the etcd client")
	}

	return OptionStorage{etcdClient}, nil
}

// GetOptions ...
func (storage OptionStorage) GetOptions(
	ctx context.Context,
	namePrefix string,
) ([]loading.Option, error) {
	response, err := storage.etcdClient.Get(ctx, namePrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, pkgerrors.WithMessage(err, "unable to get options")
	}

	var options []loading.Option
	for _, kv := range response.Kvs {
		option := loading.Option{Name: string(kv.Key), Value: kv.Value}
		options = append(options, option)
	}

	return options, nil
}

// WatchOptions ...
func (storage OptionStorage) WatchOptions(
	ctx context.Context,
	namePrefix string,
) chan loading.OptionUpdate {
	optionChannel := make(chan loading.OptionUpdate)
	watchingChannel := storage.etcdClient.Watch(ctx, namePrefix, clientv3.WithPrefix())
	go func() {
		defer close(optionChannel)

		for watchingResponse := range watchingChannel {
			if err := watchingResponse.Err(); err != nil {
				optionChannel <- loading.OptionUpdate{Error: err}
				continue
			}

			for _, event := range watchingResponse.Events {
				if event.Type == clientv3.EventTypePut {
					option := loading.Option{Name: string(event.Kv.Key), Value: event.Kv.Value}
					optionChannel <- loading.OptionUpdate{Option: option}
				}
			}
		}
	}()

	return optionChannel
}
