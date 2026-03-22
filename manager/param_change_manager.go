package manager

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"qa_test_server/db"
	"qa_test_server/model"
)

var ParamChangeManagerGlobal = &ParamChangeManager{}

type ParamChangeQuery struct {
	DeviceSn    string
	Start       time.Time
	End         time.Time
	PathKeyword string
	Offset      int
	Limit       int
}

type ParamChangeManager struct {
	retention time.Duration

	queue    chan model.DeviceParamChange
	stopCh   chan struct{}
	wg       sync.WaitGroup
	initOnce sync.Once
}

func (m *ParamChangeManager) Init(retention time.Duration) error {
	if db.DB == nil {
		return errors.New("database is unavailable")
	}
	if retention <= 0 {
		retention = 30 * 24 * time.Hour
	}
	if err := db.DB.AutoMigrate(&model.DeviceParamChange{}); err != nil {
		return err
	}

	m.initOnce.Do(func() {
		m.retention = retention
		m.queue = make(chan model.DeviceParamChange, 4096)
		m.stopCh = make(chan struct{})
		m.wg.Add(1)
		go m.loop()
	})
	return nil
}

func (m *ParamChangeManager) loop() {
	defer m.wg.Done()

	flushTicker := time.NewTicker(time.Second)
	cleanupTicker := time.NewTicker(time.Hour)
	defer flushTicker.Stop()
	defer cleanupTicker.Stop()

	batch := make([]model.DeviceParamChange, 0, 256)
	flush := func() {
		if len(batch) == 0 {
			return
		}
		if err := db.DB.Create(&batch).Error; err != nil {
			log.Printf("param-change flush failed: %v", err)
		}
		batch = batch[:0]
	}

	for {
		select {
		case <-m.stopCh:
			flush()
			return
		case item := <-m.queue:
			batch = append(batch, item)
			if len(batch) >= 256 {
				flush()
			}
		case <-flushTicker.C:
			flush()
		case <-cleanupTicker.C:
			if err := m.PurgeExpired(context.Background()); err != nil {
				log.Printf("param-change purge failed: %v", err)
			}
		}
	}
}

func (m *ParamChangeManager) Stop() {
	if m.stopCh == nil {
		return
	}
	select {
	case <-m.stopCh:
		return
	default:
		close(m.stopCh)
	}
	m.wg.Wait()
}

func (m *ParamChangeManager) Enqueue(items []model.DeviceParamChange) {
	if len(items) == 0 || m.queue == nil {
		return
	}
	for _, item := range items {
		select {
		case m.queue <- item:
		default:
			log.Printf("param-change queue full, drop log for device=%s path=%s", item.DeviceSn, item.ParamPath)
		}
	}
}

func (m *ParamChangeManager) PurgeExpired(ctx context.Context) error {
	if db.DB == nil || m.retention <= 0 {
		return nil
	}
	deadline := time.Now().Add(-m.retention)
	return db.DB.WithContext(ctx).Where("changed_at < ?", deadline).Delete(&model.DeviceParamChange{}).Error
}

func (m *ParamChangeManager) ClearAll(ctx context.Context) (int64, error) {
	if db.DB == nil {
		return 0, errors.New("database is unavailable")
	}

	if m.queue != nil {
		for {
			select {
			case <-m.queue:
				continue
			default:
				goto drained
			}
		}
	}

drained:
	res := db.DB.WithContext(ctx).Where("1 = 1").Delete(&model.DeviceParamChange{})
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (m *ParamChangeManager) Query(q ParamChangeQuery) ([]model.DeviceParamChange, int64, error) {
	if db.DB == nil {
		return nil, 0, errors.New("database is unavailable")
	}
	sn := strings.TrimSpace(q.DeviceSn)
	if sn == "" {
		return nil, 0, errors.New("device sn is required")
	}

	if q.Offset < 0 {
		q.Offset = 0
	}
	if q.Limit <= 0 {
		q.Limit = 200
	}
	if q.Limit > 5000 {
		q.Limit = 5000
	}

	end := q.End
	if end.IsZero() {
		end = time.Now()
	}
	start := q.Start
	if start.IsZero() {
		start = end.Add(-24 * time.Hour)
	}
	if start.After(end) {
		start, end = end, start
	}

	maxRange := m.retention
	if maxRange <= 0 {
		maxRange = 30 * 24 * time.Hour
	}
	if end.Sub(start) > maxRange {
		start = end.Add(-maxRange)
	}

	query := db.DB.Model(&model.DeviceParamChange{}).
		Where("device_sn = ? AND changed_at >= ? AND changed_at <= ?", sn, start, end)
	if keyword := strings.TrimSpace(q.PathKeyword); keyword != "" {
		query = query.Where("param_path LIKE ?", "%"+keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	items := make([]model.DeviceParamChange, 0, q.Limit)
	if err := query.Order("changed_at DESC").Offset(q.Offset).Limit(q.Limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
