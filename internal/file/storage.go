package file

import (
	"errors"
	"time"

	"github.com/akasrt/filensy/internal/database"
	"github.com/akasrt/filensy/internal/util/errorx"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	Create(data FileData) (FileData, error)
	Get(code string) (FileData, error)
	Delete(code string) error
	IsDuplicateErr(err error) bool
}

func NewStorage() Storage {
	return &storage{
		db: database.GetDB(),
	}
}

type storage struct {
	db *sqlx.DB
}

func (s *storage) Create(rqData FileData) (FileData, error) {
	query := createFileQuery()
	_, err := s.db.NamedExec(query, rqData)
	if err != nil {
		return FileData{}, errorx.WrapMysqlError(err)
	}

	data, err := s.Get(rqData.Code)
	if err != nil {
		return FileData{}, err
	}
	return data, nil
}

func (s *storage) Get(code string) (FileData, error) {
	var data FileData
	query := getFileQuery()
	err := s.db.Get(&data, query, code, time.Now())
	if err != nil {
		return FileData{}, errorx.WrapMysqlError(err)
	}

	return data, nil
}

func (s *storage) Delete(code string) error {
	query := deleteFileQuery()
	_, err := s.db.Exec(query, code)
	if err != nil {
		return errorx.WrapMysqlError(err)
	}

	return nil
}

func (s *storage) IsDuplicateErr(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}
