package context

import (
	"github.com/pkg/errors"
	"socket/internal/database"
)

type Database struct {
	mode      string
	drives    string
	params 	  map[string]interface{}
	processor database.Handle
}

var dbContainer = getDBContainer()

func getDBContainer () map[string]func(data map[string]interface{}) (database.Handle, error){

	db := make(map[string]func(data map[string]interface{}) (database.Handle, error))
	db[database.MYSQL] = func(data map[string]interface{}) (database.Handle, error) {
		return database.NewMysql(data)
	}

	return db
}

func (db *Database) Execute () (int, string, map[interface{}]interface{}) {

	handle, err := db.getHandle()
	if err != nil {
		return -1, err.Error(), make(map[interface{}]interface{})
	}
	data := make(map[interface{}]interface{})
	switch db.mode {
	case database.H_INSERT:
		data["insertId"], err = handle.Insert()
	}

	return 1000, "success", data
}

func (db *Database) getHandle () (database.Handle, error) {

	if handle, isOk := dbContainer[db.drives]; isOk {
		return handle(db.params)
	}
	return nil, errors.New("not defined database handle")

}

func NewStore (data map[string]interface{}) (Processor, error){
	s := &Database{}
	var isOk bool
	s.params = data
	if s.mode, isOk = data["mode"].(string); !isOk{
		return nil, errors.New("not defined mode")
	}

	if s.drives, isOk = data["drives"].(string); !isOk{
		return nil, errors.New("not defined drives")
	}

	var err error
	s.processor, err = database.NewMysql(data)
	return s, err
}

