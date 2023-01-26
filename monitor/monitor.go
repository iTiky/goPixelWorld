package monitor

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

type (
	Keeper struct {
		lock sync.Mutex
		//
		particlesCnt     uint64
		particleMovesCnt uint64
		frames           uint64
		//
		opDurations map[string]*opTracking
		opBuffer    []*opTracking
		opJobCh     chan opTrackingJob
		//
		reportPeriod time.Duration
		prevRun      time.Time
	}

	opTracking struct {
		name     string
		cnt      uint64
		duration time.Duration
	}

	opTrackingJob struct {
		name     string
		duration time.Duration
	}
)

func NewKeeper(reportPeriod time.Duration) (*Keeper, error) {
	if reportPeriod <= 0 {
		return nil, fmt.Errorf("reportPeriod must be > 0")
	}

	k := &Keeper{
		opBuffer:     make([]*opTracking, 0, 10),
		opDurations:  make(map[string]*opTracking),
		opJobCh:      make(chan opTrackingJob, 100000),
		reportPeriod: reportPeriod,
		prevRun:      time.Now(),
	}

	return k, nil
}

func (k *Keeper) AddFrame() {
	k.frames++
}

func (k *Keeper) AddParticle() {
	k.particlesCnt++
}

func (k *Keeper) RemoveParticle() {
	k.particlesCnt--
}

func (k *Keeper) AddParticleMove() {
	k.particleMovesCnt++
}

func (k *Keeper) TrackOpDuration(opName string) func() {
	startedAt := time.Now()
	return func() {
		k.opJobCh <- opTrackingJob{
			name:     opName,
			duration: time.Now().Sub(startedAt),
		}
	}
}

func (k *Keeper) Start(ctx context.Context) {
	timer := time.NewTicker(k.reportPeriod)
	go func() {
		for {
			select {
			case opJob := <-k.opJobCh:
				k.lock.Lock()

				opTrack := k.opDurations[opJob.name]
				if opTrack == nil {
					opTrack = &opTracking{
						name: opJob.name,
					}
				}
				opTrack.cnt++
				opTrack.duration += opJob.duration

				k.opDurations[opJob.name] = opTrack

				k.lock.Unlock()
			case <-timer.C:
				k.report()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (k *Keeper) report() {
	k.lock.Lock()
	defer k.lock.Unlock()

	now := time.Now()
	reportDiff := now.Sub(k.prevRun)

	framesPS := float64(k.frames) / reportDiff.Seconds()
	particleMovesPF := float64(k.particlesCnt) / float64(k.frames)

	str := strings.Builder{}

	str.WriteString("=> Monitor report:\n")
	str.WriteString(fmt.Sprintf("  Particles:\t\t%d\n", k.particlesCnt))
	str.WriteString(fmt.Sprintf("  ParticleMoves [PF]:\t%.2f\n", particleMovesPF))
	str.WriteString(fmt.Sprintf("  Frames [PS]:\t\t%.2f\n", framesPS))

	str.WriteString("  Operations:\n")
	for _, opTrack := range k.opDurations {
		k.opBuffer = append(k.opBuffer, opTrack)
	}
	sort.Slice(k.opBuffer, func(i, j int) bool {
		return k.opBuffer[i].name < k.opBuffer[j].name
	})
	for _, opTrack := range k.opBuffer {
		opCntPF := float64(opTrack.cnt) / float64(k.frames)
		opDurAvg := opTrack.duration
		if opTrack.cnt > 0 {
			opDurAvg /= time.Duration(opTrack.cnt)
		}
		opDurPF := opDurAvg * time.Duration(opCntPF)

		str.WriteString(fmt.Sprintf("    %20s:\t%.2f [PF]\t%v [avg]\t%v [total PF]\n", opTrack.name, opCntPF, opDurAvg, opDurPF))

		opTrack.cnt = 0
		opTrack.duration = 0
	}

	k.particleMovesCnt = 0
	k.frames = 0
	k.opBuffer = k.opBuffer[:0]
	k.prevRun = now

	log.Println(str.String())
}
