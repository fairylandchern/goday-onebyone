package gpool

type Job func()

// 协程池，简化版的处理逻辑，for循环监听管道信息，dispatch分发消息
type Worker struct {
	WorkerQueue chan *Worker
	JobChan     chan Job
	stopChan    chan struct{}
}

func newworker(queue chan *Worker) *Worker {
	return &Worker{
		JobChan:     make(chan Job),
		WorkerQueue: queue,
		stopChan:    make(chan struct{}),
	}
}

func (w *Worker) Start() {
	go func() {
		var job Job
		for {
			w.WorkerQueue <- w
			select {
			case job = <-w.JobChan:
				job()
				// 退出相关处理逻辑
			case <-w.stopChan:
				w.stopChan <- struct{}{}
				return
			}
		}
	}()
}

type Pool struct {
	jobQueue    chan Job
	workerQueue chan *Worker
	stopChan    chan struct{}
}

func NewPool(workerSize int, jobSize int) *Pool {
	p := &Pool{
		workerQueue: make(chan *Worker, workerSize),
		jobQueue:    make(chan Job, jobSize),
		stopChan:    make(chan struct{})}
	p.Start()
	return p
}

func (p *Pool) Start() {
	for i := 0; i < cap(p.workerQueue); i++ {
		worker := newworker(p.workerQueue)
		worker.Start()
	}
	go p.dispatch()
}

func (p *Pool) PushJob(job Job) {
	p.jobQueue <- job
}

func (p *Pool) dispatch() {
	for {
		select {
		case job := <-p.jobQueue:
			w := <-p.workerQueue
			w.JobChan <- job
			// 退出协程池相关的业务处理逻辑
		case <-p.stopChan:
			for i := 0; i < cap(p.workerQueue); i++ {
				w := <-p.workerQueue

				w.stopChan <- struct{}{}
				//log.Println("work stop after push stop")
				<-w.stopChan
				//log.Println("work stop after recv stop")
			}
			p.stopChan <- struct{}{}
			//log.Println("pool stop now send stop")
			return
		}
	}
}

// 有些地方存有疑问
func (p *Pool) Release() {
	//log.Println("pool stop here")
	p.stopChan <- struct{}{}
	//log.Println("pool stop after push stop")
	<-p.stopChan
	//log.Println("pool stop after recv stop")
}
