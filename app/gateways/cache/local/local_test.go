package local_test

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/test"
	"testing"

	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache/local"
)

func TestAddAndGet(t *testing.T) {
	test.ClearLocalCache()
	gw := local.GetInstance()

	key := "key1"
	val := test.ItemTest{
		ID:      1,
		Content: "content 1",
		Size:    45.67,
	}

	err := gw.Add(test.GetCtx(), key, val)
	if err != nil {
		t.Errorf("Add failed: %v", err)
	}

	res, err := gw.Get(test.GetCtx(), key)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if res != val {
		t.Errorf("expected %s, but got %s", test.ToString(val), res)
	}
}

func TestAddAllItems(t *testing.T) {
	gw := local.GetInstance()
	test.ClearLocalCache()

	data := map[string]interface{}{
		"key1": test.ItemTest{
			ID:      1,
			Content: "content 1",
			Size:    45.67,
		},
		"key2": test.ItemTest{
			ID:      2,
			Content: "content 2",
			Size:    23.02,
		},
	}

	err := gw.AddAllItems(context.Background(), data)
	if err != nil {
		t.Errorf("AddAllItems failed: %v", err)
	}

	for key, val := range data {
		res, err := gw.Get(test.GetCtx(), key)
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		if res != val {
			t.Errorf("expected %s, but got %s", val, res)
		}
	}
}

func TestGetAllKeysByPrefix(t *testing.T) {
	gw := local.GetInstance()
	test.ClearLocalCache()

	data := map[string]interface{}{
		"test:key:key1": test.ItemTest{
			ID:      1,
			Content: "content 1",
			Size:    45.67,
		},
		"test:key:key2": test.ItemTest{
			ID:      2,
			Content: "content 2",
			Size:    23.02,
		},
		"wrong:key:key2": test.ItemTest{
			ID:      3,
			Content: "content 3",
			Size:    173.87,
		},
	}

	err := gw.AddAllItems(test.GetCtx(), data)
	if err != nil {
		t.Errorf("AddAllItems failed: %v", err)
	}

	prefix := "test:key:"
	keys, err := gw.GetAllKeysByPrefix(test.GetCtx(), prefix)
	if err != nil {
		t.Errorf("GetAllKeysByPrefix failed: %v", err)
	}

	exp := []string{"test:key:key1", "test:key:key2"}
	for _, key := range exp {
		if !test.ContainsItem(key, keys) {
			t.Errorf("expected to contain %v, but it did not", key)
		}
	}
}

func TestGetAllItems(t *testing.T) {
	gw := local.GetInstance()
	test.ClearLocalCache()

	data := map[string]interface{}{
		"test:key:key1": test.ItemTest{
			ID:      1,
			Content: "content 1",
			Size:    45.67,
		},
		"test:key:key2": test.ItemTest{
			ID:      2,
			Content: "content 2",
			Size:    23.02,
		},
	}

	err := gw.AddAllItems(test.GetCtx(), data)
	if err != nil {
		t.Errorf("AddAllItems failed: %v", err)
	}

	items, err := gw.GetAllItems(test.GetCtx())
	if err != nil {
		t.Errorf("GetAllItems failed: %v", err)
	}

	if !test.MapEqual(items, data) {
		t.Errorf("expected %v, but got %v", data, items)
	}
}

func TestDelete(t *testing.T) {
	gw := local.GetInstance()
	test.ClearLocalCache()

	data := map[string]interface{}{
		"test:key:key1": test.ItemTest{
			ID:      1,
			Content: "content 1",
			Size:    45.67,
		},
		"test:key:key2": test.ItemTest{
			ID:      2,
			Content: "content 2",
			Size:    23.02,
		},
	}

	delKey := "test:key:key1"

	err := gw.AddAllItems(test.GetCtx(), data)
	if err != nil {
		t.Errorf("AddAllItems failed: %v", err)
	}

	err = gw.Delete(test.GetCtx(), delKey)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	_, err = gw.Get(test.GetCtx(), delKey)
	if err == nil {
		t.Errorf("expected error for Get after deletion, but got nil")
	}
}
