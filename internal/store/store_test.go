package store

import (
	"errors"
	"testing"

	"github.com/iov-one/cosmos-sdk-crud/internal/store/types"
	"github.com/iov-one/cosmos-sdk-crud/internal/test"
)

func TestStore(t *testing.T) {
	db, cdc, err := test.NewStore()
	if err != nil {
		t.Fatal("failed precondition", err)
	}
	s := NewStore(cdc, db, nil)
	obj := test.NewDeterministicObject()
	// test create
	err = s.Create(obj)
	if err != nil {
		t.Fatal(err)
	}
	// test read
	var expected = test.NewObject()
	err = s.Read(obj.PrimaryKey(), expected)
	if err != nil {
		t.Fatal(err)
	}
	if err := obj.Equals(expected); err != nil {
		t.Fatal(err)
	}
	// test update
	update := obj
	update.TestSecondaryKeyB = []byte("test-update")
	err = s.Update(update)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Read(obj.PrimaryKey(), expected)
	if err != nil {
		t.Fatal(err)
	}
	if err := update.Equals(expected); err != nil {
		t.Fatal(err)
	}
	// test cursor
	crs, err := s.Query([]types.SecondaryKey{
		update.FirstSecondaryKey(),
	}, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	// test read
	err = crs.Read(expected)
	if err != nil {
		t.Logf("%s", crs.currKey())
		t.Fatal(err)
	}
	if err := update.Equals(expected); err != nil {
		t.Fatal(err)
	}
	// test update
	update.TestSecondaryKeyA = []byte("another-update")
	err = crs.Update(update)
	if err != nil {
		t.Fatal(err)
	}
	expected.Reset()
	err = crs.Read(expected)
	if err := update.Equals(expected); err != nil {
		t.Fatal(err)
	}
	// test delete
	err = crs.Delete()
	if err != nil {
		t.Fatal(err)
	}
	err = crs.Read(expected)
	if !errors.Is(err, types.ErrNotFound) {
		t.Fatal("unexpected error", err)
	}
}
