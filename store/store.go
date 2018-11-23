/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package store

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/seeleteam/go-seele/database"

	"github.com/seeleteam/go-seele/database/leveldb"
)

// ...
const (
	DbName   = "Seele-e2e"
	BuildKey = "Seele-build"
	CoverKey = "Seele-cover"
	BenchKey = "Seele-bench"
)

// DB ...
var db database.Database

func init() {
	var err error
	if db, err = prepareDB(DbName); err != nil {
		fmt.Println("create db err:", err)
		os.Exit(1)
	}
}

func prepareDB(dbName string) (database.Database, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(usr.HomeDir, dbName)
	return leveldb.NewLevelDB(dbPath)

}

// Save the e2e test result
func Save(date, benchresult string, coverbyte []byte) {
	db.Put([]byte(date+CoverKey), coverbyte)
	db.Put([]byte(date+BenchKey), []byte(benchresult))
}

// Get the e2e test result
func Get(date string) (benchresult string, coverbyte []byte) {
	coverbyte, err := db.Get([]byte(date + CoverKey))
	if err != nil {
		fmt.Println("get cover result err:", err)
		return
	}

	benchbyte, err := db.Get([]byte(date + BenchKey))
	if err != nil {
		fmt.Println("get bench result err:", err)
		return
	}

	return string(benchbyte), coverbyte
}
