package diff

// Op represents the type of diff operation
type Op int

const (
	OpUnchanged Op = iota
	OpAdded          // New key/element in target
	OpRemoved        // Key/element removed from source
	OpChanged        // Value changed between source and target
	OpMoved          // Array element moved to a different index
)

// String returns the display name of the operation
func (o Op) String() string {
	switch o {
	case OpAdded:
		return "ADDED"
	case OpRemoved:
		return "REMOVED"
	case OpChanged:
		return "CHANGED"
	case OpMoved:
		return "MOVED"
	default:
		return "UNCHANGED"
	}
}

// DiffItem represents a single difference between two JSON values
type DiffItem struct {
	Op       Op
	Path     string
	OldValue interface{}
	NewValue interface{}
	OldIndex int // Only for OpMoved
	NewIndex int // Only for OpMoved
}
