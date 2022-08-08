package config

import "fmt"

// FileNotFoundError denotes failing to find configuration file.
type FileNotFoundError struct {
	name, locations string
}

// Error returns the formatted configuration error.
func (e FileNotFoundError) Error() string {
	return fmt.Sprintf("Config File %q Not Found in %q", e.name, e.locations)
}

// UnsupportedConfigError denotes encountering an unsupported
// configuration file type.
type UnsupportedConfigError string

// Error returns the formatted configuration error.
func (str UnsupportedConfigError) Error() string {
	return fmt.Sprintf("Unsupported Config Type %q", string(str))
}

// InvalidConfigTypeError denotes an invalid
// configuration type
type InvalidConfigTypeError string

// Error returns the formatted configuration error.
func (str InvalidConfigTypeError) Error() string {
	return string(str)
}

// ParseError denotes failing to parse configuration file.
type ParseError struct {
	err error
}

// Error returns the formatted configuration error.
func (pe ParseError) Error() string {
	return fmt.Sprintf("While parsing config: %s", pe.err.Error())
}
