package main

import (
	"fmt"
	"os"
	"time"

	pb_outputs "github.com/VU-ASE/rovercom/v2/packages/go/outputs"
	roverlib "github.com/VU-ASE/roverlib-go/v2/src"
	"github.com/rs/zerolog/log"
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
	urm      *URM09
	pollRate time.Duration
}

// NewURM09 creates a new URM09 sensor with default values
func NewURM09(pollRate time.Duration) *urm09 {
	u := &urm09{
		bus:      defaultBus,
		address:  defaultAddress,
		pollRate: pollRate,
		urm:      Initialize(defaultBus, defaultAddress),
	}
	if u.urm == nil {
		return nil
	}
	return u
}

func run(service roverlib.Service, configuration *roverlib.ServiceConfiguration) error {
	if configuration == nil {
		return fmt.Errorf("configuration cannot be accessed")
	}

	pollDelay, err := configuration.GetFloatSafe("polling-delay")
	if err != nil {
		return fmt.Errorf("failed to get configuration: %v", err)
	}
	log.Info().Msgf("Fetched runtime configuration polling delay: %f", pollDelay)

	writeStream := service.GetWriteStream("distance")
	if writeStream == nil {
		return fmt.Errorf("failed to get write stream")
	}

	// initialize the urm09 sensor
	pollRate := time.Duration(pollDelay) * time.Millisecond
	distanceSensor := NewURM09(pollRate)
	if distanceSensor == nil {
		return fmt.Errorf("failed to initialize distance sensor")
	}

	err = distanceSensor.urm.EnablePassiveMode()
	if err != nil {
		log.Err(err).Msg("Failed to enable passive mode")
	}

	for {
		// read the distance measured by the sensor
		distance, err := distanceSensor.urm.ReadDistance()
		if err != nil {
			log.Info().Msg("Failed to read distance")
			time.Sleep(distanceSensor.pollRate)
			continue
		} else if distance > maxDistanceCm {
			log.Info().Int("distance", distance).Msg("Distance too large. Setting to max distance...")
			distance = maxDistanceCm
		}

		// Send it for the actuator (and others) to use
		err = writeStream.Write(
			&pb_outputs.SensorOutput{
				SensorId:  2,
				Status:    0,
				Timestamp: uint64(time.Now().UnixMilli()),
				SensorOutput: &pb_outputs.SensorOutput_DistanceOutput{
					DistanceOutput: &pb_outputs.DistanceSensorOutput{
						Distance: float32(distance) / 100.0,
					},
				},
			},
		)
		if err != nil {
			log.Err(err).Msg("Failed to send controller output")
			continue
		}

		log.Info().Int("distance", distance).Msg("Distance read")

		time.Sleep(distanceSensor.pollRate)
	}
}

// This function gets called when roverd wants to terminate the service
func onTerminate(sig os.Signal) error {
	log.Info().Str("signal", sig.String()).Msg("Terminating service")
	return nil
}

// This is just a wrapper to run the user program
func main() {
	roverlib.Run(run, onTerminate)
}
