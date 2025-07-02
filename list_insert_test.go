package workpoolordered

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {
	proc := func(payload []byte) ([]byte, bool, error) {
		return payload, false, nil
	}

	list, errCr := NewCList(
		&ParamsCList[[]byte]{
			Processor: proc,
			Workers:   1,
		},
	)
	require.NoError(t, errCr)

	numberItems := 100000

	startInsert := time.Now()

	for ix := range numberItems {
		list.Insert([]byte(fmt.Sprintf("item-%d", ix)))
	}

	fmt.Println(
		time.Since(startInsert), // ~30ms
	)

	require.EqualValues(t,
		numberItems,
		list.length.Load(),
	)

	require.Less(t,
		list.unprocessed.Load(),
		list.length.Load(),
	)

	batchSize := 10

	processed := list.Read(batchSize)

	require.Equal(t,
		batchSize,
		len(processed),

		"returned %d items",
		len(processed),
	)

	fmt.Printf(
		"unprocessed: %d, processed returned: %d, processed all: %d.\n",
		list.unprocessed.Load(),
		len(processed),
		int64(numberItems)-list.unprocessed.Load(),
	)
}

func TestUnprocessed(t *testing.T) {
	proc := func(payload []byte) ([]byte, bool, error) {
		return payload, false, nil
	}

	list, errCr := NewCList(
		&ParamsCList[[]byte]{
			Processor:       proc,
			Workers:         1,
			WaitToStartWork: true,
		},
	)
	require.NoError(t, errCr)

	list.Insert(
		[]byte("xxx"),
	)

	require.EqualValues(t,
		1,
		list.length.Load(),
	)

	require.EqualValues(t,
		1,
		list.unprocessed.Load(),
	)
}
