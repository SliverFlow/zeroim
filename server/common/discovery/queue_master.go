package discovery

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// QueueMaster 队列master
type QueueMaster struct {
	members  map[string]kq.KqConf // 队列配置
	cli      *clientv3.Client     // etcd cli
	rootPath string               // 服务发现根路径
	observer QueueObserver        // 观察者
}

// NewQueueMaster 生成一个QueueMaster
func NewQueueMaster(rootPath string, hosts []string) (*QueueMaster, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: hosts,
	})
	if err != nil {
		return nil, err
	}

	return &QueueMaster{
		members:  make(map[string]kq.KqConf),
		cli:      cli,
		rootPath: rootPath,
	}, err
}

// Register 注册
func (m *QueueMaster) Register(o QueueObserver) {
	m.observer = o
}

// notifyUpdate 通知更新
func (m *QueueMaster) notifyUpdate(key string, kqConf kq.KqConf) {
	m.observer.Update(key, kqConf)
}

// notifyDelete 通知删除
func (m *QueueMaster) notifyDelete(key string) {
	m.observer.Delete(key)
}

// addQueueWorker 添加队列配置
func (m *QueueMaster) addQueueWorker(key string, kqConf kq.KqConf) {
	if len(kqConf.Brokers) == 0 || len(kqConf.Topic) == 0 {
		logx.Errorf("invalid kqConf: %+v", kqConf)
		return
	}

	m.members[key] = kqConf
	m.notifyUpdate(key, kqConf)
}

// updateQueueWorker 更新队列配置
func (m *QueueMaster) updateQueueWorker(key string, kqConf kq.KqConf) {
	if len(kqConf.Brokers) == 0 || len(kqConf.Topic) == 0 {
		logx.Errorf("invalid kqConf: %+v", kqConf)
		return
	}

	m.members[key] = kqConf
	m.notifyUpdate(key, kqConf)
}

// deleteQueueWorker 删除队列配置
func (m *QueueMaster) deleteQueueWorker(key string) {
	delete(m.members, key)
	m.notifyDelete(key)
}

// WatchQueueWorkers 监听队列配置
func (m *QueueMaster) WatchQueueWorkers() {
	rch := m.cli.Watch(context.Background(), m.rootPath, clientv3.WithPrefix())

	for wresp := range rch {
		if wresp.Err() != nil {
			logx.Severe(wresp.Err())
		}

		if wresp.Canceled {
			logx.Severe("watch is canceled")
		}

		for _, event := range wresp.Events {
			switch event.Type {
			case clientv3.EventTypePut:
				var kqConf kq.KqConf
				if err := json.Unmarshal(event.Kv.Value, &kqConf); err != nil {
					logx.Error(err)
					continue
				}

				if event.IsCreate() {
					m.addQueueWorker(string(event.Kv.Key), kqConf)
				} else if event.IsModify() {
					m.updateQueueWorker(string(event.Kv.Key), kqConf)
				}
			case clientv3.EventTypeDelete:
				m.deleteQueueWorker(string(event.Kv.Key))
			}
		}
	}
}
