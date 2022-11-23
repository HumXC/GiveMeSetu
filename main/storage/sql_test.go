package storage

import (
	"testing"
)

func TestGetDB(t *testing.T) {
	db, err := GetDB(testCfg)
	if err != nil {
		t.Fatal("can not create db", err)
	}
	defer db.Close()
}

func TestLibDB(t *testing.T) {
	items := []Setu{
		{
			ID:     "id1",
			Title:  "setu1",
			Ext:    ".png",
			Origin: "TestLibDB()",
		},
		{
			ID:     "id2",
			Title:  "setu2",
			Ext:    ".jpg",
			Origin: "TestLibDB()",
		},
		{
			ID:     "id3",
			Title:  "setu3",
			Ext:    ".mp3",
			Origin: "TestLibDB()",
		},
		{
			ID:     "id4",
			Title:  "setu4",
			Ext:    ".mp4",
			Origin: "TestLibDB()",
		},
		{
			ID:     "id5",
			Title:  "setu5",
			Ext:    ".png",
			Origin: "TestLibDB()",
		},
	}
	db, err := GetDB(testCfg)
	if err != nil {
		t.Fatal("can not create db", err)
	}
	setuDB := db.SetuDB
	defer db.Close()
	for _, item := range items {
		err := setuDB.Add(item)
		if err != nil {
			t.Fatalf("faild to add item: %v - ,%s", item, err)
		}
	}

	setu, err := setuDB.GetByID("id5")
	if err != nil {
		t.Fatal(err)
	}
	if setu.ID != "id5" {
		t.Fatalf("want id5, got [%s]", setu.ID)
	}
	// bad query
	_, err = setuDB.GetByID("id6")
	if err == nil {
		t.Fatal("There should be an error, but there isn't")
	}

	err = setuDB.Del("id5")
	if err != nil {
		t.Fatal(err)
	}

	_, err = setuDB.GetByID("id5")
	if err == nil {
		t.Fatal("There should be an error, but there isn't")
	}

	setus, err := setuDB.GetByIDs([]string{"id1", "id2"})
	if err != nil {
		t.Fatal(err)
	}
	if setus[0].ID != "id1" && setus[1].ID != "id2" && setus[0].Ext != items[0].Ext {
		t.Fatal("bad query result")
	}
}
