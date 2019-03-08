package database

import (
	"github.com/jmoiron/sqlx"
)

type Mysql struct {
	force bool
	isTransaction bool
	sql string
	params map[int]map[string]interface{}
	*sqlx.Tx
}

func NewMysql (data map[string]interface{}) (Handle, error) {
	mysql := &Mysql{}
	mysql.force = data["force"].(bool)
	mysql.sql = data["sql"].(string)
	mysql.isTransaction = data["isTransaction"].(bool)
	mysql.params = make(map[int]map[string]interface{})

	for idx, item := range data["params"].([]interface{}) {
		mysql.params[idx] = item.(map[string]interface{})
	}
	return mysql, nil
}


func (m *Mysql) GetRow () (map[string]string, error) {
	return nil, nil
}

func (m *Mysql) GetList () (map[int]map[string]string, error){
	return nil, nil
}

func (m *Mysql) Insert () (int64, error) {
	return 100, nil
}

func (m *Mysql) InsertAll () (int64, error) {
	return 0, nil
}

func (m *Mysql) Update () (int64, error) {
	return 0, nil
}

func (m *Mysql) Query () error {
	return nil
}

func (m *Mysql) BeginTransaction() error {
	return  nil
}

func (m *Mysql) Commit() error {
	return nil
}

func (m *Mysql) Rollback() error {
	return nil
}
