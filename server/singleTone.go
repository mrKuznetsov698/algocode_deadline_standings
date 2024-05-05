package server

import (
	"algocode_deadline_standings/configs"
	processors "algocode_deadline_standings/dataProcessor"
	"sync"
	"time"
)

type DataAccess struct {
	lock sync.RWMutex

	err error

	config     *configs.Config
	lastUpdate time.Time

	criterionTitles []*processors.CriterionTitle
	userValues      []*processors.UserValues
	stats           map[int]*processors.Stats
}

var initGlobal sync.Once
var globalDataAccess *DataAccess

func NewDataAccess(config *configs.Config) *DataAccess {
	initGlobal.Do(func() {
		globalDataAccess = &DataAccess{config: config, lastUpdate: time.Unix(0, 0)}
	})
	return globalDataAccess
}

func (data *DataAccess) tick() error {
	data.lock.Lock()
	defer data.lock.Unlock()
	if time.Since(data.lastUpdate).Seconds() > data.config.CacheTime {
		data.criterionTitles, data.userValues, data.err = processors.GetDeadlineResults(data.config)
		if data.err == nil {
			data.stats, data.err = processors.CreateStatistics(data.config, data.userValues)
		}
		data.lastUpdate = time.Now()
	}
	return data.err
}
