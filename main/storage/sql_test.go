package storage

import (
	"os"
	"path"
	"testing"
)

func TestGetDB(t *testing.T) {
	db, err := GetDB(testCfg)
	if err != nil {
		t.Fatal("can not create db", err)
	}
	defer db.Close()
	defer os.Remove(
		path.Join(
			testCfg.DataDir, testCfg.Database.Name+".db"))
}
