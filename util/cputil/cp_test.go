package cputil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

type Dummy struct {
	A1 int
	A2 int
	B1 bool
	B2 bool
	C1 map[string]interface{}
	D1 []int
}

func TestShallowCopy(t *testing.T) {
	d1 := &Dummy{
		A1: 0,
		A2: 1,
		B1: false,
		B2: true,
		C1: map[string]interface{}{
			"a": 1,
		},
		D1: []int{0, 1},
	}

	d2 := &Dummy{}
	_ = ShallowCopy(&d2, d1)
	require.Empty(t, cmp.Diff(d1, d2))

	// case: 浅拷贝会互相影响
	d1.C1["a"] = 2
	require.Equal(t, d2.C1["a"], d1.C1["a"])
}

func TestDeepCopy(t *testing.T) {
	d1 := &Dummy{
		A1: 0,
		A2: 1,
		B1: false,
		B2: true,
		C1: map[string]interface{}{
			"a": 1,
		},
		D1: []int{0, 1},
	}

	d2 := &Dummy{}
	_ = DeepCopy(&d2, d1)
	require.Empty(t, cmp.Diff(d1, d2))

	// case: 深拷贝不会互相影响
	d1.C1["a"] = 2
	require.NotEqual(t, d2.C1["a"], d1.C1["a"])
}
