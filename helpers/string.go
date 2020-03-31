package helpers

// StringPointers converts a slice of string values into a slice of string pointers
//
// This function complements aws.StringSlice but works with variadic arguments so that an array literal is not required.
func StringPointers(strings ...string) []*string {
	sp := make([]*string, len(strings))
	for i := range sp {
		sp[i] = &strings[i]
	}
	return sp
}
