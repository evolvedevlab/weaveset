package store

import (
	"testing"

	"github.com/evolvedevlab/weaveset/data"
)

func TestStore_Save(t *testing.T) {
	s := NewFileSystem("../test") // TODO: fix stuff
	err := s.Save(&data.List{
		ID:   "someid",
		Name: "Best 2025 Novels to read this summer!",
	})
	if err != nil {
		t.Error(err)
	}
}
