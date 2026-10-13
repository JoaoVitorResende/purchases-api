package purchase

import "context"

type Repository interface {
	Save(ctx context.Context, p *Purchase) (int64, error)
}