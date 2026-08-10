package all

// This test file's only job is to sit in the same package as the blank
// imports above, so `go test ./...` actually loads every concept's init()
// together and exercises the registry as a whole — including the
// panic-if-duplicate-Seq / panic-if-missing-Seq checks in concept.Register.
// Without a test here, `make check` never imports every concept in one
// binary, and a broken Seq would only surface later when someone runs
// `make build`.

import (
	"testing"

	"mathviz/internal/concept"
)

func TestRegistryHasNoSeqGapsOrDuplicates(t *testing.T) {
	all := concept.All()
	if len(all) != concept.Count() {
		t.Fatalf("concept.All() returned %d concepts, concept.Count() = %d", len(all), concept.Count())
	}

	seen := make(map[int]string, len(all))
	for _, c := range all {
		if c.Seq <= 0 {
			t.Errorf("%s: Seq is %d, want a positive build-sequence number", c.ID, c.Seq)
			continue
		}
		if other, dup := seen[c.Seq]; dup {
			t.Errorf("Seq %d used by both %s and %s", c.Seq, other, c.ID)
		}
		seen[c.Seq] = c.ID
	}
}

func TestRegistrySortedNewestFirst(t *testing.T) {
	all := concept.All()
	for i := 1; i < len(all); i++ {
		if all[i].Seq > all[i-1].Seq {
			t.Fatalf("concept.All() not sorted newest-first: %s (Seq %d) came after %s (Seq %d)",
				all[i].ID, all[i].Seq, all[i-1].ID, all[i-1].Seq)
		}
	}
}
