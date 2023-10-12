package logsmanager

type LogsManager struct {
	broadcasters map[string]*LogBroadcaster
}

func NewLogsManager() *LogsManager {
	return &LogsManager{
		broadcasters: map[string]*LogBroadcaster{},
	}
}

func (lm *LogsManager) NewLogBroadcaster(id string, pty string, options RotateWriterOptions) *LogBroadcaster {
	broadcaster := NewLogBroadcaster(pty, options)
	lm.broadcasters[id] = broadcaster

	return broadcaster
}

func (lm *LogsManager) GetLogBroadcaster(id string) *LogBroadcaster {
	broadcaster, ok := lm.broadcasters[id]
	if !ok {
		return nil
	}

	return broadcaster
}

func (lm *LogsManager) RemoveLogBroadcaster(id string) {
	lm.broadcasters[id].Stop()
	delete(lm.broadcasters, id)
}
