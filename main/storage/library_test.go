package storage

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path"
	"testing"
)

func TestGetLib(t *testing.T) {
	rootDir := testCfg.Library
	err := os.MkdirAll(rootDir, 0775)
	if err != nil {
		t.Fatal(err)
	}
	dirs1 := []string{"1", "2", "3"}
	dirs2 := []string{"11", "22"}
	dirs3 := []string{"111"}

	mkdir := func(n string) {
		err := os.Mkdir(path.Join(rootDir, n), 0775)
		if err != nil {
			t.Fatal(err)
		}
	}
	verifyDir := func(want, got string) {
		if want != got {
			t.Fatalf("want %s, got %s", want, got)
		}
	}
	for _, d1 := range dirs1 {
		mkdir(d1)
		for _, d2 := range dirs2 {
			mkdir(path.Join(d1, d2))
			for _, d3 := range dirs3 {
				mkdir(path.Join(d1, d2, d3))
			}
		}
	}

	lib, err := GetLib(rootDir)
	if err != nil {
		t.Fatal(err)
	}
	verifyDir(rootDir, lib.Dir)
	lib = lib.SubLib["1"]
	verifyDir("1", lib.Name)
}
func TestLibrary(t *testing.T) {
	rootDir := testCfg.Library
	err := os.MkdirAll(rootDir, 0775)
	if err != nil {
		t.Fatal(err)
	}
	lib, err := GetLib(testCfg.Library)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.CreateTemp("", "test-temp")
	if err != nil {
		t.Fatal(err)
	}
	temp := []byte("I an a PNG file")
	_, err = f.Write(temp)
	if err != nil {
		t.Fatal(err)
	}
	gotSum, err := lib.Add(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	setu, err := lib.GetFile(gotSum)
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(setu)
	if err != nil {
		t.Fatal(err)
	}
	for i, n := range b {
		if temp[i] != n {
			t.Fatalf("Failed to verify temp data\nwant: %v\ngot: %v", temp, b)
		}
	}
	wantSum := fmt.Sprintf("%x", md5.Sum(temp))
	if wantSum != gotSum {
		t.Fatalf("Diffrent md5 checksum \nwant: %v\ngot: %v", temp, b)
	}
}
