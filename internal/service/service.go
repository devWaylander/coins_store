package service

type Repository interface {
}

type service struct {
	repo   Repository
	jwtKey string
}

func New(repo Repository, jwtKey string) *service {
	return &service{
		repo:   repo,
		jwtKey: jwtKey,
	}
}
