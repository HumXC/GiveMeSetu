package storage

import (
	"give-me-setu/main/conf"
	"path"
)

var dataDir string = "../../test/storage"
var testCfg conf.Config = conf.Config{
	DataDir: dataDir,
	Library: path.Join(dataDir, "root-library"),
	Database: conf.Database{
		Driver: "sqldb",
		Name:   "test-db",
	},
}
