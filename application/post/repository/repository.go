package repository

import (
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	"strconv"
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/forum"
	"subd/application/post"
	"subd/application/thread"
	"subd/application/user"
	"time"
)

type pgRepository struct {
	db          *pgx.ConnPool
	reposUser   user.IRepositoryUser
	reposForum  forum.IRepositoryForum
	reposThread thread.IRepositoryThread
}

func NewPgRepository(db *pgx.ConnPool, repUser user.IRepositoryUser,
	repForum forum.IRepositoryForum, repThread thread.IRepositoryThread) post.IRepositoryPost {
	return &pgRepository{db: db,
		reposUser:   repUser,
		reposForum:  repForum,
		reposThread: repThread}
}

func (p pgRepository) GetById(id int64) (models.Post, errors.Err) {
	buf := models.Post{ID: id}

	if err := p.db.QueryRow(`select p.parent, p.author, p.message, 
		p.isedited,p.forum, p.thread, p.created from posts p where id = $1`, id).Scan(&buf.Parent,
		&buf.Author, &buf.Message, &buf.IsEdited, &buf.Forum, &buf.Thread, &buf.Created); err != nil {
		return models.Post{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.PostNotFoundMsg}
	}
	return buf, nil
}

func (p pgRepository) CreatePost(posts models.PostsList, slugId string) (*models.PostsList, errors.Err) {
	var (
		thread models.Thread
		err    errors.Err
		buffer string
	)
	size := len(posts)

	if id, _ := strconv.Atoi(slugId); id != 0 {
		thread, err = p.reposThread.GetById(id)
	} else {
		thread, err = p.reposThread.GetBySlug(slugId)
	}
	if err != nil {
		return nil, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.ThreadNotFoundMsg}
	}
	if size == 0 {
		return nil, nil
	}
	parents := make(map[int64]int)
	id := 0
	for _, single := range posts {
		if _, exist := parents[single.Parent]; !exist && single.Parent != 0 {
			if err := p.db.QueryRow(`select p.thread from posts p where p.id = $1`, single.Parent).Scan(&id); err != nil || id != thread.ID {
				return nil, errors.RespErr{StatusCode: errors.ConflictCode, Message: errors.PostNotFoundMsg}
			}
			parents[single.Parent] = thread.ID
		}
	}
	var values string

	for i := range posts {
		values += fmt.Sprintf("('%s')", posts[i].Author)
		if i < size-1 {
			values += ", "
		}
	}

	query := fmt.Sprintf(`with ctetable(nick) as (values %s)
		select nickname from users inner join ctetable on ctetable.nick = users.nickname `, values)

	usersNicks := make([]string, size)

	res, err1 := p.db.Query(query)
	if err1 != nil {
		return nil, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(err1.Error())}
	}
	defer res.Close()
	i := 0
	for res.Next() {
		_ = res.Scan(&usersNicks[i])
		i++
	}
	if len(usersNicks) != size {
		return nil, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.UserNotFoundMsg}
	}

	buffer = ""
	for i = range usersNicks {
		buffer += fmt.Sprintf("('%s', '%s')", thread.Forum, usersNicks[i])
		if i < size-1 {
			buffer += ", "
		}
	}
	query = fmt.Sprintf(`insert into users_on_forum values %s on conflict do nothing;`, buffer)
	if _, err := p.db.Exec(query); err != nil {
		msg := err.Error()
		if errors.UserNotFound(msg) {
			return nil, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.UserNotFoundMsg}
		}
		return nil, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	}

	timeCreate := strfmt.DateTime(time.Now())
	buffer = ""
	for i = range posts {
		posts[i].Forum = thread.Forum
		posts[i].Thread = thread.ID
		posts[i].Created = timeCreate
		buffer += fmt.Sprintf("('%d', '%s', '%s', '%t', '%s', '%d', '%s')",
			posts[i].Parent, posts[i].Author, posts[i].Message, posts[i].IsEdited, posts[i].Forum,
			posts[i].Thread, timeCreate)
		if i < size-1 {
			buffer += ", "
		}
	}

	sql := fmt.Sprintf(`insert into posts (parent, author, message, isedited, forum, 
			thread, created) values %s returning id`, buffer)
	if res2, err2 := p.db.Query(sql); err2 != nil {
		msg := err2.Error()
		if errors.UserNotFound(msg) {
			return nil, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.UserNotFoundMsg}
		}
		return nil, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(msg)}
	} else {
		i = 0
		for res2.Next() {
			_ = res2.Scan(&posts[i].ID)
			i++
		}
	}
	if _, err := p.db.Exec("update forums set posts = posts + $1 where slug = $2", size, thread.Forum); err != nil {
		return nil, errors.RespErr{StatusCode: errors.ServerErrorCode, Message: []byte(err.Error())}
	}

	return &posts, nil
}

func (p pgRepository) UpdatePost(postUpdate models.PostUpdate) (models.Post, errors.Err) {
	var postOld models.Post
	if postUpdate.Message == "" {
		return p.GetById(postUpdate.ID)
	}
	query := fmt.Sprintf(`update posts 
		set message = '%s',
		isedited = posts.isedited or posts.message != '%s'
		where posts.id = '%d'
		returning id, parent, author, isedited, forum, thread, created`, postUpdate.Message, postUpdate.Message, postUpdate.ID)
	err := p.db.QueryRow(query).Scan(&postOld.ID, &postOld.Parent, &postOld.Author,
		&postOld.IsEdited, &postOld.Forum, &postOld.Thread, &postOld.Created)
	if err != nil || postOld.ID == 0 {
		return models.Post{}, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.PostNotFoundMsg}
	}
	postOld.Message = postUpdate.Message
	return postOld, nil
}

func (p pgRepository) GetFull(params models.PostGetParams) (models.PostFull, errors.Err) {
	postFull := models.PostFull{}
	bufPost := models.Post{ID: params.PostId}

	if err := p.db.QueryRow(`select p.parent, p.author, p.message, 
		p.isedited,p.forum, p.thread, p.created from posts p where id = $1`, params.PostId).Scan(&bufPost.Parent,
		&bufPost.Author, &bufPost.Message, &bufPost.IsEdited, &bufPost.Forum, &bufPost.Thread, &bufPost.Created); err != nil {
		return postFull, errors.RespErr{StatusCode: errors.NotFoundCode, Message: errors.PostNotFoundMsg}
	}
	postFull.Post = &bufPost
	if params.HaveUser {
		bufUser, err := p.reposUser.GetUser(bufPost.Author)
		if err != nil {
			return postFull, err
		}
		postFull.Author = &bufUser
	}
	if params.HaveForum {
		bufForum, err := p.reposForum.GetBySlug(bufPost.Forum)
		if err != nil {
			return postFull, err
		}
		postFull.Forum = &bufForum
	}
	if params.HaveThread {
		bufThread, err := p.reposThread.GetById(bufPost.Thread)
		if err != nil {
			return postFull, err
		}
		postFull.Thread = &bufThread
	}

	return postFull, nil
}
