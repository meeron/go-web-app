package shared

func PanicOnErr(err error) {
	if err == nil {
		return
	}

	panic(err)
}

func Unwrap[T any](result T, err error) T {
	PanicOnErr(err)
	return result
}
