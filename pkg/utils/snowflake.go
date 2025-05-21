package utils

import (
	"github.com/Alf_Grindel/save/pkg/utils/hlog"
	"sync"
	"time"
)

const (
	epoch          = 903024000000
	timestampBits  = 41
	machineIDBits  = 10
	sequenceBits   = 12
	maxTimestamp   = -1 ^ (-1 << timestampBits)
	maxMachineID   = -1 ^ (-1 << machineIDBits)
	maxSequenceNum = -1 ^ (-1 << sequenceBits)
	machineShift   = sequenceBits
	timestampShift = sequenceBits + machineIDBits
)

type Snowflake struct {
	mu          sync.Mutex
	timestamp   int64
	machineID   int64
	sequenceNum int64
}

func NewSnowflake(machineId int64) *Snowflake {
	if machineId < 0 || machineId > maxMachineID {
		hlog.Fatalf("machine must between 0 and %d\n", maxMachineID-1)
	}
	return &Snowflake{
		timestamp:   0,
		machineID:   machineId,
		sequenceNum: 0,
	}
}

func (s *Snowflake) GenerateID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UnixNano() / 1e6
	if s.timestamp == now {
		s.sequenceNum = (s.sequenceNum + 1) & maxSequenceNum
		if s.sequenceNum == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequenceNum = 0
	}

	t := now - epoch
	if t > maxTimestamp {
		hlog.Fatalf("epoch must between 0 and %d\n", maxTimestamp-1)
	}
	s.timestamp = now
	r := ((t) << timestampShift) | (s.machineID << machineShift) | s.sequenceNum
	return r
}
