package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"unsafe"
)

const ( // alias
	Ldate         = log.Ldate
	Ltime         = log.Ltime
	Lmicroseconds = log.Lmicroseconds
	Llongfile     = log.Llongfile
	Lshortfile    = log.Lshortfile
	LUTC          = log.LUTC
	LstdFlags     = Ldate | Ltime
)

var ( // alias
	New = log.New
)

var (
	mu sync.Mutex
	gw io.Writer
	f  *RotatedFile
)

var Output = log.Output

func SetOutput(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	if ww, ok := w.(*RotatedFile); ok {
		f = ww
	}
	gw = w
	log.SetOutput(w)
}

func SetFlags(flag int) {
	log.SetFlags(flag)
}

func Flags() int {
	return log.Flags()
}

func GetOutput() io.Writer {
	mu.Lock()
	defer mu.Unlock()
	if gw != nil {
		return gw
	}
	return os.Stderr
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func outputAndPutBuffer(buf *bytes.Buffer) {
	b := buf.Bytes()
	Output(3, *(*string)(unsafe.Pointer(&b)))
	buf.Reset()
	bufPool.Put(buf)
}

func Printf(format string, v ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.WriteString("INFO ")
	fmt.Fprintf(buf, format, v...)
	outputAndPutBuffer(buf)
}

func Println(v ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.WriteString("INFO ")
	fmt.Fprintln(buf, v...)
	outputAndPutBuffer(buf)
}

func Infof(format string, v ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.WriteString("INFO ")
	fmt.Fprintf(buf, format, v...)
	outputAndPutBuffer(buf)
}

func Infoln(v ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.WriteString("INFO ")
	fmt.Fprintln(buf, v...)
	outputAndPutBuffer(buf)
}

func Warnf(format string, v ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.WriteString("WARN ")
	fmt.Fprintf(buf, format, v...)
	outputAndPutBuffer(buf)
}

func Warnln(v ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.WriteString("WARN ")
	fmt.Fprintln(buf, v...)
	outputAndPutBuffer(buf)
}

func Errorf(format string, v ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.WriteString("ERROR ")
	fmt.Fprintf(buf, format, v...)
	outputAndPutBuffer(buf)
}

func Errorln(v ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.WriteString("ERROR ")
	fmt.Fprintln(buf, v...)
	outputAndPutBuffer(buf)
}

func Fatalf(format string, v ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.WriteString("FATAL ")
	fmt.Fprintf(buf, format, v...)
	outputAndPutBuffer(buf)
	if f != nil {
		f.Close() // wait flush
	}
	os.Exit(1)
}

func Fatalln(v ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.WriteString("FATAL ")
	fmt.Fprintln(buf, v...)
	outputAndPutBuffer(buf)
	if f != nil {
		f.Close() // wait flush
	}
	os.Exit(1)
}

func XInfo(aa ...interface{}) {
	x := xbufPool.Get().(*xbuffer)
	x.A("INFO").A(aa...).outputAndFree()
}

func XWarn(aa ...interface{}) {
	x := xbufPool.Get().(*xbuffer)
	x.A("WARN").A(aa...).outputAndFree()
}

func XError(aa ...interface{}) {
	x := xbufPool.Get().(*xbuffer)
	x.A("ERROR").A(aa...).outputAndFree()
}

func XInfof(format string, aa ...interface{}) {
	x := xbufPool.Get().(*xbuffer)
	x.A("INFO").Printf(format, aa...).outputAndFree()
}

func XWarnf(format string, aa ...interface{}) {
	x := xbufPool.Get().(*xbuffer)
	x.A("WARN").Printf(format, aa...).outputAndFree()
}

func XErrorf(format string, aa ...interface{}) {
	x := xbufPool.Get().(*xbuffer)
	x.A("ERROR").Printf(format, aa...).outputAndFree()
}

var ( // alias
	Print = Println
	Info  = Infoln
	Warn  = Warnln
	Error = Errorln
	Fatal = Fatalln
)

func (x *xbuffer) outputAndFree() {
	Output(3, *(*string)(unsafe.Pointer(&x.buf)))
	x.Reset()
	xbufPool.Put(x)
}
