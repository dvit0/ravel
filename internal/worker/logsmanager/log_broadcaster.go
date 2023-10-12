package logsmanager

import (
	"context"
	"io"
	"os"

	"github.com/charmbracelet/log"
	"github.com/containerd/console"
)

type Subscriber struct {
	writer io.Writer
	ctx    context.Context
}

type LogBroadcaster struct {
	pty         string
	writer      *RotateWriter
	closed      bool
	subscribers []*Subscriber
}

func NewLogBroadcaster(pty string, options RotateWriterOptions) *LogBroadcaster {
	rotateWriter, err := NewRotateWriter(options)
	if err != nil {
		log.Fatal("Failed to create rotate writer", "err", err)

	}
	return &LogBroadcaster{
		pty:         pty,
		writer:      rotateWriter,
		closed:      false,
		subscribers: []*Subscriber{},
	}
}

func (lb *LogBroadcaster) Start() {
	file, err := os.Open(lb.pty)
	if err != nil {
		log.Error("Failed to open pty", "err", err)
	}
	pty, err := console.ConsoleFromFile(file)
	if err != nil {
		log.Error("Failed to open console on pty", "err", err)
	}

	log.Info("Starting log broadcaster", "pty", pty)

	go func(lb *LogBroadcaster) {
		for !lb.closed {
			if lb == nil || lb.writer == nil {
				log.Error("Log broadcaster is nil")
				return
			}
			buf := make([]byte, 1024)
			n, err := pty.Read(buf)
			if err != nil {
				log.Error("Failed to read from pty", "err", err)
			}
			lb.writer.Write(buf[:n])

			for i, subscriber := range lb.subscribers {
				// if the context is done, remove the subscriber
				if subscriber.ctx.Err() != nil {
					lb.subscribers = append(lb.subscribers[:i], lb.subscribers[i+1:]...)
					continue
				}
				_, err := subscriber.writer.Write(buf[:n])
				if err != nil {
					log.Error("Failed to write to subscriber", "err", err)
				}

			}
		}
	}(lb)
}

func (lb *LogBroadcaster) Subscribe(ctx context.Context, writer io.Writer) {
	subscriber := &Subscriber{
		writer: writer,
		ctx:    ctx,
	}

	lb.subscribers = append(lb.subscribers, subscriber)

	<-ctx.Done()
	log.Info("Removing subscriber")
}

func (lb *LogBroadcaster) Stop() {
	lb.closed = true
	lb.writer.Close()
}
