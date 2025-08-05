## Positioning example

This example first runs the motor in speedmode with two different speeds / directions and then runs the motor in positioning mode where it toggles between 0 and +90° five times.

## Cross compiling and running the Go example code

1. Set the following environment variables (syntax may vary depending on your OS flavour),

```sh
GOOS=linux 
GOARCH=arm64 
GOARM=8
```

2. Compile the example program 

```sh
go build -o position
```

3. Upload the binary to your board using scp (host name is "overlord" in the example)

```sh
scp position overlord:
```

4. Bring up the CAN bus

```sh
sudo ip link set can0 up type can bitrate 1000000
```

4. ssh to your board and run the scan program (chmod +x if it isn't executable)

```sh
ssh overlord
chmod +x position
./scan
```
