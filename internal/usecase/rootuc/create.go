package rootuc

import (
	"context"
	"github.com/puny-activity/file-service/internal/entity/root"
)

func (u *UseCase) Create(ctx context.Context, rootToCreate root.Root) (root.Root, error) {
	return root.Root{}, nil
}
