package htmlq

type stringWriter struct {
	s string
}

func (sw *stringWriter) Write(bs []byte) (int, error) {
	sw.s += string(bs)
	return len(bs), nil
}

func (sw *stringWriter) String() string {
	return sw.s
}
