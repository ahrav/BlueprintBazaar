package parkinglot

import (
	"container/heap"
	"sync"
	"time"
)

// ParkingLot (singleton)
// ParkingLotConfig (builder)
// ParkingSpot
// ParkingSpotStrategy -> NearestSpotStrategy
// ParkingLotManager
// PaymentManager
// PaymentStrategy -> FlatRate / TimeBased

type SpotType string

const (
	CompactSpotType     SpotType = "compact"
	ReservedSpotType    SpotType = "reserved"
	HandicappedSpotType SpotType = "handicapped"
)

type ParkingSpot struct {
	ID         string
	Type       SpotType
	IsOccupied bool
	StartTime  time.Time
}

func (ps *ParkingSpot) IsAvailable() bool {
	return !ps.IsOccupied
}

func (ps *ParkingSpot) Vacate() error {
	ps.IsOccupied = false
	return nil
}

func (ps *ParkingSpot) Allocate() error {
	ps.StartTime = time.Now()
	ps.IsOccupied = true
	return nil
}

type ParkingSpotStrategy interface {
	FindSpot() *ParkingSpot
}

type NearestSpotStrategy struct {
	allSpots       heap.Interface // Minheap
	availableSpots map[string]*ParkingSpot
	occupiedSpots  map[string]*ParkingSpot
}

func (nss *NearestSpotStrategy) FindSpot() *ParkingSpot {
	// implement
	return nil
}

type PaymentStrategy interface {
	CalculatePayment(spot *ParkingSpot) int
}

type FlatRatePayment struct {
	rates map[SpotType]int
}

func (frp *FlatRatePayment) CalculatePayment(spot *ParkingSpot) int {
	return frp.rates[spot.Type]
}

type PaymentManager struct {
	strategy PaymentStrategy
}

func (pm *PaymentManager) CalculateFee(spot *ParkingSpot) int {
	return pm.strategy.CalculatePayment(spot)
}

func (pm *PaymentManager) ProcessPayment(fee int) error {
	return nil
}

type ParkingManager struct {
	strategy ParkingSpotStrategy
}

func (pm *ParkingManager) AllocateSpot() *ParkingSpot {
	return pm.strategy.FindSpot()
}

func (pm *ParkingManager) Free(spot *ParkingSpot) error {
	return spot.Vacate()
}

type ParkingLot struct {
	spots      []*ParkingSpot
	parkingMgr *ParkingManager
	paymentMgr *PaymentManager
}

var instance *ParkingLot
var once sync.Once

func GetParkingLot(spots []*ParkingSpot, paymentMgr *PaymentManager, parkingMgr *ParkingManager) *ParkingLot {
	once.Do(func() {
		instance = &ParkingLot{
			spots:      spots,
			parkingMgr: parkingMgr,
			paymentMgr: paymentMgr,
		}
	})
	return instance
}

type ParkingLotConfig struct {
	NumCompactSpots     int
	NumReservedSpots    int
	NumHandicappedSpots int
	PaymentStrategy     PaymentStrategy
	ParkingStrategy     ParkingSpotStrategy
}

type ParkingLotBuilder struct {
	config ParkingLotConfig
}

func NewParkingLotBuilder(cfg ParkingLotConfig) ParkingLotBuilder {
	return ParkingLotBuilder{config: cfg}
}

func (pb *ParkingLotBuilder) Build() *ParkingLot {
	// Create spots.

	var spots []*ParkingSpot
	parkingMgr := &ParkingManager{strategy: pb.config.ParkingStrategy}
	paymentMgr := &PaymentManager{strategy: pb.config.PaymentStrategy}
	return GetParkingLot(spots, paymentMgr, parkingMgr)
}
