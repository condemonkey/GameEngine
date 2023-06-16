package core

// 월드 레벨관리(static, instance)
type World struct {
	levels        map[string]Level
	addChan       chan Level
	removeChan    chan Level
	removeAllChan chan bool
}

func NewWorld() *World {
	return &World{
		levels:        make(map[string]Level),
		addChan:       make(chan Level),
		removeChan:    make(chan Level),
		removeAllChan: make(chan bool),
	}
}

func (w *World) RegisterLevel(level Level) {
	w.levels[level.Name()] = level
}

// 각 레벨들은 상호작용 X 고루틴으로 멀티처리
func (w *World) CreateLevel(name string) {
	go func() {
		level := w.levels[name].Create()
		level.Load()
		level.OnLoad()
		w.addChan <- level
	}()
}

// 각 레벨들은 상호작용 X 고루틴으로 멀티처리
// Level에서 호출함
func (w *World) RemoveLevel(level Level) {
	go func() {
		level.OnDestroy()
		w.removeChan <- level
	}()
}

func (w *World) Run() chan bool {
	closeChan := make(chan bool)
	go func() {
		for {
			select {
			case level := <-w.addChan:
				w.levels[level.Name()] = level
				go func() {
					<-level.Run()
					w.RemoveLevel(level)
				}()
			case level := <-w.removeChan:
				delete(w.levels, level.Name())
			case <-w.removeAllChan:
				for _, level := range w.levels {
					level.Stop()
				}
				closeChan <- true
			}
		}
	}()
	return closeChan
}

func (w *World) Stop() {
	w.removeAllChan <- true
}
