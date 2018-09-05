package defaults

// String returns s1 and if empty then returns s2
func String(s1, s2 string) string {
	if s1 == "" {
		return s2
	}
	return s1
}

// Int returns i1 if != 0 othervise returns i2
func Int(i1, i2 int) int {
	if i1 == 0 {
		return i2
	}
	return i1
}

// Int64 returns i1 if != 0 othervise returns i2
func Int64(i1, i2 int64) int64 {
	if i1 == 0 {
		return i2
	}
	return i1
}

// Float32 returns i1 if != 0 othervise returns i2
func Float32(i1, i2 float32) float32 {
	if i1 == 0.0 {
		return i2
	}
	return i1
}

// Float64 returns i1 if != 0 othervise returns i2
func Float64(i1, i2 float64) float64 {
	if i1 == 0.0 {
		return i2
	}
	return i1
}
