package domain

import "time"

type RepositoryLastScan struct {
	Repository      RepositoryID
	LastScanTime    time.Time
	PreconditionOut GradleUpdatePreconditionOut
}
