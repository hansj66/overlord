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
2025/08/01 16:47:14 INFO BNO055 heading="Heading: 357.75°, Roll: 4066.62°, Pitch: 11.19°\n"
2025/08/01 16:47:14 INFO BNO055 heading="Heading: 359.56°, Roll: 4037.38°, Pitch: 43.06°\n"
2025/08/01 16:47:15 INFO BNO055 heading="Heading: 358.31°, Roll: 4027.44°, Pitch: 21.00°\n"
2025/08/01 16:47:15 INFO BNO055 heading="Heading: 340.94°, Roll: 4033.19°, Pitch: 4028.31°\n"
2025/08/01 16:47:16 INFO BNO055 heading="Heading: 76.75°, Roll: 4048.81°, Pitch: 4020.38°\n"
2025/08/01 16:47:16 INFO BNO055 heading="Heading: 81.56°, Roll: 4067.94°, Pitch: 4029.06°\n"
2025/08/01 16:47:17 INFO BNO055 heading="Heading: 73.81°, Roll: 4089.69°, Pitch: 4020.50°\n"
2025/08/01 16:47:17 INFO BNO055 heading="Heading: 74.81°, Roll: 4085.62°, Pitch: 4060.06°\n"
2025/08/01 16:47:18 INFO BNO055 heading="Heading: 82.56°, Roll: 4089.88°, Pitch: 1.00°\n"
```



### Resources

* [BNO055 Datasheet](https://www.bosch-sensortec.com/media/boschsensortec/downloads/datasheets/bst-bno055-ds000.pdf)
* [BNO055 no-os driver](https://github.com/boschsensortec/BNO055_SensorAPI)