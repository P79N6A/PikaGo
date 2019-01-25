package log

import (
	"strconv"
	"sync"
	"unsafe"
)

type xbuffer struct {
	buf []byte
}

var xbufPool = sync.Pool{
	New: func() interface{} {
		ret := new(xbuffer)
		ret.buf = make([]byte, 0, 500)
		return ret
	},
}

func (x *xbuffer) Reset() *xbuffer {
	x.buf = x.buf[:0]
	return x
}

func (x *xbuffer) Bytes() []byte {
	return x.buf
}

func (x *xbuffer) space() {
	if len(x.buf) != 0 {
		x.buf = append(x.buf, ' ')
	}
}

func appendFloat(dst []byte, f float64) []byte {
	return strconv.AppendFloat(dst, f, 'f', 6, 64)
}

func appendInt(dst []byte, i int64) []byte {
	return strconv.AppendInt(dst, i, 10)
}

func appendPt(dst []byte, p uintptr) []byte {
	if p == 0 {
		return append(dst, '<', 'n', 'i', 'l', '>')
	}
	dst = append(dst, '0', 'x')
	return strconv.AppendUint(dst, uint64(p), 16)
}

func (x *xbuffer) A(aa ...interface{}) *xbuffer {
	for i := range aa {
		x.space()
		switch v := aa[i].(type) {
		case float32:
			x.buf = appendFloat(x.buf, float64(v))
		case float64:
			x.buf = appendFloat(x.buf, v)
		case int:
			x.buf = appendInt(x.buf, int64(v))
		case int8:
			x.buf = appendInt(x.buf, int64(v))
		case int16:
			x.buf = appendInt(x.buf, int64(v))
		case int32:
			x.buf = appendInt(x.buf, int64(v))
		case int64:
			x.buf = appendInt(x.buf, int64(v))
		case uint:
			x.buf = appendInt(x.buf, int64(v))
		case uint8:
			x.buf = appendInt(x.buf, int64(v))
		case uint16:
			x.buf = appendInt(x.buf, int64(v))
		case uint32:
			x.buf = appendInt(x.buf, int64(v))
		case uint64:
			x.buf = appendInt(x.buf, int64(v))
		case string:
			// XXX: https://github.com/golang/go/issues/15730
			p := uintptr(unsafe.Pointer(&v))
			s := *(*string)(unsafe.Pointer(p)) // go vet: I am sure what I am doing
			x.buf = append(x.buf, s...)
		case *float32:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *float64:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *int:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *int8:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *int16:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *int32:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *int64:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *uint:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *uint8:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *uint16:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *uint32:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *uint64:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *string:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		}
	}
	return x
}

// Printf supports a simpfy format:
// %v	represents a var
//
// for float. format: %[.prec]specifier
// %.3f is the same as strconv.FormatFloat(v, "f", 3, 64)
//
// for int. format: %specifier
// %d	base 10 (signed)
// %x	base 16 (unsigned)
// %o	base 8  (unsigned)
//
func (x *xbuffer) Printf(format string, aa ...interface{}) *xbuffer {
	// %[flags][width][.precision]specifier
	end := len(format)
	aIdx := 0

	for i := 0; i < end; {
		lasti := i
		for i < end && format[i] != '%' {
			i++
		}
		if i > lasti {
			x.buf = append(x.buf, format[lasti:i]...)
		}
		if i >= end {
			break // done
		}
		i++ // skip %

		// ignore width
		for i < end && format[i] >= '0' && format[i] <= '9' {
			i++
		}

		// precision
		precision := -1
		if i < end && format[i] == '.' {
			precision = 0
			i++ // skip .
			for i < end && format[i] >= '0' && format[i] <= '9' {
				precision = precision*10 + int(format[i]-'0')
				i++
			}
			if precision > 1e6 || precision < -1e6 {
				precision = -1
			}
		}
		// specifier
		spec := byte('v')
		if i < end {
			spec = format[i]
			i++ // skip spec
		}
		if spec == '%' { // %%
			x.buf = append(x.buf, '%')
			continue
		}
		if aIdx >= len(aa) {
			x.buf = append(x.buf, '%', '!', spec, '(', 'M', 'I', 'S', 'S', 'I', 'N', 'G', ')')
			continue
		}

		a := aa[aIdx]
		aIdx++

		var tFloat, tInt bool
		var vf float64
		var vi int64
		switch v := a.(type) {
		case float32:
			vf, tFloat = float64(v), true
		case float64:
			vf, tFloat = v, true
		case int:
			vi, tInt = int64(v), true
		case int8:
			vi, tInt = int64(v), true
		case int16:
			vi, tInt = int64(v), true
		case int32:
			vi, tInt = int64(v), true
		case int64:
			vi, tInt = int64(v), true
		case uint:
			vi, tInt = int64(v), true
		case uint8:
			vi, tInt = int64(v), true
		case uint16:
			vi, tInt = int64(v), true
		case uint32:
			vi, tInt = int64(v), true
		case uint64:
			vi, tInt = int64(v), true
		case string:
			// XXX: https://github.com/golang/go/issues/15730
			p := uintptr(unsafe.Pointer(&v))
			s := *(*string)(unsafe.Pointer(p))
			x.buf = append(x.buf, s...)
		case *float32:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *float64:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *int:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *int8:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *int16:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *int32:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *int64:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *uint:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *uint8:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *uint16:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *uint32:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *uint64:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		case *string:
			x.buf = appendPt(x.buf, uintptr(unsafe.Pointer(v)))
		default:
			if v == nil {
				x.buf = append(x.buf, '<', 'n', 'i', 'l', '>')
			} else {
				x.buf = append(x.buf, '<', '?', '>')
			}
		}
		if tInt {
			switch spec {
			case 'x':
				x.buf = strconv.AppendUint(x.buf, uint64(vi), 16)
			case 'o':
				x.buf = strconv.AppendUint(x.buf, uint64(vi), 8)
			default:
				x.buf = strconv.AppendInt(x.buf, vi, 10)
			}
		} else if tFloat {
			switch spec {
			case 'f', 'e', 'E', 'g', 'G':
				x.buf = strconv.AppendFloat(x.buf, vf, spec, precision, 64)
			default:
				x.buf = strconv.AppendFloat(x.buf, vf, 'f', precision, 64)
			}
		}
	}
	return x
}
