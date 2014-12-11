package parallel

type Task func() error
type TaskCallback func(interface{}, error)

type Parallel struct {
	taskChan  chan taskCompletion
	tasks     []taskWrapper
	remaining int
}

type taskWrapper struct {
	task Task
	tag  interface{}
}

type taskCompletion struct {
	tag interface{}
	err error
}

func New() *Parallel {
	return &Parallel{
		make(chan taskCompletion),
		[]taskWrapper{},
		0,
	}
}

func (p *Parallel) Submit(task Task, tag interface{}) {
	p.tasks = append(p.tasks, taskWrapper{task, tag})
}

func (p *Parallel) Exec(callback TaskCallback) {
	p.remaining = len(p.tasks)
	if p.remaining == 0 {
		return
	}
	for _, w := range p.tasks {
		wrapper := w
		go func() {
			if err := wrapper.task(); err != nil {
				p.taskChan <- taskCompletion{wrapper.tag, err}
				return
			}
			p.taskChan <- taskCompletion{wrapper.tag, nil}

		}()
	}

	for {
		c := <-p.taskChan
		callback(c.tag, c.err)
		p.remaining--
		if p.remaining == 0 {
			close(p.taskChan)
			return
		}
	}
}
