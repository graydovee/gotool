package utils

import (
	"fmt"
	"testing"
)

type KV struct {
	K string `json:"key"`
	v string `json:"value"`
}

func TestNewSet(t *testing.T) {

	s1 := NewSet[KV]()
	s2 := NewSet[KV]()
	s3 := NewSet[KV]()

	s1.Add(KV{"1", "a"})
	s1.Add(KV{"2", "b"})

	s2.Add(KV{"2", "b"})
	s2.Add(KV{"3", "c"})

	s3.Add(KV{"1", "a"})
	s3.Add(KV{"2", "b"})
	s3.Add(KV{"3", "c"})

	print(s1.Union(s2), t)

	print(s3.Difference(s2), t)

	print(s1.Difference(s2), t)

	print(s1.SymmetricDifference(s2), t)

	print(s1.Intersect(s2), t)
}

func print(set Set[KV], t *testing.T) {
	json, err := set.MarshalJSON()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Println(string(json))
}
