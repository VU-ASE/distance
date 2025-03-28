package main

import (
	"encoding/binary"
	"github.com/rs/zerolog/log"
	smbus "github.com/corrupt/go-smbus"
)

const (
	modeRegister          = 0x07
	commandRegister       = 0x08
	highOrderBitsRegister = 0x03
	lowOrderBitsRegister  = 0x04
)

// Device wraps an I2C connection to a MPU6050 device.
// based on: https://wiki.dfrobot.com/URM09_Ultrasonic_Sensor_(Gravity-I2C)_(V1.0)_SKU_SEN0304
type URM09 struct {
	bus *smbus.SMBus
}

// New creates a new URM09 connection. The I2C bus must already be
// configured.
//
// This function only creates the Device object, it does not touch the device.
func new(bus *smbus.SMBus) *URM09 {
	return &URM09{bus}
}

func (u *URM09) ReadDistance() (int, error) {
	// Initialize a ranging request
	u.bus.Write_byte_data(commandRegister, 0x01)

	// Read the results
	lowBits, err := u.bus.Read_byte_data(lowOrderBitsRegister)
	if err != nil {
		return 0, err
	}
	highBits, err := u.bus.Read_byte_data(highOrderBitsRegister)
	if err != nil {
		return 0, err
	}

	// Put it in little endian representation
	data := make([]byte, 2)
	data[0] = lowBits
	data[1] = highBits

	distance := binary.LittleEndian.Uint16(data)
	return int(distance), nil
}

func (u *URM09) EnablePassiveMode() error {
	// Set the board in passive mode so that it does a range detection on request,
	// instead of continuously
	return u.bus.Write_byte_data(modeRegister, 0x0)
}

// Initialize a URM09 device on the specified I2C bus and address (default is 0x11)
func Initialize(bus uint, address uint8) *URM09 {
	i2c, err := smbus.New(bus, address)
	if err != nil {
		log.Err(err).Msg("Failed to initialize i2c")
		return nil
	}

	urm := new(i2c)
	if err != nil {
		log.Err(err).Msg("Failed to initialize urm on the bus")
		return nil
	}

	err = urm.EnablePassiveMode()
	if err != nil {
		log.Err(err).Msg("Failed to enable passive mode")
		return nil
	}

	log.Debug().Msgf("Initialized URM09 controller on bus %d (0x%x)", bus, address)
	return urm
}