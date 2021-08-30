package util

type ByteSize uint64

// Common units for calculating size
const (
	BYTE ByteSize = 1
	KB            = 1024 * BYTE
	MB            = 1024 * KB
	GB            = 1024 * MB
	TB            = 1024 * GB
)

func (s ByteSize) Add(arg ByteSize) ByteSize {
	return s + arg
}

func (s ByteSize) Del(arg ByteSize) ByteSize {
	return s - arg
}

func (s ByteSize) Mul(arg ByteSize) ByteSize {
	return s * arg
}

func (s ByteSize) Div(arg ByteSize) ByteSize {
	return s / arg
}

func (s ByteSize) TB() ByteSize {
	return s / TB
}

func (s ByteSize) GB() ByteSize {
	return s / GB
}

func (s ByteSize) MB() ByteSize {
	return s / MB
}

func (s ByteSize) KB() ByteSize {
	return s / KB
}
