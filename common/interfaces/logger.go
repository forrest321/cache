package interfaces

// Logger interface for logging messages
type Logger interface {

	// Printf formats and prints the given format string with optional arguments.
	//
	// The format string can contain verbs for formatting different types. The
	// verbs are derived from the fmt package in Go standard library. The
	// arguments can be passed using an optional variadic argument of type
	// interface{}. The method will use these arguments to replace any verb
	// placeholders in the format string. The formatted string will be printed
	// using the underlying logger implementation.
	//
	// Example usage:
	//   logger := log.New(os.Stdout, "cache: ", log.LstdFlags)
	//   logger.Printf("Active cleanup: expired item removed: %s", key)
	//   logger.Printf("Lazy cleanup: expired item removed: %s", key)
	//   logger.Printf("Value: %d, Error: %v", value, err)
	Printf(format string, v ...interface{})

	// Println formats its arguments and writes to standard output.
	// It receives a variadic parameter, v of type ...interface{} which represents the values to be printed.
	// The function does not return any value.
	Println(v ...interface{})

	// Fatal logs a fatal error message and exits the program
	Fatal(v ...any)
}
