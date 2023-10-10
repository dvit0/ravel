package logger

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/valyentdev/ravel/internal/utils"
)

type RotateWriter struct {
	filename      string
	directory     string
	maxSizeByFile int
	maxFiles      int
	lock          sync.Mutex
	fp            *os.File
}

type RotateWriterOptions struct {
	Filename      string
	Directory     string
	MaxSizeByFile int // in bytes
	MaxFiles      int
}

// Make a new RotateWriter. Return nil if error occurs during setup.
func NewRotateWriter(options RotateWriterOptions) *RotateWriter {
	w := &RotateWriter{
		filename:      options.Filename,
		directory:     options.Directory,
		maxFiles:      options.MaxFiles,
		maxSizeByFile: options.MaxSizeByFile,
	}
	os.Mkdir(options.Directory, 0755)
	err := w.rotate()
	if err != nil {
		return nil
	}
	return w
}

func (w *RotateWriter) afterRotate() error {
	return w.removeOldestFile()
}

func (w *RotateWriter) shouldRotate() bool {
	stats, err := w.fp.Stat()
	if err != nil {
		return false
	}

	return stats.Size() > 1*utils.MB
}

func (w *RotateWriter) removeOldestFile() error {
	// List all files in directory
	files, err := os.ReadDir(w.directory)

	if len(files) > w.maxFiles {
		return nil
	}

	if err != nil {
		return err
	}

	// Find the oldest file
	var oldestFile os.DirEntry
	for _, file := range files {
		if oldestFile == nil {
			oldestFile = file
			continue
		}
		fileInfo, err := file.Info()
		if err != nil {
			return err
		}

		oldestFileInfo, err := oldestFile.Info()
		if err != nil {
			return err
		}

		if fileInfo.ModTime().Before(oldestFileInfo.ModTime()) {
			oldestFile = file
		}
	}

	// Remove the oldest file
	return os.Remove(oldestFile.Name())
}

// Write satisfies the io.Writer interface.
func (w *RotateWriter) Write(output []byte) (int, error) {
	log.Println("Writing to log", string(output))
	w.lock.Lock()
	defer w.lock.Unlock()
	defer func() {
		if w.shouldRotate() {
			w.rotate()
			w.afterRotate()
		}
	}()

	return w.fp.Write(output)
}

// Perform the actual act of rotating and reopening file.
func (w *RotateWriter) rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	filePath := filepath.Clean(w.directory + "/" + w.filename)

	// Close existing file if open
	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			return
		}
	}
	// Rename dest file if it already exists
	_, err = os.Stat(filePath)
	if err == nil {
		newFilePath := filePath + "." + time.Now().Format(time.RFC3339)
		err = os.Rename(filePath, newFilePath)
		if err != nil {
			return
		}
	}

	// Create a file.
	w.fp, err = os.Create(filePath)
	return
}
