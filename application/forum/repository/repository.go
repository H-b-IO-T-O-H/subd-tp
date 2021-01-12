package repository

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
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

func (p pgRepository) CreateForum(forumNew models.Forum) (models.Forum, errors.Err) {
	_, err := p.db.Exec("insert into forums (title, \"user\", slug) values ($1, $2, $3)",
		forumNew.Title, forumNew.User, forumNew.Slug)
	if err != nil {
		msg := err.Error()
		if errors.UserNotFound(msg) {
			return models.Forum{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.UserNotFoundMsg}
		} else if errors.RecordExists(msg) {
			return models.Forum{}, errors.RespErr{StatusCode: errors.ConflictCode}
		}
		return models.Forum{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(err.Error())}
	}
	if err := p.db.QueryRow(`select u.nickname from users u where u.nickname = $1`, forumNew.User).Scan(&forumNew.User); err != nil {
		return models.Forum{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(err.Error())}
	}
	return forumNew, nil
}

func (p pgRepository) GetBySlug(slug string) (models.Forum, errors.Err) {
	var buf models.Forum
	err := p.db.QueryRow("select * from forums where slug = $1", slug).
		Scan(&buf.Title, &buf.User, &buf.Slug, &buf.Posts, &buf.Threads)
	if err != nil {
		msg := err.Error()
		if errors.EmptyResult(msg) {
			return models.Forum{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.ForumNotFoundMsg}
		}
		return models.Forum{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	}
	if err := p.db.QueryRow(`select u.nickname from users u where u.nickname = $1`, buf.User).Scan(&buf.User); err != nil {
		return models.Forum{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(err.Error())}
	}
	return buf, nil
}

func (p pgRepository) GetUsers(query models.QueryParams) (models.UsersList, errors.Err) {
	isForumExist := true
	err := p.db.QueryRow("select exists(select 1 from forums where slug=$1)", query.Slug).Scan(&isForumExist)
	if err != nil || !isForumExist {
		return models.UsersList{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.ForumNotFoundMsg}
	}

	sql := "select uf.nick, u.fullname, u.about, u.email from users_on_forum uf join users u on uf.nick=u.nickname where slug="
	sql = fmt.Sprintf("%s '%s' ", sql, query.Slug)
	if query.Since != "" {
		if !query.Desc {
			sql = fmt.Sprintf("%s and u.nickname > '%s' ", sql, query.Since)
		} else {
			sql = fmt.Sprintf("%s and u.nickname < '%s' ", sql, query.Since)
		}
	}
	if query.Desc {
		sql = fmt.Sprintf("%s order by nickname DESC", sql)
	} else {
		sql = fmt.Sprintf("%s order by nickname ASC", sql)
	}
	if query.Limit > 0 {
		sql = fmt.Sprintf("%s limit %d", sql, query.Limit)
	}
	rows, err := p.db.Query(sql)
	if err != nil {
		return models.UsersList{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: errors.ServerErrorMsg}
	}
	defer rows.Close()

	usersList := models.UsersList{}
	one := models.User{}
	for rows.Next() {
		if err := rows.Scan(&one.Nickname, &one.FullName, &one.About, &one.Email); err != nil {
			return models.UsersList{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: errors.ServerErrorMsg}
		}
		usersList = append(usersList, one)
	}
	return usersList, nil
}

func (p pgRepository) GetThreads(query models.QueryParams) (models.ThreadsList, errors.Err) {
	isForumExist := true
	err := p.db.QueryRow("select exists(select 1 from forums where slug=$1)", query.Slug).Scan(&isForumExist)
	if err != nil || !isForumExist {
		return models.ThreadsList{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.ForumNotFoundMsg}
	}
	sql := "select * from threads t where forum ="
	sql = fmt.Sprintf("%s '%s' ", sql, query.Slug)
	if query.Since != "" {
		if !query.Desc {
			sql = fmt.Sprintf("%s and t.created >= '%s' ", sql, query.Since)
		} else {
			sql = fmt.Sprintf("%s and t.created <= '%s' ", sql, query.Since)
		}
	}
	if query.Desc {
		sql = fmt.Sprintf("%s order by created DESC", sql)
	} else {
		sql = fmt.Sprintf("%s order by created ASC", sql)
	}
	if query.Limit > 0 {
		sql = fmt.Sprintf("%s limit %d", sql, query.Limit)
	}
	rows, err := p.db.Query(sql)
	if err != nil {
		return models.ThreadsList{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: errors.ServerErrorMsg}
	}
	defer rows.Close()

	threadsList := models.ThreadsList{}
	one := models.Thread{}
	nullSlug := &pgtype.Varchar{}
	for rows.Next() {
		if err := rows.Scan(&one.ID, &one.Title, &one.Author, &one.Forum, &one.Message, &one.Votes, nullSlug, &one.Created); err != nil {
			return models.ThreadsList{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: errors.ServerErrorMsg}
		}
		one.Slug = nullSlug.String
		threadsList = append(threadsList, one)
	}
	return threadsList, nil
}
