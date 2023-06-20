package repo

import (
	"fmt"
	"id-maker/internal/entity"
	"id-maker/pkg/mysql"
	"time"
)

type SegmentRepo struct {
	*mysql.Mysql
}

func New(mysql *mysql.Mysql) *SegmentRepo {
	return &SegmentRepo{mysql}
}

func (r *SegmentRepo) GetList() ([]entity.Segments, error) {
	var s []entity.Segments
	err := r.Engine.Find(&s)
	if err != nil {
		return s, fmt.Errorf("SegmentRepo - GetList - Find: %w", err)
	}
	return s, nil
}

func (r *SegmentRepo) Add(s *entity.Segments) error {
	var (
		exist bool
		err   error
	)
	exist, err = r.Engine.Where("biz_tag = ?", s.BizTag).Exist(&entity.Segments{})
	if err != nil {
		return fmt.Errorf("SegmentRepo - Add - Exist: %w", err)
	}
	if exist {
		return fmt.Errorf("Tag Already Exist")
	}
	_, err = r.Engine.Insert(s)
	if err != nil {
		return fmt.Errorf("SegmentRepo - Add - Insert: %w", err)
	}
	return nil
}

func (r *SegmentRepo) GetNextId(tag string) (*entity.Segments, error) {
	var (
		err error
		id  = &entity.Segments{}
		tx  = r.Engine.Prepare()
	)

	_, err = tx.Exec("update segments set max_id=max_id+step, update_time = ? where biz_tag = ?", time.Now(), tag)
	if err != nil {
		_ = tx.Rollback()
		return id, fmt.Errorf("SegmentRepo - GetNextId - Exec: %w", err)
	}

	_, err = tx.Where("biz_tag = ?", tag).Get(id)
	if err != nil {
		_ = tx.Rollback()
		return id, fmt.Errorf("SegmentRepo - GetNextId - Get: %w", err)
	}

	err = tx.Commit()

	return id, nil

}
