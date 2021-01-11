package repository

import (
	"fmt"
	"github.com/jackc/pgx"
	"strconv"
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/thread"
)

type pgRepository struct {
	db *pgx.ConnPool
}

func NewPgRepository(db *pgx.ConnPool) thread.IRepositoryThread {
	return &pgRepository{db: db}
}

func (p pgRepository) CreateThread(thread models.Thread) (int, errors.Err) {
	var err error
	//TODO:transactions
	//if thread.Slug == "" {
	//	err = p.db.QueryRow("insert into threads (title, author, forum, message, created) values ($1, $2, $3, $4, $5) returning id",
	//		thread.Title, thread.Author, thread.Forum, thread.Message, thread.Created).Scan(&thread.ID)
	//} else {
	//	err = p.db.QueryRow("insert into threads (title, author, forum, message, slug, created) values ($1, $2, $3, $4, $5, $6) returning id",
	//		thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created).Scan(&thread.ID)
	//}
	//tx, _ := p.db.Begin()
	//defer tx.Rollback()
	err = p.db.QueryRow("insert into threads (title, author, forum, message, slug, created) values ($1, $2, $3, $4, $5, $6) returning id",
		thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created).Scan(&thread.ID)
	if err != nil {
		msg := err.Error()
		if errors.RecordExists(msg) {
			return -1, errors.RespErr{StatusCode: errors.ConflictCode}
		} else if errors.UserNotFound(msg) {
			return -1, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.UserNotFoundMsg}
		}
		return -1, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	}
	//err = tx.Commit()
	//if err != nil {
	//	return -1, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: errors.ServerErrorMsg}
	//}
	return thread.ID, nil
}

func (p pgRepository) GetBySlug(slug string) (models.Thread, errors.Err) {
	buf := models.Thread{Slug: slug}
	err := p.db.QueryRow("select t.id, t.title, t.author, t.forum, t.message, t.votes, t.created from threads as t where t.slug = $1", slug).
		Scan(&buf.ID, &buf.Title, &buf.Author, &buf.Forum, &buf.Message, &buf.Votes, &buf.Created)
	if err != nil {
		msg := err.Error()
		if errors.EmptyResult(msg) {
			return models.Thread{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.UserNotFoundMsg}
		}
		return models.Thread{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	}
	return buf, nil
}

func (p pgRepository) GetById(id int) (models.Thread, errors.Err) {
	buf := models.Thread{ID: id}
	err := p.db.QueryRow("select t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created from threads as t where t.id = $1", id).
		Scan(&buf.Title, &buf.Author, &buf.Forum, &buf.Message, &buf.Votes, &buf.Slug, &buf.Created)
	if err != nil {
		msg := err.Error()
		if errors.EmptyResult(msg) {
			return models.Thread{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.UserNotFoundMsg}
		}
		return models.Thread{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	}
	return buf, nil
}

func (p pgRepository) UpdateBySlugOrId(threadNew models.Thread, slugId string) (models.Thread, errors.Err) {
	var (
		threadOld models.Thread
		err       errors.Err
	)
	if id, _ := strconv.Atoi(slugId); id != 0 {
		threadOld, err = p.GetById(id)
	} else {
		threadOld, err = p.GetBySlug(slugId)
	}
	if err != nil {
		return models.Thread{}, err
	}

	sql := "update threads set "
	needUpdate := false
	if threadNew.Title != "" && threadNew.Title != threadOld.Title {
		sql = fmt.Sprintf("%s title = '%s' ", sql, threadNew.Title)
		needUpdate = true
		threadOld.Title = threadNew.Title
	}
	if threadNew.Message != "" && threadNew.Message != threadOld.Message {
		if needUpdate {
			sql += ","
		}
		sql = fmt.Sprintf("%s message = '%s' ", sql, threadNew.Message)
		threadOld.Message = threadNew.Message
		needUpdate = true
	}
	if !needUpdate {
		return threadOld, nil
	}
	if _, err := p.db.Exec(sql); err != nil {
		msg := err.Error()
		if errors.UserNotFound(msg) {
			return models.Thread{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.UserNotFoundMsg}
		}
		return models.Thread{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	}
	return threadOld, nil
}

func (p pgRepository) UpsertVote(vote models.Vote) (models.Thread, errors.Err) {
	id := 0
	voice := true
	if vote.Voice == -1 {
		voice = false
	}
	if id, _ = strconv.Atoi(vote.SlugOrId); id != 0 {
		isThreadExist := true
		if err := p.db.QueryRow("select exists(select 1 from threads where id=$1)", id).Scan(&isThreadExist);
			err != nil || !isThreadExist {
			return models.Thread{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.ThreadNotFoundMsg}
		}
	} else {
		if err := p.db.QueryRow("select t.id from threads t where t.slug = $1", vote.SlugOrId).Scan(&id); err != nil {
			return models.Thread{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.ThreadNotFoundMsg}
		}
	}
	if _, err := p.db.Exec("insert into votes (nick, voice, thread_id) values ($1, $2, $3) on conflict (nick, thread_id) do update set voice = $2", vote.Nick, voice, id);
		err != nil {
		msg := err.Error()
		if errors.UserNotFound(msg) {
			return models.Thread{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.ThreadNotFoundMsg}
		}
		return models.Thread{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	}
	return p.GetById(id)
}
