package devmapper

// ZeroTable represents information needed for 'zero' target creation
type ZeroTable struct {
	StartSector, Length uint64
}

func (z ZeroTable) startSector() uint64 {
	return z.StartSector
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
