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
## Implementation Details

We have implemented the communication with the URM09 sensor directly in Go which can be found under [`src/urm09.go`](https://github.com/VU-ASE/distance/blob/main/src/urm09.go). It thus does not deppend on any 3rd party drivers.
