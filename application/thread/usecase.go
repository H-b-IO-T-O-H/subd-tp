package thread

import (
	"subd/application/common/errors"
	"subd/application/common/models"
)

type IUseCaseThread interface {
	CreateThread(thread models.Thread) (models.Thread, errors.Err)
	GetBySlugOrId(slug string, id int) (models.Thread, errors.Err)
	UpdateBySlugOrId(threadNew models.Thread, slugId string) (models.Thread, errors.Err)
	CreateVote(vote models.Vote) (models.Thread, errors.Err)
	GetThreadPosts(params models.QueryPostParams) (models.PostsList, errors.Err)
}