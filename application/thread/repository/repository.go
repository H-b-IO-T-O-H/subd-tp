package repository

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"strconv"
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/common/utils"
	"subd/application/thread"
	"time"
)

type pgRepository struct {
	db *pgx.ConnPool
}

func NewPgRepository(db *pgx.ConnPool) thread.IRepositoryThread {
	return &pgRepository{db: db}
}

func (p pgRepository) CreateThread(thread models.Thread) (models.Thread, errors.Err) {
	var err error
	if thread.Created.IsZero() {
		thread.Created = time.Now()
	}
	isForumExist := true
	if err := p.db.QueryRow("select exists(select 1 from forums where slug=$1)", thread.Forum).Scan(&isForumExist);
		err != nil || !isForumExist {
		return models.Thread{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.ForumNotFoundMsg}
	}
	query := ""
	if thread.Slug == "" {
		query = fmt.Sprintf(`insert into threads (title, author, forum, message, created)
			values('%s','%s','%s','%s','%s') returning id, forum`,
			thread.Title, thread.Author, thread.Forum, thread.Message, thread.Created.Format(time.RFC3339Nano))
	} else {
		query = fmt.Sprintf(`insert into threads (title, author, forum, message, slug, created)
			values('%s','%s',(select slug as forum from forums where slug = '%s'),'%s','%s','%s') returning id, forum`,
			thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created.Format(time.RFC3339Nano))
	}

	err = p.db.QueryRow(query).Scan(&thread.ID, &thread.Forum)
	if err != nil {
		msg := err.Error()
		if errors.RecordExists(msg) {
			return models.Thread{}, errors.RespErr{StatusCode: errors.ConflictCode}
		} else if errors.UserNotFound(msg) {
			return models.Thread{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.UserNotFoundMsg}
		} else if errors.ForumNotFound(msg) {
			return models.Thread{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.ForumNotFoundMsg}
		}
		return models.Thread{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	}
	return thread, nil
}

func (p pgRepository) GetBySlug(slug string) (models.Thread, errors.Err) {
	buf := models.Thread{}
	nullSlug := &pgtype.Varchar{}
	err := p.db.QueryRow("select t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created from threads as t where t.slug = $1", slug).
		Scan(&buf.ID, &buf.Title, &buf.Author, &buf.Forum, &buf.Message, &buf.Votes, nullSlug, &buf.Created)
	if err != nil {
		msg := err.Error()
		if errors.EmptyResult(msg) {
			return models.Thread{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.UserNotFoundMsg}
		}
		return models.Thread{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	}
	buf.Slug = nullSlug.String
	//_ = p.db.QueryRow("select f.slug from forum f where f.slug = $1", buf.Forum).Scan(&buf.Forum)
	return buf, nil
}

func (p pgRepository) GetById(id int) (models.Thread, errors.Err) {
	buf := models.Thread{ID: id}
	nullSlug := &pgtype.Varchar{}
	err := p.db.QueryRow("select t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created from threads as t where t.id = $1", id).
		Scan(&buf.Title, &buf.Author, &buf.Forum, &buf.Message, &buf.Votes, nullSlug, &buf.Created)
	if err != nil {
		msg := err.Error()
		if errors.EmptyResult(msg) {
			return models.Thread{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.UserNotFoundMsg}
		}
		return models.Thread{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	}
	buf.Slug = nullSlug.String
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
	if _, err := p.db.Exec("insert into votes (nick, voice, thread_id) values ($1, $2, $3) on conflict (nick, thread_id) do update set voice = $2 where votes.voice <> $2", vote.Nick, voice, id);
		err != nil {
		msg := err.Error()
		if errors.UserNotFound(msg) {
			return models.Thread{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.ThreadNotFoundMsg}
		}
		return models.Thread{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	}
	return p.GetById(id)
}

func (p pgRepository) GetThreadPosts(params models.QueryPostParams) (models.PostsList, errors.Err) {
	var (
		threadOld models.Thread
		err       errors.Err
	)
	sortedPosts := models.PostsList{}
	if id, _ := strconv.Atoi(params.SlugId); id != 0 {
		threadOld, err = p.GetById(id)
	} else {
		threadOld, err = p.GetBySlug(params.SlugId)
	}
	if err != nil {
		return models.PostsList{}, err
	}
	sql := prepareQueryForSort(threadOld.ID, params)
	res, err1 := p.db.Query(sql)
	if err1 != nil {
		return models.PostsList{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(err1.Error())}
	}
	defer res.Close()

	buf := models.Post{}
	for res.Next() {
		if err := res.Scan(&buf.ID, &buf.Parent, &buf.Author, &buf.Message, &buf.IsEdited, &buf.Forum, &buf.Thread, &buf.Created); err != nil {
			return models.PostsList{}, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(err.Error())}
		}
		sortedPosts = append(sortedPosts, buf)
	}
	return sortedPosts, nil
}

func prepareQueryForSort(threadId int, params models.QueryPostParams) string {

	sql := fmt.Sprintf("select id, parent, author, message, isedited, forum, thread, created from posts where thread = '%d'", threadId)
	//tree
	if params.Sort == utils.Tree {
		if params.Since != 0 {
			if params.Desc {
				sql = fmt.Sprintf("%s and path < (select path from posts where id = %d) ", sql, params.Since)
			} else {
				sql = fmt.Sprintf("%s and path > (select path from posts where id = %d) ", sql, params.Since)
			}
		}

		if params.Desc {
			sql = fmt.Sprintf("%s order by path desc, id desc ", sql)
		} else {
			sql = fmt.Sprintf("%s order by path, id ", sql)
		}
		if params.Limit != 0 {
			sql = fmt.Sprintf("%s limit %d", sql, params.Limit)
		}
	}
	//parent_tree
	if params.Sort == utils.ParentTree {
		sql = fmt.Sprintf("(select id from posts where thread = %d and parent = 0 ", threadId)

		if params.Since != 0 {
			if params.Desc {
				sql = fmt.Sprintf("%s and path[1] < (select path[1] from posts where id = %d) ", sql, params.Since)
			} else {
				sql = fmt.Sprintf("%s and path[1] > (select path[1] from posts where id = %d) ", sql, params.Since)
			}
		}
		if params.Desc {
			sql = fmt.Sprintf("%s order by id desc ", sql)
		} else {
			sql = fmt.Sprintf("%s order by id ", sql)
		}
		if params.Limit != 0 {
			sql = fmt.Sprintf("%s limit %d", sql, params.Limit)
		}
		sql += ")"
		sql = fmt.Sprintf(`select id, parent, author, message, isedited, forum, thread, created from posts where path[1] in %s `, sql)
		if params.Desc {
			sql = fmt.Sprintf("%s order by path[1] desc, path, id ", sql)
		} else {
			sql = fmt.Sprintf("%s order by path ", sql)
		}
	}
	//flat
	if params.Sort == utils.Flat {
		if params.Since != 0 {
			if params.Desc {
				sql = fmt.Sprintf("%s and id < '%d'", sql, params.Since)
			} else {
				sql = fmt.Sprintf("%s and id > '%d'", sql, params.Since)
			}
		}
		if params.Desc {
			sql = fmt.Sprintf("%s order by id desc ", sql)
		} else {
			sql = fmt.Sprintf("%s order by id ", sql)
		}
		if params.Limit != 0 {
			sql = fmt.Sprintf("%s limit %d", sql, params.Limit)
		}
	}

	return sql
}
