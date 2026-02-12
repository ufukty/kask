package memo

import (
	"fmt"
	"runtime"
)

type measurement struct {
	checkpoint string
	alloc      uint64
	sys        uint64
}

func stats() *runtime.MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return &m
}

func measure(checkpoint string) *measurement {
	s := stats()
	return &measurement{
		checkpoint: checkpoint,
		alloc:      s.Alloc,
		sys:        s.Sys,
	}
}

type Accountant struct {
	measurements []*measurement
}

func (a *Accountant) Check(checkpoint string) {
	a.measurements = append(a.measurements, measure(checkpoint))
}

func (a *Accountant) PeakAlloc() (*measurement, error) {
	var m *measurement
	for _, m2 := range a.measurements {
		if m == nil || m2.alloc > m.alloc {
			m = m2
		}
	}
	if m == nil {
		return nil, fmt.Errorf("no measurement found. is there any checkpoints?")
	}
	return m, nil
}

func NewAccountant() *Accountant {
	return &Accountant{
		measurements: []*measurement{},
	}
}
