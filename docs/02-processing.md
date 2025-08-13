# Processing

This service performs the following steps:

1. Send an I2C signal to the URM09 sensor to read the latest distance value
2. This distance is encoded in the [`DistanceSensorOutput` message](https://github.com/VU-ASE/rovercom/blob/c1d6569558e26d323fecc17d01117dbd089609cc/definitions/outputs/distance.proto#L11) which is written to the [`distance` stream](https://github.com/VU-ASE/distance/blob/6f71c8d766f969d5ad76069a870723ac3e49caa4/service.yaml#L14)

## Driver Implementation

We have implemented the communication with the URM09 sensor directly in Go which can be found under [`src/urm09.go`](https://github.com/VU-ASE/distance/blob/main/src/urm09.go). It thus does not deppend on any 3rd party drivers.

## Refresh Rate

You will most likely want to use this service together with an imaging service to brake the Rover when it is close to hitting a wall. Be aware that the refresh rate of this sensor is quite low (but [configurable](https://github.com/VU-ASE/distance/blob/6f71c8d766f969d5ad76069a870723ac3e49caa4/service.yaml#L17)), so you might slow down your service if you read from this service in your main loop.

We advise to use a threaded solution in case you want to use this service. See the "knowing more" series on the ASE website for more information.