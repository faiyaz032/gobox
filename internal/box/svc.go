package box

type Svc struct {
	repo Repo
}

func NewSvc(repo Repo) *Svc {
	return &Svc{
		repo: repo,
	}
}
