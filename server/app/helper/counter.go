package helper

// Counter is a simple writer that counts the number of bytes written to it.
type Counter struct {
	Total *int64
}

// Write implements the io.Writer interface.
// It adds the number of bytes in p to the total.
func (c *Counter) Write(p []byte) (int, error) {
	*c.Total += int64(len(p))
	return len(p), nil
}
