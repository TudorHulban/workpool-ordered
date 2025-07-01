package workpoolordered

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type payload []byte

func TestOneWorkerList(t *testing.T) {
	proc := func(payload payload) (payload, bool, error) {
		return payload, false, nil
	}

	list, errCr := NewCList(
		&ParamsCList[payload]{
			Processor:       proc,
			Workers:         1,
			WaitToStartWork: true,
		},
	)
	require.NoError(t, errCr)

	list.Insert(
		[]byte("a"),
	)
	list.Insert(
		[]byte("b"),
	)
	list.Insert(
		[]byte("c"),
	)

	require.EqualValues(t, 3, list.length.Load())
	require.EqualValues(t, 3, list.unprocessed.Load())

	time.Sleep(300 * time.Millisecond)

	elements := list.Read()
	require.NotEmpty(t, elements)

	// Should be processed in order: a, b, c (tail to head)
	require.Equal(t, []byte("a"), elements[0])
	require.Equal(t, []byte("b"), elements[1])
	require.Equal(t, []byte("c"), elements[2])

	fmt.Println("Processed elements:", elements)
}

func TestManyWorkersList(t *testing.T) {
	proc := func(payload payload) (payload, bool, error) {
		return payload, false, nil
	}

	list, errCr := NewCList(
		&ParamsCList[payload]{
			Processor: proc,
		},
	)
	require.NoError(t, errCr)

	list.Insert(
		[]byte("a"),
	)
	list.Insert(
		[]byte("b"),
	)
	list.Insert(
		[]byte("c"),
	)

	require.EqualValues(t, 3, list.length.Load())
	require.EqualValues(t, 3, list.unprocessed.Load())

	time.Sleep(300 * time.Millisecond)

	elements := list.Read()
	require.NotEmpty(t, elements)

	// Should be processed in order: a, b, c (tail to head)
	require.Equal(t, []byte("a"), elements[0])
	require.Equal(t, []byte("b"), elements[1])
	require.Equal(t, []byte("c"), elements[2])

	fmt.Println("Processed elements:", elements)
}
