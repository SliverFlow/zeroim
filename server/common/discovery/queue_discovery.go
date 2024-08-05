package discovery

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/discov"
)

// QueueObserver 队列观察者
type QueueObserver interface {
	Update(string, kq.KqConf)
	Delete(string)
}

// QueueDiscoveryProc 启动一个队列发现进程
func QueueDiscoveryProc(conf discov.EtcdConf, qo QueueObserver) {
	master, err := NewQueueMaster(conf.Key, conf.Hosts)
	if err != nil {
		panic(err)
	}
	master.Register(qo)
	master.WatchQueueWorkers()
}
