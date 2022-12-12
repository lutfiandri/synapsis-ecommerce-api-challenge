package helper

func PanicIfNeeded(err error) {
	if err != nil {
		panic(err)
	}
}
