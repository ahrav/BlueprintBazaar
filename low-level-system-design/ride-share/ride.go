package rideshare

import (
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
)

type Price struct {
	Amount   int64 // Smallest units (eg cents)
	Currency string
}

type Driver struct {
	id        string
	name      string
	available bool
	location  Location
}

func (d *Driver) UpdateAvailability(availability bool) {
	d.available = availability
}

func (d *Driver) UpdateLocation(location Location) {
	d.location = location
}

type Rider struct {
	id   string
	name string
}

type Location struct {
	latitude  float64
	longitude float64
}

func (l *Location) Distance(dst Location) float64 {
	return math.Sqrt(math.Pow(l.latitude-dst.latitude, 2) + math.Pow(l.longitude-dst.longitude, 2))
}

type TripStatus int

const (
	TripStatusInProgess = iota
	TripStatusFinished
)

type Trip struct {
	id       string
	status   TripStatus
	rider    Rider
	driver   Driver
	src      Location
	dst      Location
	distance int
	cost     Price
}

func NewTrip(rider Rider, driver Driver, src, dst Location, price Price) Trip {
	dist := math.Sqrt(math.Pow(src.latitude-dst.latitude, 2) + math.Pow(src.longitude-dst.longitude, 2))
	return Trip{
		id:       uuid.New().String(),
		status:   TripStatusInProgess,
		rider:    rider,
		driver:   driver,
		src:      src,
		dst:      dst,
		distance: int(dist),
		cost:     price,
	}
}

func (t *Trip) EndTrip() {
	t.status = TripStatusFinished
}

type DriverService struct {
	drivers map[string]*Driver
}

func NewDriverService() *DriverService {
	return &DriverService{drivers: make(map[string]*Driver)}
}

var ErrDriverExists = errors.New("Unable to create new Driver, Driver exists")

func (ds *DriverService) CreateDriver(name string) (*Driver, error) {
	if _, ok := ds.drivers[name]; ok {
		return nil, ErrDriverExists
	}

	return &Driver{
		id:        uuid.New().String(),
		name:      name,
		available: true,
	}, nil
}

func (ds *DriverService) GetDriver(id string) (*Driver, error) {
	// Implement.
	return nil, nil
}

func (ds *DriverService) UpdateAvailability(id string, isAvailable bool) (bool, error) {
	// Implement.
	return true, nil
}

func (ds *DriverService) UpdateLocation(id string, location Location) (bool, error) {
	// Implement.
	return true, nil
}

var ErrNoNearbyDrivers = errors.New("No nearby drivers available")

func (ds *DriverService) AvailableDrivers(location Location) ([]*Driver, error) {
	var result []*Driver
	for _, driver := range ds.drivers {
		if driver.available && driver.location.Distance(location) <= 10 {
			result = append(result, driver)
		}
	}

	if len(result) == 0 {
		return nil, ErrNoNearbyDrivers
	}

	// Implement.
	return result, nil
}

type RiderService struct {
	riders map[string]*Rider
}

func NewRiderService() *RiderService {
	return &RiderService{riders: make(map[string]*Rider)}
}

func (rs *RiderService) CreateRider(name string) (*Rider, error) {
	// Implement
	return nil, nil
}

func (rs *RiderService) GetRiders(id string) (*Rider, error) {
	// Implement
	return nil, nil
}

type PricingStrategy interface {
	FindPrice(src, dst Location) (Price, error)
}

type DriverMatchingStrategy interface {
	MatchDriver(rider *Rider, drivers []*Driver, src, dst Location) (*Driver, error)
}

type TripService struct {
	trips          map[string][]Trip
	driversManager DriverService
	ridersManager  RiderService

	priceMatcher  PricingStrategy
	driverMatcher DriverMatchingStrategy
}

func NewTripService(
	driversManager DriverService,
	ridersManager RiderService,
	priceMatcher PricingStrategy,
	driverMatcher DriverMatchingStrategy,
) *TripService {
	return &TripService{
		driversManager: driversManager,
		ridersManager:  ridersManager,
		priceMatcher:   priceMatcher,
		driverMatcher:  driverMatcher,
	}
}

func (ts *TripService) CreateTrip(rider *Rider, src, dst Location) error {
	drivers, err := ts.driversManager.AvailableDrivers(src)
	if err != nil {
		return err
	}

	driver, err := ts.driverMatcher.MatchDriver(rider, drivers, src, dst)
	if err != nil {
		return err
	}

	price, err := ts.priceMatcher.FindPrice(src, dst)
	if err != nil {
		return err
	}

	trip := NewTrip(*rider, *driver, src, dst, price)
	ts.trips[rider.id] = append(ts.trips[rider.id], trip)
	driver.UpdateAvailability(false)

	return nil
}

func (ts *TripService) EndTrip(driver *Driver) error {
	if driver.available {
		return fmt.Errorf("Driver has no current trip")
	}

	driver.UpdateAvailability(true)
	return nil
}
