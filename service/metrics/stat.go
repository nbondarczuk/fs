package metrics

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	. "fs/service/config"
)

type Stat struct {
	Alloc      uint64
	TotalAlloc uint64
	Sys        uint64
	NumGC      uint32
}

const (
	ServiceAlive   int32 = 1
	ServiceHealthy int32 = 2
	ServiceError   int32 = 3
	ServiceDead    int32 = 4
)

var (
	alive, healthy int32
	requestId      string
)

//
// bToMb connverts bytes to Mbytes
//
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

//
// Info provides mem stat info
//
func StatInfo() Stat {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	s := Stat{
		bToMb(m.Alloc),
		bToMb(m.TotalAlloc),
		bToMb(m.Sys),
		m.NumGC,
	}

	return s
}

func SetAlive(s int32) {
	Log.Debug("Alive status: -> " + fmt.Sprintf("%d", s))
	atomic.StoreInt32(&alive, s)
}

func IsAlive() bool {
	var val = atomic.LoadInt32(&alive)
	Log.Debug("Alive status: " + fmt.Sprintf("%d", val))
	return val == ServiceAlive
}

func SetHealthy(s int32) {
	Log.Debug("Health status: -> " + fmt.Sprintf("%d", s))
	atomic.StoreInt32(&healthy, s)
}

func IsHealthy() bool {
	var val = atomic.LoadInt32(&healthy)
	Log.Debug("Health status: " + fmt.Sprintf("%d", val))
	return val == ServiceHealthy
}

func RequestId() string {
	requestId = fmt.Sprintf("%d", time.Now().UnixNano())
	Log.Debug("Current requestId: " + requestId)
	return requestId
}
