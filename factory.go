package logger

func Default() Interface {
	return New(Config{})
}

func Fake() Interface {
	return fake{}
}
