package usecase

import (
	"context"
	"errors"
	"id-maker/internal/entity"
	"id-maker/pkg/snowflake"
	"sync"
	"time"
)

type IdArray struct {
	Cur   int64
	Start int64
	End   int64
}

type BizAlloc struct {
	Mu      sync.RWMutex
	BazTag  string
	IdArray []*IdArray
	GetDb   bool
}

type Alloc struct {
	Mu        sync.RWMutex
	BizTagMap map[string]*BizAlloc
}

func (uc *SegmentUseCase) NewAllocId() (a *Alloc, err error) {
	var res []entity.Segments
	if res, err = uc.repo.GetList(); err != nil {
		return
	}

	a = &Alloc{
		BizTagMap: make(map[string]*BizAlloc),
	}

	for _, v := range res {
		a.BizTagMap[v.BizTag] = &BizAlloc{
			BazTag:  v.BizTag,
			IdArray: make([]*IdArray, 0),
			GetDb:   false,
		}
	}

	return

}

func (uc *SegmentUseCase) NewAllocSnowFlakeId() (*snowflake.Worker, error) {
	return snowflake.NewWorker(1)
}

func (b *BizAlloc) GetId(uc *SegmentUseCase) (id int64, err error) {
	var (
		canGetId    bool
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	)
	b.Mu.Lock()
	if b.LeftIdCount() > 0 {
		canGetId = true
		id = b.PopId()
	}

	if len(b.IdArray) <= 1 && !b.GetDb {
		b.GetDb = true
		b.Mu.Unlock()
		go b.getIdArray(cancel, uc)
	} else {
		b.Mu.Unlock()
		defer cancel()
	}

	if canGetId {
		return
	}
	select {
	case <-ctx.Done():
	}
	b.Mu.Lock()
	if b.LeftIdCount() > 0 {
		id = b.PopId()
	} else {
		err = errors.New("get id error")
	}
	b.Mu.Unlock()
	return
}

func (b *BizAlloc) getIdArray(cancel context.CancelFunc, uc *SegmentUseCase) {
	var (
		tryNum int
		ids    *entity.Segments
		err    error
	)
	defer cancel()
	for {
		if tryNum >= 3 {
			b.GetDb = false
			break
		}
		b.Mu.Lock()

		if len(b.IdArray) <= 1 {
			b.Mu.Unlock()
			ids, err = uc.repo.GetNextId(b.BazTag)
			if err != nil {
				tryNum++
			} else {
				tryNum = 0
				b.Mu.Lock()
				b.IdArray = append(b.IdArray, &IdArray{Start: ids.MaxId, End: ids.MaxId + ids.Step})
				if len(b.IdArray) > 1 {
					b.GetDb = false
					b.Mu.Unlock()
					break
				} else {
					b.Mu.Unlock()
				}
			}

		} else {
			b.Mu.Unlock()
		}
	}
}

func (b *BizAlloc) LeftIdCount() (count int64) {
	for _, v := range b.IdArray {
		arr := v

		count += arr.End - arr.Start - v.Cur
	}
	return count
}

func (b *BizAlloc) PopId() (id int64) {
	id = b.IdArray[0].Start + b.IdArray[0].Cur
	b.IdArray[0].Cur++
	if id+1 >= b.IdArray[0].End {
		b.IdArray = append(b.IdArray[:0], b.IdArray[1:]...)
	}
	return
}
