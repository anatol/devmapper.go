package devmapper

// ZeroTable represents information needed for 'zero' target creation
type ZeroTable struct {
	Start  uint64
	Length uint64
}

func (z ZeroTable) start() uint64 {
	return z.Start
}

func (z ZeroTable) length() uint64 {
	return z.Length
}

func (z ZeroTable) buildSpec() string {
	return ""
}

func (z ZeroTable) targetType() string {
	return "zero"
}
