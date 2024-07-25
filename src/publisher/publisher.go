package publisher

import (
	"time"

	pb_module_outputs "github.com/VU-ASE/pkg-CommunicationDefinitions/v2/packages/go/outputs"
	"google.golang.org/protobuf/proto"

	shareddata "vu/ase/distance/src/shareddata"

	zmq "github.com/pebbe/zmq4"
	"github.com/rs/zerolog/log"
)

// PubDistance is a subscriber that listens to the mod-imaging distance output
type PubDistance struct {
	ingress shareddata.DistanceChan

	// ZMQ distance server
	serverString string
}

// Create new publisher instance and return its pointer
func NewPubDistance(serverAddr string, in shareddata.DistanceChan) *PubDistance {
	p := &PubDistance{}
	var err = p.Init(serverAddr, in)
	if err != nil {
		return nil
	}
	return p
}

func (p *PubDistance) Init(serverAddr string, in shareddata.DistanceChan) error {
	p.serverString = serverAddr
	p.ingress = in

	return nil // No possible errors
}

// Run starts the distance client, calling and handling errors of Start()
func (p *PubDistance) Run() {
	log.Info().Msg("Starting distance client")
	err := p.Start()
	if err != nil {
		log.Err(err).Msg("Failed to start distance client")
	}
}

// Start starts the distance publisher, reading from the ingress
func (p *PubDistance) Start() error {
	publisher, _ := zmq.NewSocket(zmq.PUB)
	defer publisher.Close()

	err := publisher.Bind(p.serverString)
	if err != nil {
		log.Error().Err(err).Str("address", p.serverString).Msg("Failed to bind to zmq address")
		return err
	}

	log.Info().Str("address", "tcp://*:5556").Msg("Publishing zmq messages")

	for {
		distance := <-p.ingress

		// Create the protobuf message
		protoSensorDistanceData := &pb_module_outputs.SensorOutput{
			SensorId:  1010,
			Timestamp: uint64(time.Now().UnixMilli()),
			Status:    0,
			SensorOutput: &pb_module_outputs.SensorOutput_DistanceOutput{
				DistanceOutput: distance,
			},
		}

		distanceBytes, err := proto.Marshal(protoSensorDistanceData)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal protobuf message")
			return err
		}

		_, err = publisher.SendBytes(distanceBytes, 0)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send protobuf message")
			return err
		}
	}
}
