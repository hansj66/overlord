# IMU

Note: AFAIK there is no kernel driver available for the BNO055 IMU, so the low level i2c interface must be used to communicate witht the IMU.

## Cross compiling the Go example code

1. Set the following environment variables (syntax may vary depending on your OS flavour),

```sh
GOOS=linux 
GOARCH=arm64 
GOARM=8
```

2. Compile the example program 

```sh
go build -o euler
```

3. Upload the binary to your board using scp (host name is "overlord" in the example)

```sh
scp euler overlord:
```

4. ssh to your board and run the euler program (chmod +x if it isn't executable)

```sh
ssh overlord
chmod +x euler
./euler
```

Move the carrier board around to see the heading vector change.
Example output:

```sh
BNO055 sensor detected, chip ID is correct (0xA0)
Heading: 360.00°, Roll: 0.38°, Pitch: 0.31°
Heading: 360.00°, Roll: 0.38°, Pitch: 0.31°
Heading: 359.00°, Roll: 0.38°, Pitch: 0.19°
Heading: 344.62°, Roll: 4078.50°, Pitch: 4083.62°
Heading: 240.50°, Roll: 4045.25°, Pitch: 4056.69°
Heading: 238.50°, Roll: 4044.38°, Pitch: 4044.31°
Heading: 332.50°, Roll: 4058.62°, Pitch: 4068.19°
Heading: 329.62°, Roll: 4090.62°, Pitch: 25.81°
Heading: 330.44°, Roll: 4094.88°, Pitch: 34.25°
Heading: 330.94°, Roll: 4094.81°, Pitch: 24.50°
Heading: 332.69°, Roll: 0.19°, Pitch: 0.12°
```



### Resources

* [BNO055 Datasheet](https://www.bosch-sensortec.com/media/boschsensortec/downloads/datasheets/bst-bno055-ds000.pdf)
* [BNO055 no-os driver](https://github.com/boschsensortec/BNO055_SensorAPI)