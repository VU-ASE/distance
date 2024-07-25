package shareddata

import (
	pb_module_outputs "github.com/VU-ASE/pkg-CommunicationDefinitions/v2/packages/go/outputs"
)

type DistanceValue = pb_module_outputs.DistanceSensorOutput
type DistanceChan = chan *DistanceValue
