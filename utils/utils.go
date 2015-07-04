package utils

func SubmitAsync(e error, c chan error) {
	go func() {
		c <- e
	}()
}
