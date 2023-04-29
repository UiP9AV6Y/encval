package log

const (
	OFF Verbosity = iota
	ERROR
	WARNING
	INFO
	DEBUG
	TRACE
)

var verbosityNames = []string{
	"OFF",
	"ERROR",
	"WARNING",
	"INFO",
	"DEBUG",
	"TRACE",
}

func (v Verbosity) Validate() bool {
	return int(v) < len(verbosityNames)
}

// Compare returns -1, 0, or +1 depending on whether v is 'less than',
// 'equal to', or 'greater than' the other verbosity.
func (v Verbosity) Compare(other Verbosity) int {
	if v == other {
		return 0
	}

	if v < other {
		return -1
	}

	return 1
}

func (v Verbosity) String() string {
	i := len(verbosityNames)
	if int(v) < i {
		return verbosityNames[v]
	}

	return verbosityNames[i-1]
}
