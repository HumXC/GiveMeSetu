package storage

import (
	"crypto/md5"
	"fmt"
	"give-me-setu/util"
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
	verifyDir(path.Join(rootDir, "1"), lib.SubLib["1"].Dir)
	lib = lib.SubLib["1"]
	verifyDir("1", lib.Name)
}
func TestLibrary(t *testing.T) {
	createTemp := func(s string) (string, []byte) {
		f, err := os.CreateTemp("", "test-temp")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		temp := []byte(s)
		_, err = f.Write(temp)
		if err != nil {
			t.Fatal(err)
		}
		return f.Name(), temp
	}
	rootDir := testCfg.Library
	err := os.MkdirAll(rootDir, 0775)
	if err != nil {
		t.Fatal(err)
	}
	lib, err := GetLib(testCfg.Library)
	if err != nil {
		t.Fatal(err)
	}
	setus := make([]string, 0)

	f1, b1 := createTemp("I am a file")
	f2, _ := createTemp("I am anorther file")
	t.Run("Add-GetFile", func(t *testing.T) {
		gotSum, err := lib.Add(f1)
		if err != nil {
			t.Fatal(err)
		}
		_, err = lib.Add(f2)
		if err != nil {
			t.Fatal(err)
		}
		for k := range lib.Setus {
			setus = append(setus, k)
		}
		if len(setus) != 2 {
			t.Fatalf("want %d, got %d", 2, len(setus))
		}
		setu, err := lib.GetFile(gotSum)
		if err != nil {
			t.Fatal(err)
		}
		defer setu.Close()
		b, err := io.ReadAll(setu)
		if err != nil {
			t.Fatal(err)
		}
		if string(b1) != string(b) {
			t.Fatalf("Failed to verify temp data\nwant: %v\ngot: %v", b1, b)
		}
		wantSum := fmt.Sprintf("%x", md5.Sum(b1))
		if wantSum != gotSum {
			t.Fatalf("Diffrent md5 checksum \nwant: %v\ngot: %v", b1, b)
		}
	})
	t.Run("Del", func(t *testing.T) {
		ok := lib.Rm(setus[0])
		if !ok {
			t.Fatal("there should be true, but itn't")
		}
		ok = lib.Rm(setus[0])
		if ok {
			t.Fatal("there should be false, but itn't")
		}
		_, err := lib.GetFile(setus[0])
		if err == nil {
			t.Fatal("there should be an error, but itn't")
		}
	})
	t.Run("CreateLib", func(t *testing.T) {
		create := func(name string) {
			err := lib.CreateLib(name)
			if err != nil {
				t.Fatal(err)
			}
		}
		create("lib1")
		create("lib2")
		create("lib3")
		lib1 := lib.SubLib["lib1"]
		if lib1 == nil {
			t.Fatal("failed to create a new library")
		}
		want := "lib1"
		if lib1.Name != want {
			t.Fatalf("want %s, got %s", want, lib1.Name)
		}
	})
	t.Run("RmLib", func(t *testing.T) {
		sb := lib.SubLib["lib1"]
		ok := lib.RmLib("lib1")
		if !ok {
			t.Fatal("there should be true, but itn't")
		}
		ok = lib.RmLib("lib1")
		if ok {
			t.Fatal("there should be false, but itn't")
		}
		if !util.IsExist(sb.Dir) {
			t.Fatalf("dont delete lib's dir: %s", sb.Dir)
		}
	})
}
