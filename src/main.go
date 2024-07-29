package main

import (
	"time"
	"os"
	"vu/ase/distance/src/publisher"
	"vu/ase/distance/src/sensor"

	pb_core_messages "github.com/VU-ASE/rovercom/packages/go/core"
	servicerunner "github.com/VU-ASE/roverlib/src"

	"github.com/rs/zerolog/log"
)

func run(
	service servicerunner.ResolvedService,
	sysMan servicerunner.CoreInfo,
	initialTuning *pb_core_messages.TuningState) error {

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
}

func onTuningState(newtuning *pb_core_messages.TuningState) {
	log.Info().Str("Value", newtuning.String()).Msg("Received tuning state from system manager")
}

func onTerminate(signal os.Signal) {

}

func main() {
	servicerunner.Run(run, onTuningState, onTerminate, false)
}
