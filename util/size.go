package util

type SIZE uint64

// 定义文件单位
const (
	BYTE SIZE = 1
	KB   SIZE = 1024 * BYTE
	MB   SIZE = 1024 * KB
	GB   SIZE = 1024 * MB
	TB   SIZE = 1024 * GB
)

func (s SIZE) Add(arg SIZE) SIZE {
	return s + arg
}

func (s SIZE) Del(arg SIZE) SIZE {
	return s - arg
}

func (s SIZE) Mul(arg SIZE) SIZE {
	return s * arg
}

func (s SIZE) Div(arg SIZE) SIZE {
	return s / arg
}

func (s SIZE) TB() SIZE {
	return s / TB
}

func (s SIZE) GB() SIZE {
	return s / GB
}

func (s SIZE) MB() SIZE {
	return s / MB
}

func (s SIZE) KB() SIZE {
	return s / KB
}
