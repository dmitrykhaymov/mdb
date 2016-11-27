package dal

import (
	"github.com/Bnei-Baruch/mdb/rest"

	"database/sql"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func Init() bool {
	url := viper.GetString("mdb.url")
	db, err := sql.Open("postgres", url)
	checkErr(err)
	defer db.Close()
	return true
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type DalError struct {
    err string
}

func (e DalError) Error() string {
    return e.err
}

func CaptureStart(rest.CaptureStart) (bool, error) {

    // Implementation goes here...

    return false, DalError{err: "not implemented"}
}
