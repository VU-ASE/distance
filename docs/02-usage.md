# Usage


## Output

This service produces a `distance` stream, containing an output value from the sensor in **meters**. It is contained in `DistanceSensorOutput` within a `SensorOutput_DistanceOutput` object. Code responsible for writing this object can be seen below:

```
err = writeStream.Write(
    &pb_outputs.SensorOutput{
        SensorId:  2,
        Status: 0,
        Timestamp: uint64(time.Now().UnixMilli()),
        SensorOutput: &pb_outputs.SensorOutput_DistanceOutput{
            DistanceOutput: &pb_outputs.DistanceSensorOutput{
                Distance: float32(distance) / 100.0,
            },
        },
    },
)
```
## Driver

This service requires a URM09 driver to work. The driver contains functions necessary to initialize the sensor and read data from it. This is abstracted away and you do not need to worry about it. In order to initialize the sensor, first an instance of `URM09` is defined by passing the default values for bus number and the address of the sensor. Keep in mind that changing these values may result in you (or the debix) not being able to find the device on the bus. To read the distance value, use `ReadDistance()` method. 