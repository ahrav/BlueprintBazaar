package implementations

import (
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	nodeBits = 10
	stepBits = 12
	maxStep  = 1<<stepBits - 1
)

var (
	mutex  sync.Mutex // Protects the following variables
	lastTs int64      // Timestamp of last ID generation
	step   int64      // Counter for IDs generated in the same millisecond
)

func generateID() int64 {
	mutex.Lock()
	defer mutex.Unlock()

	ts := time.Now().UnixMilli()

	// if the time in millisecond is the same as last time, and we haven't exhausted the step
	// in this millisecond, increase the step.
	if ts == lastTs && step < maxStep {
		step++
	} else {
		step = 0
		lastTs = ts
	}

	res := ts << (nodeBits + stepBits)

	// NodeID use the MAC address.
	// If the MAC address is not available, use a random number.
	nodeID := getNodeID()
	res |= nodeID << stepBits

	// Step is the number of IDs generated in the same millisecond.
	// If the step is greater than the maximum value, wait until the next millisecond.
	step := rand.Int63n(maxStep)
	return res | step
}

// getNodeID returns the MAC address of the machine if available, otherwise
// returns a random number.
func getNodeID() int64 {
	// Get the MAC address of the machine.
	interfaces, err := net.Interfaces()
	if err != nil {
		return rand.Int63()
	}

	// Get the first non-loopback interface.
	for _, i := range interfaces {
		if i.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Get the MAC address.
		mac := i.HardwareAddr
		if len(mac) == 0 {
			continue
		}

		// Return the MAC address.
		return int64(mac[0])<<40 |
			int64(mac[1])<<32 |
			int64(mac[2])<<24 |
			int64(mac[3])<<16 |
			int64(mac[4])<<8 |
			int64(mac[5])
	}

	// If no MAC address is available, return a random number.
	return rand.Int63()
}
