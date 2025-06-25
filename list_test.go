package workpoolordered

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOneWorkerList(t *testing.T) {
	proc := func(payload []byte) ([]byte, error) {
		return payload, nil
	}

	list := NewDLinkedList(proc, 1)

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

	go list.Process(context.Background())

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
	proc := func(payload []byte) ([]byte, error) {
		return payload, nil
	}

	list := NewDLinkedList(proc, 4)

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

	go list.Process(context.Background())

	time.Sleep(300 * time.Millisecond)

	elements := list.Read()
	require.NotEmpty(t, elements)

	// Should be processed in order: a, b, c (tail to head)
	require.Equal(t, []byte("a"), elements[0])
	require.Equal(t, []byte("b"), elements[1])
	require.Equal(t, []byte("c"), elements[2])

	fmt.Println("Processed elements:", elements)
}
