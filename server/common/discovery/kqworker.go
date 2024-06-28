package discovery

import clientv3 "go.etcd.io/etcd/client/v3"

type KqWorker struct {
	Key    string
	client *clientv3.Client
}
