package repository

import (
	"fmt"
	"github.com/jackc/pgx"
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/user"
)

type pgRepository struct {
	db *pgx.ConnPool
}

func NewPgRepository(db *pgx.ConnPool) user.IRepositoryUser {
	return &pgRepository{db: db}
}

func (p pgRepository) CreateUser(userNew models.User) errors.Err {
	_, err := p.db.Exec("insert into users (nickname, fullname, about, email) values ($1, $2, $3, $4)",
		userNew.Nickname, userNew.FullName, userNew.About, userNew.Email)
	if err != nil {
		msg := err.Error()
		if errors.RecordExists(msg) {
			return errors.RespErr{StatusCode: errors.ConflictCode, Message: errors.UserAlreadyExists}
		}
		return errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	}
	return nil
}

func (p pgRepository) GetUser(nickname string) (models.User, errors.Err) {
	var buf models.User
	err := p.db.QueryRow("select * from users where nickname = $1", nickname).
		Scan(&buf.Nickname, &buf.FullName, &buf.About, &buf.Email)
	if err != nil {
		msg := err.Error()
		if errors.EmptyResult(msg) {
			return models.User{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.UserNotFoundMsg}
		}
		return models.User{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	}
	return buf, nil
}

func (p pgRepository) UpdateUser(userNew models.User) (models.User, errors.Err) {
	var userOld models.User
	userOld, err := p.GetUser(userNew.Nickname)
	if err != nil {
		return models.User{}, err
	}

	sql := "update users set "
	needUpdate := false
	if userNew.Email != "" && userNew.Email != userOld.Email {
		sql = fmt.Sprintf("%s email = '%s' ", sql, userNew.Email)
		needUpdate = true
	}
	if userNew.FullName != "" && userNew.FullName != userOld.FullName {
		if needUpdate {
			sql += ","
		}
		sql = fmt.Sprintf("%s fullname = '%s' ", sql, userNew.FullName)
		needUpdate = true
	}
	if userNew.About != "" && userNew.About != userOld.About {
		if needUpdate {
			sql += ","
		}
		sql = fmt.Sprintf("%s about = '%s' ", sql, userNew.About)
	}
	if !needUpdate {
		return userOld, nil
	}
	if _, err := p.db.Exec(sql); err != nil {
		msg := err.Error()
		if errors.RecordExists(msg) {
			return models.User{}, errors.RespErr{StatusCode: errors.ConflictCode, Message: errors.UserAlreadyExists}
		} else if errors.UserNotFound(msg) {
			return models.User{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.UserNotFoundMsg}
		}
		return models.User{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	}
	return p.GetUser(userNew.Nickname)
}
