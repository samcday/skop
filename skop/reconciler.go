package skop

import "context"

type Reconciler interface {
	Reconcile(ctx context.Context, op *Operator, res Resource) error
}

type ReconcilerFunc func(ctx context.Context, op *Operator, res Resource) error

func (f ReconcilerFunc) Reconcile(ctx context.Context, op *Operator, res Resource) error {
	if op.dryRun {
		ctx = context.WithValue(ctx, "dryrun", true)
	}
	return f(ctx, op, res)
}
