package main

import (
	"time"
	"vu/ase/distance/src/publisher"
	"vu/ase/distance/src/sensor"

	pb_systemmanager_messages "github.com/VU-ASE/pkg-CommunicationDefinitions/v2/packages/go/systemmanager"
	servicerunner "github.com/VU-ASE/pkg-ServiceRunner/v2/src"

	"github.com/rs/zerolog/log"
)

func run(
	service servicerunner.ResolvedService,
	sysMan servicerunner.SystemManagerInfo,
	initialTuning *pb_systemmanager_messages.TuningState) error {

	// Get the polling rate from the service yaml
	pollDelay, err := servicerunner.GetTuningInt("polling-delay", initialTuning)
	if err != nil {
		return err
	}
	pollRate := time.Duration(pollDelay) * time.Millisecond

	// Start the publisher
	pubAddr, err := service.GetOutputAddress("distance")
	if err != nil {
		return err
	}

	distanceSensor := sensor.NewURM09(pollRate)
	pub := publisher.NewPubDistance(pubAddr, distanceSensor.GetOutgress())

	go distanceSensor.Run()
	go pub.Run()

	select {}

	return nil /* Unreachable */
}

func onTuningState(newtuning *pb_systemmanager_messages.TuningState) {
	log.Info().Str("Value", newtuning.String()).Msg("Received tuning state from system manager")
}

func main() {
	servicerunner.Run(run, onTuningState, false)
}
