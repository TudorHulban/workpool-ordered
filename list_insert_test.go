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

	processed := list.Read()

	require.NotEmpty(t,
		processed,
	)

	fmt.Println(
		list.unprocessed.Load(),
		len(processed),
	)
}
