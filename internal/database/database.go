package database

type Handle interface {
	GetRow () (map[string]string, error)
	GetList () (map[int]map[string]string, error)
	Insert () (int64, error)
	InsertAll () (int64, error)
	Update () (int64, error)
	Query () error
	BeginTransaction() error
	Commit() error
	Rollback() error
}


const (
	MYSQL = "mysql"
	H_GET_ROW = "getrow"
	H_GET_LIST = "getlist"
	H_INSERT = "insert"
	H_INSERT_ALL = "insertall"
	H_UPDATE = "update"
	H_QUERY = "query"
	H_BEGIN_TRAN = "begintran"
	H_COMMIT = "commit"
	H_ROLLBACK = "rollback"

)

