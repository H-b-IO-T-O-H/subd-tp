package usecase

import (
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/post"
)

type PostUseCase struct {
	reposPost post.IRepositoryPost
}

func NewPostUseCase(repPost post.IRepositoryPost) post.IUseCasePost {
	return &PostUseCase{reposPost: repPost}
}

func (p PostUseCase) GetPost(params models.PostGetParams) (models.PostFull, errors.Err) {
	return p.reposPost.GetFull(params)
}

func (p PostUseCase) CreatePost(posts models.PostsList, slugId string) (models.PostsList, errors.Err) {
	return p.reposPost.CreatePost(posts, slugId)
}

func (p PostUseCase) UpdatePost(postUpdate models.PostUpdate) (models.Post, errors.Err) {
	return p.reposPost.UpdatePost(postUpdate)
}
