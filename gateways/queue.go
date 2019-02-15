package gateways

import (
	"context"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/handlers/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"google.golang.org/appengine/taskqueue"
)

type Queue struct {
	dig.In
	RouteResolver handlers.RouteResolver
}

func (q *Queue) EnqueueSendUpdate(ctx context.Context, id git.RepositoryID) error {
	t := taskqueue.Task{
		Method: "POST",
		Path:   q.RouteResolver.TaskSendUpdate(id),
	}
	if _, err := taskqueue.Add(ctx, &t, "SendUpdate"); err != nil {
		return errors.Wrapf(err, "error while adding a SendUpdate task for the repository %s", id)
	}
	return nil
}
