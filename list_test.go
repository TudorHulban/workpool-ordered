package workpoolordered

import (
	"context"
	"fmt"
	"testing"

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

	require.Equal(t, 3, list.length)
	require.Equal(t, 3, list.unprocessed)

	require.NoError(t,
		list.Process(context.Background()),
	)
	require.Zero(t, list.unprocessed)

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

	require.Equal(t, 3, list.length)
	require.Equal(t, 3, list.unprocessed)

	require.NoError(t,
		list.Process(context.Background()),
	)
	require.Zero(t, list.unprocessed)

	elements := list.Read()
	require.NotEmpty(t, elements)

	// Should be processed in order: a, b, c (tail to head)
	require.Equal(t, []byte("a"), elements[0])
	require.Equal(t, []byte("b"), elements[1])
	require.Equal(t, []byte("c"), elements[2])

	fmt.Println("Processed elements:", elements)
}
