package thread

import (
	"subd/application/common/errors"
	"subd/application/common/models"
)

type IRepositoryThread interface {
	CreateThread(thread models.Thread) (models.Thread, errors.Err)
	GetBySlug(slug string) (models.Thread, errors.Err)
	GetById(id int) (models.Thread, errors.Err)
	UpdateBySlugOrId(threadNew models.Thread, slugId string) (models.Thread, errors.Err)
	UpsertVote(vote models.Vote) (models.Thread, errors.Err)
	GetThreadPosts(params models.QueryPostParams) (models.PostsList, errors.Err)
}
