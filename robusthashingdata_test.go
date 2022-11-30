package care

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func slices(t *testing.T) {
	t.Parallel()
	const items = 100
	s := make([]int, 0, items)
	for i := 0; i < items; i++ {
		s = append(s, i)
	}

	rand.Seed(time.Now().UnixNano())

	b := strings.Builder{}
	set := make(map[string]struct{})
	const attempts = 100

	for k := 0; k < attempts; k++ {
		rand.Shuffle(items, func(i, j int) { s[i], s[j] = s[j], s[i] })
		err := robustHashingData(s, &b)
		if err != nil {
			t.Fatalf("Failed to call robustHashingData. %v", err)
		}
		data := b.String()
		set[data] = struct{}{}
		b.Reset()
	}

	if len(set) != 1 {
		t.Error("Failed to make data for robust hashing from slices.")
	}
}

func maps(t *testing.T) {
	t.Parallel()
	const items = 100
	m := make(map[int][]int)
	for i := 0; i < items; i++ {
		m[i] = make([]int, 0, items)
		for j := 0; j < items; j++ {
			m[i] = append(m[i], j)
		}
	}

	rand.Seed(time.Now().UnixNano())

	b := strings.Builder{}
	set := make(map[string]struct{})
	const attempts = 100

	for l := 0; l < attempts; l++ {
		for k := 0; k < items; k++ {
			rand.Shuffle(items, func(i, j int) { m[k][i], m[k][j] = m[k][j], m[k][i] })
		}
		err := robustHashingData(m, &b)
		if err != nil {
			t.Fatalf("Failed to call robustHashingData. %v", err)
		}
		data := b.String()
		set[data] = struct{}{}
		b.Reset()
	}

	if len(set) != 1 {
		t.Error("Failed to make data for robust hashing from maps.")
	}
}

func structs(t *testing.T) {
	t.Parallel()

	type nested struct {
		field1 int
		field2 float32
		field3 *int
		field4 map[int]int
	}

	testStruct := struct {
		field1 string
		field2 string
		field3 int
		field4 []int
		field5 [2]bool
		field6 map[string]string
		field7 map[string][]nested
	}{
		field1: "String 1",
		field2: "String 2",
		field3: 100500,
		field4: []int{1, 2, 3, 4, 5},
		field5: [2]bool{true, false},
		field6: map[string]string{
			"s1": "s1",
			"s2": "s2",
			"s3": "s3",
		},
		field7: map[string][]nested{
			"item1": {nested{
				field1: 1,
				field2: 3.14,
				field3: nil,
				field4: map[int]int{
					1: 1,
					2: 2,
					3: 3,
				},
			}},
			"item2": {nested{
				field1: 2,
				field2: 9.81,
				field3: nil,
				field4: map[int]int{
					33: 55,
					77: 55,
				},
			},
				nested{
					field1: 3,
					field2: 22 / 7.0,
					field3: nil,
					field4: map[int]int{
						100: 500,
						500: 100},
				}},
		},
	}

	rand.Seed(time.Now().UnixNano())

	b := strings.Builder{}
	set := make(map[string]struct{})
	const attempts = 100

	for k := 0; k < attempts; k++ {
		for _, v := range testStruct.field7 {
			rand.Shuffle(len(v), func(i, j int) { v[i], v[j] = v[j], v[i] })
		}

		err := robustHashingData(testStruct, &b)
		if err != nil {
			t.Fatalf("Failed to call robustHashingData. %v", err)
		}
		data := b.String()
		set[data] = struct{}{}
		b.Reset()
	}

	if len(set) != 1 {
		t.Error("Failed to make data for robust hashing from structs.")
	}
}

func Test_robustHashingData(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{name: "slices", test: slices},
		{name: "maps", test: maps},
		{name: "structs", test: structs},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
