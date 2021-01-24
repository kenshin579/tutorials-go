package go_testing

import "io"

// readN reads at most n bytes from r and returns them as a string.
func readN(r io.Reader, n int) (string, error) {
	buf := make([]byte, n)
	m, err := r.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:m]), nil
}
