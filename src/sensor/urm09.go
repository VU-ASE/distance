package sensor

import (
	"time"

	urm09_driver "github.com/MrBuggy-Amsterdam/go-urm09driver"
	"github.com/rs/zerolog/log"

	shareddata "vu/ase/distance/src/shareddata"
)

const (
	defaultBus     = 5
	defaultAddress = 0x11
	maxDistanceCm  = 500
)

// urm09 is a struct that represents the URM09 sensor wrapper for this project
type urm09 struct {
	bus      uint
	address  uint8
	urm      *urm09_driver.URM09
	pollRate time.Duration

	// An outgress to write values to
	outgress shareddata.DistanceChan
}

// NewURM09 creates a new URM09 sensor with default values
func NewURM09(pollRate time.Duration) *urm09 {
	u := &urm09{
		bus:      defaultBus,
		address:  defaultAddress,
		pollRate: pollRate,
		urm:      urm09_driver.Initialize(defaultBus, defaultAddress),
	}

	u.outgress = make(shareddata.DistanceChan)

	return u
}

// GetOutgress returns the outgress channel
func (u *urm09) GetOutgress() shareddata.DistanceChan {
	return u.outgress
}

// Run starts reading the distance sensor and writing the values to the outgress channel
func (u *urm09) Run() {
	urm := urm09_driver.Initialize(u.bus, uint8(u.address))

	// Enable passive mode for  on-demand ranging
	err := urm.EnablePassiveMode()
	if err != nil {
		log.Err(err).Msg("Failed to enable passive mode")
		// Do something with the error
	}

	for {
		distance, err := urm.ReadDistance()
		if err != nil {
			log.Err(err).Msg("Failed to read distance")
			time.Sleep(u.pollRate)
			continue
		} else if distance > maxDistanceCm {
			log.Debug().Int("distance", distance).Msg("Distance too large. Setting to max distance...")
			distance = maxDistanceCm
		}

		log.Debug().Int("distance", distance).Msg("Distance read")
		u.outgress <- &shareddata.DistanceValue{
			Distance: float32(distance),
		}

		time.Sleep(u.pollRate)
	}
}
