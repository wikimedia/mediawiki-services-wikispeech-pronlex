package line

func equals(expect map[Field]string, result map[Field]string) bool {
	if len(expect) != len(result) {
		return false
	}
	for f, expS := range expect {
		resS := result[f]
		if resS != expS {
			return false
		}
	}
	return true
}

// Equals compares two line.Format instances
func (f Format) Equals(other Format) bool {
	if f.Name != other.Name {
		return false
	}
	if f.FieldSep != other.FieldSep {
		return false
	}
	if f.NFields != other.NFields {
		return false
	}
	if len(f.Fields) != len(other.Fields) {
		return false
	}
	for f, expS := range f.Fields {
		resS := other.Fields[f]
		if resS != expS {
			return false
		}
	}
	return true
}

type stringSlice []string

func (a stringSlice) Len() int      { return len(a) }
func (a stringSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a stringSlice) Less(i, j int) bool {
	return a[i] < a[j]
}
