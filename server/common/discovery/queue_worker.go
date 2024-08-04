package discovery

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// QueueWorker 一个队列的worker
type QueueWorker struct {
	cli    *clientv3.Client // etcd cli
	kqConf kq.KqConf        // 队列配置
	key    string           // 队列key
}

// NewQueueWorker 生成一个QueueWorker
func NewQueueWorker(key string, endpoints []string, kqConf kq.KqConf) (*QueueWorker, error) {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * 3,
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}

	return &QueueWorker{
		cli:    client,
		key:    key,
		kqConf: kqConf,
	}, nil
}

// HeartBeat 保持心跳
func (qw *QueueWorker) HeartBeat() {
	v, err := json.Marshal(qw.kqConf)
	if err != nil {
		panic(err)
	}

	qw.register(string(v))
}

// register 注册
func (qw *QueueWorker) register(value string) {
	// 申请租约
	leaseResp, err := qw.cli.Grant(context.Background(), 45)
	if err != nil {
		panic(err)
	}

	// 获取租约 id
	leaseId := leaseResp.ID
	logx.Infof("查看 leaseID :+ %v", leaseId)

	// 获取 kv api子集
	kv := clientv3.NewKV(qw.cli)

	// put 一个 key ，让他和租约进行关联，实现 10 秒自动过期
	putResp, err := kv.Put(context.TODO(), qw.key, value, clientv3.WithLease(leaseId))
	if err != nil {
		panic(err)
	}
	logx.Infof("查看 putResp :+ %v", putResp)

	// 自动续约
	leaseKeepAliveChan, err := qw.cli.KeepAlive(context.Background(), leaseId)
	if err != nil {
		panic(err)
	}

	// 处理续约应答
	go func() {
		for {
			select {
			case keepRes, ok := <-leaseKeepAliveChan:
				if !ok {
					logx.Infof("租约已经失效：%x", leaseId)
					qw.register(value)
					return
				} else {
					// 每秒续租一次，会有一次应答
					logx.Infof("收到自动续租应答：%v", keepRes.ID)
				}
			}
		}
	}()
}
