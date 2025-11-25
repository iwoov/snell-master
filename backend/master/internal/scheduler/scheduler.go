package scheduler

import "time"

// Task represents a cancellable repeating job.
type Task struct {
	stop chan struct{}
}

// Stop halts the task.
func (t *Task) Stop() {
	if t == nil {
		return
	}
	select {
	case <-t.stop:
		return
	default:
		close(t.stop)
	}
}

func newTask(initialDelay, interval time.Duration, fn func()) *Task {
	task := &Task{stop: make(chan struct{})}
	go func() {
		timer := time.NewTimer(initialDelay)
		defer timer.Stop()
		for {
			select {
			case <-task.stop:
				return
			case <-timer.C:
				fn()
				if interval <= 0 {
					timer.Reset(initialDelay)
				} else {
					timer.Reset(interval)
				}
			}
		}
	}()
	return task
}

// Manager 统一管理多个任务以便关闭。
type Manager struct {
	tasks []*Task
}

// NewManager 创建任务管理器。
func NewManager() *Manager {
	return &Manager{tasks: make([]*Task, 0)}
}

// Add 注册任务。
func (m *Manager) Add(task *Task) {
	if task == nil {
		return
	}
	m.tasks = append(m.tasks, task)
}

// Stop 停止所有任务。
func (m *Manager) Stop() {
	for _, task := range m.tasks {
		task.Stop()
	}
}
