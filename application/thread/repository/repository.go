package repository

import (
	"github.com/jackc/pgx"
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/forum"
)

type pgRepository struct {
	db *pgx.ConnPool
}

func NewPgRepository(db *pgx.ConnPool) forum.IRepositoryForum {
	return &pgRepository{db: db}
}

func (p pgRepository) CreateForum(forumNew models.Forum) errors.Err {
	_, err := p.db.Exec("insert into forums (title, \"user\", slug) values ($1, $2, $3)",
		forumNew.Title, forumNew.User, forumNew.Slug)
	if err != nil {
		msg := err.Error()
		if errors.UserNotFound(msg) {
			return errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.UserNotFoundMsg}
		} else if errors.RecordExists(msg) {
			return errors.RespErr{StatusCode: errors.ConflictCode}
		}
		return errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(err.Error())}
	}
	return nil
}

func (p pgRepository) GetBySlug(slug string) (models.Forum, errors.Err) {
	var buf models.Forum
	err := p.db.QueryRow("select f.Title, f.User, f.Posts, f.Threads from forums as f where slug = $1", slug).
		Scan(&buf.Title, &buf.User, &buf.Posts, &buf.Threads)
	if err != nil {
		msg := err.Error()
		if errors.EmptyResult(msg) {
			return models.Forum{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.ForumNotFoundMsg}
		}
		return models.Forum{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	}
	buf.Slug = slug
	return buf, nil
}

func (p pgRepository) GetUsers(uri models.QueryParams) (models.UsersList, errors.Err) {
	//sql := "select * from users_on_forum uf where uf"
	return models.UsersList{Users: []models.User{models.User{Nickname: "aaa", About: "a"}}}, nil
}
