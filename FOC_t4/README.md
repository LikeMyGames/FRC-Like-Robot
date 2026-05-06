# FOC_t4

This is a simple foc motor controller project designed to run with a PJRC Teensy 4.0 board. This project was designed to run on my custom designed printed circuit board and uses the FlexPWM modules of the Teensy 4.0 for accurate, fast switching of the main motor driving MOSFETS.

- [FOC\_t4](#foc_t4)
  - [Breakdown](#breakdown)
    - [Power In/Out, Phase Outputs](#power-inout-phase-outputs)
    - [Motor Driving MOSFETs](#motor-driving-mosfets)
    - [Gate Drivers](#gate-drivers)
    - [Current Sensing](#current-sensing)
    - [PJRC Teensy 4.0](#pjrc-teensy-40)
    - [CAN FD/2.0 Communication](#can-fd20-communication)
    - [Can Message Format](#can-message-format)
  - [CAN API Classes and Indecies](#can-api-classes-and-indecies)
    - [API Class Reference](#api-class-reference)
    - [API Index Reference](#api-index-reference)
      - [Slot 0 -\> Slot 3: (Classes 1-4)](#slot-0---slot-3-classes-1-4)
      - [Read: (Class 5)](#read-class-5)
    - [Full List](#full-list)
  - [FOC](#foc)

## Breakdown

The physical motor controller itself is comprised of six main parts.

- Power In/Out, Phase Outputs
- Motor Driving MOSFETs
- Gate Drivers
- Current Sensing
- PJRC Teensy 4.0
- CAN FD/2.0 Communication

### Power In/Out, Phase Outputs

This is the large lever locking connector block on the left side of the pcb (with capacitor rowing facing up and looking at the board from a top down perspective). This lever locking connector block is where the main power for the motor controller enters the board, along with where the phase wires of the motor are connected in order to be driven by the board

### Motor Driving MOSFETs

This the row of six rectangles that are arranged vertically on the board. These components switch power on and off to simulate different voltage going into the motor through a process called PWM. PWM stands for Pulse Width Modulation; this is a process where a digital signal (0 or 1) is quickly switched back and forth. For example: the state of the signal is 1 for the first 10 microseconds, but 0 for the other 900 microseconds. When looked on a macro instead of micro scale, this simulates a different voltage than the one going in. In the previous example, if the true voltage was 12V, then that "simulated" voltage would be 10% of 12V, or 1.2V.

### Gate Drivers

The gate drivers are the three chips that sit directly to the left of the MOSFETs. Due to the nature of how mosfets work, the require a relatively high voltage applied to their gate pin (the one that switches whether or not power flows through the MOSFET). Because of this, a gate driver is needed to conver the low voltage signals coming out of the Teensy 4.0 (3.3V) into something that can drive the MOSFETS (~22V).

### Current Sensing

This is needed to beable to properly run the [FOC control algorithm](foc). These are small, surface mount components placed in between some of the gate drivers. These components measure the voltage drop over a shunt resistor of a known resistance. This votlage drop is then amplified by some predetermined amount specific to the ship (either 50x or 100x in this case). This can then be read through an analog read on one of the Teensy 4.0 ADC pins. Ohm's law:

```math
I=\frac{V}{R}
```

to calculate the current that most be flowing through the resistor to get the particular voltage drop.

### PJRC Teensy 4.0

This is the main brain of the board. It runs the FOC control algorithm, interprets and sends out messages from and to the CAN communication system, and is able to connect to a gui to be configured (WIP).

### CAN FD/2.0 Communication

This is the communication protocal that the motor controller uses to talk with the rest of the project its being used in. For example, on a robot, it would be connected to the CAN bus that would be connected to the main processor so that it could receive commands. The main chip used on the board is made by Texas Instruments and supports both the CanFD spec and the Can2.0 spec. This means that the board is compatible with both CanFD and Can2.0 buses. The CAN interface runs at a speed of 1,000,000 bauds/second make sure the bus you connect the device to matches this speed. Possible feature addition of being able to change the CAN interface speed.

### Can Message Format

The CAN messages that this device is designed to handle are based off of the [CanFD and Can2.0](https://en.wikipedia.org/wiki/CAN_bus#Base_frame_format) spec. However, in the arbitration field of the frame, the device uses the extended version following the same structure as the [FIRST Robotics Competition](https://docs.wpilib.org/en/stable/docs/software/can-devices/can-addressing.html). An example of the arbitration portion of the frame is as follows.

|                     | Device Type | Manufacturer | API Class | Index  | Device ID |
| :-----------------: | :---------: | :----------: | :-------: | :----: | :-------: |
| **Value (Binary)**  |   0b00010   |  0b00000000  | 0b000000  | 0b0000 |  0b01010  |
| **Value (Decimal)** |      2      |     255      |     0     |   0    |    10     |
|    **Bit Width**    |      5      |      8       |     6     |   4    |     5     |

The example above refers to a motor controller that refers to with a CanId of 10 in the robot's main code. Depending on the API Class and Index values, different actions will be performed

## CAN API Classes and Indecies

### API Class Reference

| API Class |        Name         |                                Description                                |
| :-------: | :-----------------: | :-----------------------------------------------------------------------: |
|     0     |     Information     |                 Provides basic information of the device                  |
|     1     | Slot 0 Position PID | Sets the values stored in the ControlSlot.SLot0 position for position PID |
|     2     | Slot 0 Velocity PID | Sets the values stored in the ControlSlot.SLot0 position for velocity PID |
|     3     |  Slot 0 Torque PID  |  Sets the values stored in the ControlSlot.SLot0 position for torque PID  |
|     4     | Slot 1 Position PID | Sets the values stored in the ControlSlot.SLot0 position for position PID |
|     5     | Slot 1 Velocity PID | Sets the values stored in the ControlSlot.SLot0 position for velocity PID |
|     6     |  Slot 1 Torque PID  |  Sets the values stored in the ControlSlot.SLot0 position for torque PID  |
|     7     | Slot 2 Position PID | Sets the values stored in the ControlSlot.SLot0 position for position PID |
|     8     | Slot 2 Velocity PID | Sets the values stored in the ControlSlot.SLot0 position for velocity PID |
|     9     |  Slot 2 Torque PID  |  Sets the values stored in the ControlSlot.SLot0 position for torque PID  |
|    10     | Slot 3 Position PID | Sets the values stored in the ControlSlot.SLot0 position for position PID |
|    11     | Slot 3 Velocity PID | Sets the values stored in the ControlSlot.SLot0 position for velocity PID |
|    12     |  Slot 3 Torque PID  |  Sets the values stored in the ControlSlot.SLot0 position for torque PID  |
|    13     |        Read         |                   Reads various values from the device                    |

<!-- |     2     |   Slot 1    | Sets the values stored in the ControlSlot.Slot1 position |
|     3     |   Slot 2    | Sets the values stored in the ControlSlot.Slot2 position |
|     4     |   Slot 3    | Sets the values stored in the ControlSlot.Slot3 position | -->


<!-- |     6     |             |                                                          |
|     7     |             |                                                          | -->

### API Index Reference

#### Slot 0 -> Slot 3: (Classes 1-4)

| API Index |    Name    |                          Description                          |
| :-------: | :--------: | :-----------------------------------------------------------: |
|     0     |   Set P    |      Sets the P value in the PID controller of the slot       |
|     1     |   Set I    |      Sets the I value in the PID controller of the slot       |
|     2     |   Set D    |      Sets the D value in the PID Controller of the slot       |
|     3     | Set I Zone |    Sets the I Zone value in the PID Controller of hte slot    |
|     4     |   Set FF   | Sets the Feed Forward value in the PID Controller of the slot |

#### Read: (Class 5)

| API Index |           Name            |                                          Description                                          |
| :-------: | :-----------------------: | :-------------------------------------------------------------------------------------------: |
|     0     |   Read Internal Encoder   | Reads the value of the encoder connected to the internal encoder port on the motor controller |
|     1     |   Read External Encoder   | Reads the value of the econder connected to the external encoder port on the motor controller |
|     2     |   Set Internal Encoder    |      Sets the encoder to the specified position by applying an offset to the zero value       |
|     3     |   Set External Encoder    |      Sets the encoder to the specified position by applying an offset to the zero value       |
|     4     | Set External Encoder Type |   Sets the type of encoder plugged into the external encoder port (Absolute or Quadrature)    |

### Full List

|       Command Name        | API Class | API Index | Input / Return Types | Input / Return Meanings |
| :-----------------------: | :-------: | :-------: | :------------------: | :---------------------: |
|          Version          |     0     |     0     |   int8, int8, int8   |   major, minor, patch   |
|     Set Control Type      |   1 - 4   |     0     |         int8         |    ControlType value    |
|           Set P           |   1 - 4   |     1     |       float32        |         P value         |
|           Set I           |   1 - 4   |     2     |       float32        |         I value         |
|           Set D           |   1 - 4   |     3     |       float32        |         S value         |
|        Set I Zone         |   1 - 4   |     4     |       float32        |      I Zone value       |
|          Set FF           |   1 - 4   |     5     |       float32        |        FF value         |
|   Read Internal Encoder   |     5     |     0     |       float32        |    Encoder position     |
|   Read External Encoder   |     5     |     1     |       float32        |    Encoder position     |
|   Set Internal Encoder    |     5     |     2     |       float32        |    Encoder position     |
|   Set External Encoder    |     5     |     3     |       float32        |    Encoder position     |
| Set External Encoder Type |     5     |     4     |         int8         |      Encoder type       |

## FOC

A control algorithm to control the position, speed, or torque of a motor. It uses PID controllers to keep the q and d axises of a motor in constant orthagonality. This is the most efficient way to convert the power going into the motor controller into work down by the motor.

The algorithm analyses the current being applied to each phase of a motor.
Each phase is assigned an 1 out of 3 axises on a graph where all 3 axis are orthangonal to each other. This creates a graph where the current of each phase are on their own independent axises.

The algorithm then applies a transformation on to the 3 axis reference frame. in order to create a 2 axis reference frame.

Another transformation is then applied to this reference frame by an amount read from the internal encoder of the motor. This creates a 2 axis reference frame. One of the axis, D, cooresponds directly to the magnetic flux generated by the motor, while the other, Q, cooresponds directly to the torque generated by the motor.

The algorithm then uses a PID controller to bring the Q axis to a certain value. In order to get voltages to apply to motor phases, all of the transformations are done in reverse.

The fact that an encoder is used in the calculations of how to apply voltage to the motor makes this a closed loop control system. This means that the motor reacts to changes on the output.

FOC stands for Field Oriented Control. This is able to be seen directly by how it works. The algorithm focuses on orienting the torque and flux axises to be constantly orthagonal. This makes this algorithm one of the most efficient and quiet out of most motor control schemes.
