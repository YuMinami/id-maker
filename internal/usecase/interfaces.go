package usecase

import "id-maker/internal/entity"

type (
	Segment interface {
		CreateTag(segments *entity.Segments) error
		GetId(string) (int64, error)
		SnowFlakeGetId() int64
	}
	SegmentRepo interface {
		GetList() ([]entity.Segments, error)
		GetNextId(string) (*entity.Segments, error)
		Add(*entity.Segments) error
	}
)
