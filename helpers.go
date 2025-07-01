package workpoolordered

import "bytes"

func NotInComparableSliceSource[T comparable](source []T, elements ...T) []T {
	result := make([]T, 0)

outer:
	for _, element := range elements {
		for _, elemSource := range source {
			if elemSource == element {
				continue outer
			}
		}

		result = append(result, element)
	}

	return result
}

func NotInByteSliceSource(source [][]byte, elements ...[]byte) [][]byte {
	result := make([][]byte, 0)

outer:
	for _, element := range elements {
		for _, elemSource := range source {
			if bytes.Equal(elemSource, element) {
				continue outer
			}
		}
		result = append(result, element)
	}

	return result
}
