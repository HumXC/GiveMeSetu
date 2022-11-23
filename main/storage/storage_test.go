package storage

import (
	"give-me-setu/main/conf"
	"os"
	"path"
	"testing"
)

func clear() {
	os.RemoveAll(dataDir)
	os.Mkdir(dataDir, 0775)
}
func TestMain(m *testing.M) {
	clear()
	defer clear()
}

var dataDir string = "../../test/storage"
var testCfg conf.Config = conf.Config{
	DataDir: dataDir,
	Library: path.Join(dataDir, "root-library"),
	Database: conf.Database{
		Driver: "sqldb",
		Name:   "test-db",
	},
}
