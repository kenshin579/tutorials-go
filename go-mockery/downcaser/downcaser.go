package downcaser

type downcaser interface {
	Downcase(string) (string, error)
}
