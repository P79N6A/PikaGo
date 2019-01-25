package log

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrLogBufferFull = errors.New("gopkg/pkg/log: buffer full")
	ErrLogClosed     = errors.New("gopkg/pkg/log: closed")
)

// SetupLogOutput configures logger of log.Output
func SetupLog2RotatedFile(filename string, maxFileSize int64, maxRotateNum int) (*RotatedFile, error) {
	c := DefaultRotatedConfig()
	c.MaxRotatedFileSize = maxFileSize
	c.MaxRotatedFileNum = maxRotateNum
	f, err := NewRotatedFile(filename, c)
	if err != nil {
		return nil, err
	}
	SetOutput(f)
	return f, nil
}

// RotatedFile represents the rotated file for log
type RotatedFile struct {
	filename string

	c *RotatedConfig

	mu      sync.Mutex
	wbuf    []byte
	wnotify chan struct{}
	closed  bool

	f *os.File

	curSize int64
	curIdx  int

	wg sync.WaitGroup
}

type RotatedConfig struct {
	FileMode           os.FileMode
	MaxRotatedFileSize int64
	MaxRotatedFileNum  int
	AsyncBufferSize    int
	AsyncWrite         bool
}

func DefaultRotatedConfig() *RotatedConfig {
	return &RotatedConfig{
		FileMode:           os.FileMode(0644),
		MaxRotatedFileSize: 200 << 20, // 200MiB
		MaxRotatedFileNum:  10,
		AsyncBufferSize:    1 << 20, // 1MiB
		AsyncWrite:         true,
	}
}

// NewRotatedFile creates instance of RotatedFile
func NewRotatedFile(filename string, c *RotatedConfig) (*RotatedFile, error) {
	if filename == "" {
		return nil, errors.New("filename empty")
	}
	if c == nil {
		c = DefaultRotatedConfig()
	}
	if c.FileMode == 0 {
		c.FileMode = os.FileMode(0644)
	}

	f := RotatedFile{filename: filename, c: c}

	basename := filepath.Base(filename)
	dir := filepath.Dir(filename)
	os.MkdirAll(dir, 0755)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var t time.Time
	for _, fi := range files {
		if fi.IsDir() || fi.Name() == basename {
			continue
		}
		if !strings.HasPrefix(fi.Name(), basename) {
			continue
		}
		if fi.ModTime().Before(t) {
			continue
		}
		t = fi.ModTime()
		s := strings.TrimPrefix(fi.Name(), basename+".")
		if n, err := strconv.Atoi(s); err == nil {
			f.curIdx = n
			f.curSize = fi.Size()
		}
	}
	flags := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	if f.curSize > c.MaxRotatedFileSize {
		flags |= os.O_TRUNC
		f.curSize = 0
		f.curIdx = (f.curIdx + 1) % f.c.MaxRotatedFileNum
	}
	_f, err := f.openfile(flags)
	if err != nil {
		return nil, err
	}
	f.f = _f
	f.MkSymlink()

	if c.AsyncWrite {
		if c.AsyncBufferSize <= 0 {
			c.AsyncBufferSize = 1 << 20
		}
		f.wbuf = make([]byte, 0, c.AsyncBufferSize) // 1MB
		f.wnotify = make(chan struct{}, 1)
		f.wg.Add(1)
		go f.asyncloop(f.wnotify)
	}
	return &f, nil
}

// Write implements io.Writer
func (f *RotatedFile) Write(b []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.c.AsyncWrite {
		return f.asyncwrite(b)
	}
	return f.writeAndRotate(b)
}

func (f *RotatedFile) Close() error {
	var wnotify chan struct{}
	f.mu.Lock()
	wnotify, f.wnotify = f.wnotify, nil
	f.mu.Unlock()
	if wnotify == nil {
		return ErrLogClosed
	}
	for f.Buffered() > 0 {
		time.Sleep(time.Millisecond)
	}
	close(wnotify)
	f.wg.Wait()
	return nil
}

func (f *RotatedFile) Buffered() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.wbuf)
}

// CurrentFilename returns current filename that should be appended to
func (f *RotatedFile) CurrentFilename() string {
	rn := f.c.MaxRotatedFileNum
	idx := f.curIdx
	suffix := ""
	if rn <= 100 {
		suffix = fmt.Sprintf(".%02d", idx)
	} else if rn <= 1000 {
		suffix = fmt.Sprintf(".%03d", idx)
	} else if rn <= 10000 {
		suffix = fmt.Sprintf(".%04d", idx)
	} else {
		suffix = fmt.Sprintf(".%05d", idx)
	}
	return f.filename + suffix
}

// MkSymlink creates sym link by CurrentFilename
func (f *RotatedFile) MkSymlink() {
	os.Remove(f.filename)
	os.Symlink(filepath.Base(f.CurrentFilename()), f.filename)
}

func (f *RotatedFile) asyncwrite(b []byte) (n int, err error) {
	if f.wnotify == nil {
		return 0, ErrLogClosed
	}
	if len(f.wbuf)+len(b) > cap(f.wbuf) {
		return 0, ErrLogBufferFull
	}
	f.wbuf = append(f.wbuf, b...)
	select {
	case f.wnotify <- struct{}{}:
	default:
	}
	return len(b), nil
}

// Run starts async write loop
func (f *RotatedFile) asyncloop(wnotify <-chan struct{}) {
	defer f.wg.Done()

	buf := make([]byte, 0, f.c.AsyncBufferSize)
	for range wnotify {
		f.mu.Lock() // swap buffer
		buf = buf[:0]
		buf, f.wbuf = f.wbuf, buf
		f.mu.Unlock()
		f.writeAndRotate(buf)
	}
	f.f.Close()
}

func (f *RotatedFile) writeAndRotate(b []byte) (n int, err error) {
	n, err = f.f.Write(b)
	f.curSize += int64(n)
	if f.curSize > f.c.MaxRotatedFileSize {
		f.curSize = 0
		f.curIdx = (f.curIdx + 1) % f.c.MaxRotatedFileNum
		nextf, er := f.openfile(os.O_WRONLY | os.O_APPEND | os.O_CREATE | os.O_TRUNC)
		if er != nil {
			f.f.Truncate(0) // XXX: avoid ulimit grow
			return n, er
		}
		f.f.Close()
		f.f = nextf
		f.MkSymlink()
	}
	return
}

func (f *RotatedFile) openfile(flag int) (*os.File, error) {
	ret, err := os.OpenFile(f.CurrentFilename(), flag, f.c.FileMode)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
