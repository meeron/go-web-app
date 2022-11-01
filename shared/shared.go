package shared

func PanicOnErr(err error) {
	if err == nil {
		return
	}

	panic(err)
}
