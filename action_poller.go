package worker

import (
	"strconv"
	"sync"

	"github.com/mohitkumar/orchy-worker/client"
)

type nWorker struct {
	worker *worker
	num    int
}
type actionPoller struct {
	workers      []*nWorker
	pollerWorker []*pollerWorker
	clusterConf  *client.ClusterConf
	config       WorkerConfiguration
}

func newActionPoller(conf WorkerConfiguration, clusterConf *client.ClusterConf) *actionPoller {
	return &actionPoller{
		config:      conf,
		clusterConf: clusterConf,
	}
}

func (tp *actionPoller) registerWorker(worker *worker, numWorkers int) {
	tp.workers = append(tp.workers, &nWorker{worker: worker, num: numWorkers})
}

func (tp *actionPoller) start(wg *sync.WaitGroup) {
	for _, w := range tp.workers {
		for i := 0; i < w.num; i++ {
			client, err := client.NewRpcClient(tp.config.ServerUrl, tp.clusterConf)
			if err != nil {
				panic(err)
			}
			pw := &pollerWorker{
				worker:     w.worker,
				stop:       make(chan struct{}),
				client:     client,
				wg:         wg,
				workerName: w.worker.GetName() + "_" + strconv.Itoa(i),
			}
			tp.pollerWorker = append(tp.pollerWorker, pw)
			pw.Start()
		}
	}
}

func (tp *actionPoller) stop() {
	for _, w := range tp.pollerWorker {
		w.Stop()
	}
}
