package proto

import (
	"bufio"
	"fmt"
	"math/big"
	"strconv"
)

func (w *RESP) Read() (interface{}, error) {
	seg, err := w.readLine()
	if err != nil {
		return nil, err
	}

	switch typ := seg[0]; typ {
	case RespNil:
		return nil, Nil
	case RespError:
		return nil, RedisError(seg[1:])
	case RespStatus:
		return string(seg[1:]), nil
	case RespInt:
		return strconv.ParseInt(string(seg[1:]), 10, 64)
	case RespFloat:
		return strconv.ParseFloat(string(seg[1:]), 64)
	case RespBool:
		switch string(seg[1:]) {
		case "t":
			return true, nil
		case "f":
			return false, nil
		}
		return nil, UnexpectedData
	case RespBigInt:
		i := new(big.Int)
		if i, ok := i.SetString(string(seg[1:]), 10); ok {
			return i, nil
		}
		return i, nil
	case RespString:
		n, err := byteToInt(seg[1:])
		bs := make([]byte, n)
		rn, err := w.Reader.Read(bs)
		if err != nil {
			return nil, err
		} else if n != int64(rn) {
			return nil, UnexpectedData
		}

		return string(bs), nil
	case RespArray, RespSet, RespPush:
		n, err := byteToInt(seg[1:])
		if err != nil {
			return nil, err
		}
		return w.readSlice(n)
	case RespMap:
		n, err := byteToInt(seg[1:])
		if err != nil {
			return nil, err
		}
		return w.readMap(n)
	default:
		return nil, fmt.Errorf("redis: unknwon proto %s", string(typ))
	}
}

func (w *RESP) readLine() ([]byte, error) {
	b, err := w.Reader.ReadSlice('\n')
	if err != nil {
		if err != bufio.ErrBufferFull {
			return nil, err
		}

		full := make([]byte, len(b))
		copy(full, b)

		b, err = w.Reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		full = append(full, b...) //nolint:makezero
		b = full
	}
	if len(b) <= 2 || b[len(b)-1] != '\n' || b[len(b)-2] != '\r' {
		return nil, fmt.Errorf("redis: invalid reply: %q", b)
	}
	return b[:len(b)-2], nil
}

func (w *RESP) readSlice(n int64) ([]interface{}, error) {
	val := make([]interface{}, n)
	for i := 0; i < len(val); i++ {
		v, err := w.Read()
		if err != nil {
			if err == Nil {
				val[i] = nil
				continue
			}
			if err, ok := err.(RedisError); ok {
				val[i] = err
				continue
			}
			return nil, err
		}
		val[i] = v
	}
	return val, nil
}

func (w *RESP) readMap(n int64) (map[interface{}]interface{}, error) {
	m := make(map[interface{}]interface{}, n)
	for i := 0; i < len(m); i++ {
		k, err := w.Read()
		if err != nil {
			return nil, err
		}
		v, err := w.Read()
		if err != nil {
			if err == Nil {
				m[k] = nil
				continue
			}
			if err, ok := err.(RedisError); ok {
				m[k] = err
				continue
			}
			return nil, err
		}
		m[k] = v
	}
	return m, nil
}
