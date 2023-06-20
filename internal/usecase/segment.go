package usecase

import (
	"id-maker/pkg/snowflake"
)

type SegmentUseCase struct {
	repo SegmentRepo
	//alloc     *Alloc
	snowFlake *snowflake.Worker
}
