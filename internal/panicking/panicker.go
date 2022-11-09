package panicking

// Panicker abstracts panic for our tests and such.
type Panicker interface {
	Panic(any)
	Panicf(format string, args ...any)
}
