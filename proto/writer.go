package proto

import (
	"context"
	"encoding"
	"fmt"
	"net"
	"strconv"
	"time"
)

func (w *RESP) Write(ctx context.Context, args []interface{}) error {
	// *n
	bs := append(strconv.AppendUint([]byte{RespArray}, uint64(len(args)), 10), '\r', '\n')
	if _, err := w.Writer.Write(bs); err != nil {
		return err
	}

	// arg
	for _, arg := range args {
		if err := w.writeArg(arg); err != nil {
			return err
		}
	}

	return nil
}

func (w *RESP) writeArg(v interface{}) error {
	switch v := v.(type) {
	case nil:
		return w.writeString("")
	case string:
		return w.writeString(v)
	case []byte:
		return w.writeBytes(v)
	case int:
		return w.writeInt(int64(v))
	case int8:
		return w.writeInt(int64(v))
	case int16:
		return w.writeInt(int64(v))
	case int32:
		return w.writeInt(int64(v))
	case int64:
		return w.writeInt(v)
	case uint:
		return w.writeUint(uint64(v))
	case uint8:
		return w.writeUint(uint64(v))
	case uint16:
		return w.writeUint(uint64(v))
	case uint32:
		return w.writeUint(uint64(v))
	case uint64:
		return w.writeUint(v)
	case float32:
		return w.writeFloat(float64(v))
	case float64:
		return w.writeFloat(v)
	case bool:
		var bi int64
		if v {
			bi = 1
		}
		return w.writeInt(bi)
	case time.Time:
		return w.writeBytes(v.AppendFormat([]byte{}, time.RFC3339Nano))
	case time.Duration:
		return w.writeInt(v.Nanoseconds())
	case encoding.BinaryMarshaler:
		b, err := v.MarshalBinary()
		if err != nil {
			return err
		}
		return w.writeBytes(b)
	case net.IP:
		return w.writeBytes(v)
	default:
		return fmt.Errorf("redis: can't marshal %T (implement encoding.BinaryMarshaler)", v)
	}
}

func (w *RESP) writeBytes(b []byte) error {
	bs := strconv.AppendInt([]byte{RespString}, int64(len(b)), 10)
	bs = append(bs, b...)
	bs = append(bs, '\r', '\n')
	if _, err := w.Writer.Write(bs); err != nil {
		return err
	}
	return nil
}

func (w *RESP) writeString(s string) error {
	return w.writeBytes([]byte(s))
}

func (w *RESP) writeUint(n uint64) error {
	bs := strconv.AppendUint([]byte{}, n, 10)
	return w.writeBytes(bs)
}

func (w *RESP) writeInt(n int64) error {
	bs := strconv.AppendInt([]byte{}, n, 10)
	return w.writeBytes(bs)
}

func (w *RESP) writeFloat(f float64) error {
	bs := strconv.AppendFloat([]byte{}, f, 'f', -1, 64)
	return w.writeBytes(bs)
}
