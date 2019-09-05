package goprofiler

import (
	"fmt"
	"stash.tutu.ru/golang/log"
	"sync"
	"time"
)

type Profiler struct {
	startedTags map[string]time.Time
	Tags        map[string][]int64
	startedTime time.Time
	mtx *sync.Mutex
}

func (p *Profiler) Start(tag string) {
	if _, ok := p.startedTags[tag]; ok {
		p.Stop(tag)
	}
	p.startedTags[tag] = time.Now()
}

func (p *Profiler) Stop(tag string) {
	if _, ok := p.startedTags[tag]; ok {
		p.Tags[tag] = append(p.Tags[tag], time.Now().UnixNano()-p.startedTags[tag].UnixNano())
		delete(p.startedTags, tag)
	}
}

func (p *Profiler) Print() {
	for tag := range p.startedTags {
		p.Stop(tag)
	}

	var fullDuration int64 = time.Now().UnixNano() - p.startedTime.UnixNano()

	for tag, durations := range p.Tags {
		var durationSum int64 = 0
		for _, duration := range durations {
			durationSum += duration
		}
		log.Logger.Info().Msg(fmt.Sprintf(`%s: %.2f%%`, tag, float64(durationSum)/float64(fullDuration)*100))
	}

	for tag, durations := range p.Tags {
		var durationSum int64 = 0
		for _, duration := range durations {
			durationSum += duration
		}
		log.Logger.Info().Msg(fmt.Sprintf(`%s: %.2f s`, tag, float64(durationSum)/1000000000))
	}

	log.Logger.Info().Msg(fmt.Sprintf(`%s: %.2f s`, "total", float64(fullDuration)/1000000000))

	p.Reset()
}

func (p *Profiler) Reset() {
	p.startedTime = time.Now()
	p.Tags = make(map[string][]int64)
	p.startedTags = make(map[string]time.Time)
}

var profiler *Profiler

func GetProfiler() *Profiler {
	if profiler == nil {
		profiler = &Profiler{
			make(map[string]time.Time),
			make(map[string][]int64),
			time.Now(),
			 &sync.Mutex{},
		}
	}
	return profiler
}
