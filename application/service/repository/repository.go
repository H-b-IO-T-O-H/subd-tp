package repository

import (
	"github.com/jackc/pgx"
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/service"
)

type pgRepository struct {
	db *pgx.ConnPool
}

func NewPgRepository(db *pgx.ConnPool) service.IRepositoryService {
	return &pgRepository{db: db}
}

func (p pgRepository) GetStatus() (models.ServiceStatus, errors.Err) {
	var buf models.ServiceStatus
	err := p.db.QueryRow("select * from (select count(*) from users) as u "+
		"cross join (select count(*) from forums) as f "+
		"cross join (select count(*) from threads) as t "+
		"cross join (select count(*) from posts) as p").Scan(&buf.UsersCnt, &buf.ForumsCnt, &buf.ThreadsCnt, &buf.PostsCnt)
	if err != nil {
		return models.ServiceStatus{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: errors.ServerErrorMsg}
	}
	return buf, nil
}

func (p pgRepository) Clear() errors.Err {
	_, err := p.db.Exec("truncate forums, users, threads, posts, votes, users_on_forum cascade")
	if err != nil {
		return errors.RespErr{StatusCode: errors.ServerErrorCode, Message: errors.ServerErrorMsg}
	}
	return nil
}
