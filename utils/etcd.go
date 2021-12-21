package utils

import (
	"context"
	"github.com/ropon/work_api/conf"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func EtCdSet(key, val string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	_, err = conf.Etcd.Put(ctx, key, val)
	cancel()
	return
}

func EtCdGet(key string) (resp *clientv3.GetResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err = conf.Etcd.Get(ctx, key)
	cancel()
	return
}
