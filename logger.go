package main

type Logger struct {
	internal string
}

func (l *Logger) Write(bs []byte) (int, error) {
	l.internal += string(bs)
	return len(bs), nil
}
func (l *Logger) String() string {
	return l.internal
}
