# Dynamixel Protocol 1.0 network scan

This example demonstrates how to scan the valid servo ID range. This is done by sending a ping packet to all addresses between 0 and 253. A response packet indicates that a servo has been found in the network

This example is using Dynamixel protocol 1.0 for sending ping packets. The example can be modified for using protocol 2.0 by changing the frame format and CRC calculation. Protocol 2.0 ping is described in section 5.1.3.2 in this manual: https://emanual.robotis.com/docs/en/dxl/protocol2/


## Cross compiling the Go example code

1. Set the following environment variables (syntax may vary depending on your OS flavour),

```sh
GOOS=linux 
GOARCH=arm64 
GOARM=8
```

2. Compile the example program 

```sh
go build -o scan
```

3. Upload the binary to your board using scp (host name is "overlord" in the example)

```sh
scp scan overlord:
```

4. ssh to your board and run the scan program (chmod +x if it isn't executable)

```sh
ssh overlord
chmod +x scan
./scan
```

Assuming that the network contains one servo with ID=1, the output will be on the form:

```sh
2025/08/01 14:09:41 INFO Opening serialport port=/dev/ttyAMA0
2025/08/01 14:09:41 INFO Serial port opened successfully
2025/08/01 14:09:41 INFO DXL DIR pin=GPIO22
PING id:0 - no response
PING id:1 - response packet: FFFF010200FC
PING id:2 - no response
PING id:3 - no response
...
```

The same protocol can be used for reading and writing to EEPROM and RAM control table parameters
Protocol specifications:
 - 1.0: https://emanual.robotis.com/docs/en/dxl/protocol1/
 - 2.0: https://emanual.robotis.com/docs/en/dxl/protocol2/

Note:

1. The carrier board has to be powered through the "COMPUTE POWAH" connector (5-36V).
1. The dynamixel servos has to be powered through the "DYNAMIXEL POWAH". This power rail only feeds the
   servos via the dynamixel bus. Please refer to the servo manuals regarding min/max voltages. Exceeding
   the rated servo voltages will permanently damage your servos.
