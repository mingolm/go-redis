package proto

var (
	Nil            = RedisError("redis: nil")
	UnexpectedData = RedisError("redis: unexpected data")
)

type RedisError string

func (e RedisError) Error() string {
	return string(e)
}
