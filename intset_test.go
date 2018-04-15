package intset

import (
	"fmt"
	"testing"
)

var (
	fibonacci = []int{0, 1, 2, 3, 5, 8}     // Fibonacci numbers < 10
	primes    = []int{2, 3, 5, 7}           // prime numbers < 10
	union     = []int{0, 1, 2, 3, 5, 7, 8}  // fibonnacci ∪ prime
	empty     = []int{}
	single    = []int{100}  // set with element > 63
)

func TestAdd(t *testing.T) {
	s := IntSet{}
	s.Add(13)
	for _, n := range fibonacci {
		s.Add(n)
	}
	want := "{0 1 2 3 5 8 13}"
	got := s.String()
	if want != got {
		t.Errorf("%s != %s", got, want)
	}
}

func TestFromSlice(t *testing.T) {
	s := FromSlice(primes)
	want := "{2 3 5 7}"
	got := s.String()
	if want != got {
		t.Errorf("%s != %s", got, want)
	}
}

func TestLen(t *testing.T) {
	s := FromSlice(fibonacci)
	want := 6
	got := s.Len()
	if want != got {
		t.Errorf("%d != %d", got, want)
	}
}

func TestLen_Union(t *testing.T) {
	want := FromSlice(union).Len()
	set := FromSlice(primes)
	set.UnionWith(FromSlice(fibonacci))
	got := set.Len()
	if want != got {
		t.Errorf("%d != %d", got, want)
	}

}

func TestLen_UnionSecondWord(t *testing.T) {
	// UnionWith set that uses 2 words: {100}
	set := FromSlice(fibonacci)
	set.UnionWith(FromSlice(single))
	want := 7
	got := set.Len()
	if want != got {
		t.Errorf("%d != %d", got, want)
	}

}

func TestHas(t *testing.T) {
	var testCases = []struct {
		n    int
		set  *IntSet
		want bool
	}{
		{3, FromSlice(primes), true},
		{4, FromSlice(primes), false},
		{1, FromSlice(empty), false},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d in %v", tc.n, tc.set), func(t *testing.T) {
			got := tc.set.Has(tc.n)
			if got != tc.want {
				t.Errorf("tc.set.Has(%d) = %v; want %v", tc.n, got, tc.want)
			}
		})
	}
}

func TestBitCount(t *testing.T) {
	var testCases = []struct {
		word uint64
		want int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{10, 2},
		{11, 3},
		{12, 2},
		{13, 3},
		{15, 4},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%b has %d", tc.word, tc.want), func(t *testing.T) {
			got := bitCount(tc.word)
			if got != tc.want {
				t.Errorf("bitCount(%b) = %d; want %d", tc.word, got, tc.want)
			}
		})
	}
}

func TestElems(t *testing.T) {
	var testCases = []struct {
		set  *IntSet
		want []int
	}{
		{FromSlice(empty), empty},
		{FromSlice(single), single},
		{FromSlice(primes), primes},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v elems %v", tc.set, tc.want), func(t *testing.T) {
			got := tc.set.Elems()
			equal := true;
			if len(got) != len(tc.want) {
				equal = false
			} else {
				for i, v := range got {
					if v != tc.want[i] {
						equal = false
						break
					}
				}
			}
			if !equal {
				t.Errorf("%v.Elems() = %v; want %v", tc.set, got, tc.want)
			}
		})
	}
}

func Example() {
	var x, y IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	fmt.Println(x.String()) // "{1 9 144}"

	y.Add(9)
	y.Add(42)
	fmt.Println(y.String()) // "{9 42}"

	x.UnionWith(&y)
	fmt.Println(x.String()) // "{1 9 42 144}"

	fmt.Println(x.Has(9), x.Has(123)) // "true false"

	// Output:
	// {1 9 144}
	// {9 42}
	// {1 9 42 144}
	// true false
}

func Example_string_representations() {
	var x IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	x.Add(42)

	fmt.Println(&x)         // "{1 9 42 144}"
	fmt.Println(x.String()) // "{1 9 42 144}"
	fmt.Println(x)          // "{[4398046511618 0 65536] 4}"

	// Output:
	// {1 9 42 144}
	// {1 9 42 144}
	// {[4398046511618 0 65536] 4}
}
