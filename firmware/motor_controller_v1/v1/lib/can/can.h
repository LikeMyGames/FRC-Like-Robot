#ifndef CAN_H_
#define CAN_H_

#include <Arduino.h>
#include <FlexCAN_T4.h>
#include <string>

// Can Controller Config
#define NUM_TX_MAILBOXES 2
#define NUM_RX_MAILBOXES 6

// Can Id masks
#define API_CLASS_MASK 0b1111110000000000
#define API_INDEX_MASK 0b1111000000
#define DEVICE_ID_MASK 0b111111

// API Classes
#define INFORMATION_API_CLASS 0
#define SLOT_0_POSITION_PID_API_CLASS 1
#define SLOT_0_VELOCITY_PID_API_CLASS 2
#define SLOT_0_TORQUE_PID_API_CLASS 3
#define SLOT_1_POSITION_PID_API_CLASS 4
#define SLOT_1_VELOCITY_PID_API_CLASS 5
#define SLOT_1_TORQUE_PID_API_CLASS 6
#define SLOT_2_POSITION_PID_API_CLASS 7
#define SLOT_2_VELOCITY_PID_API_CLASS 8
#define SLOT_2_TORQUE_PID_API_CLASS 9
#define SLOT_3_POSITION_PID_API_CLASS 10
#define SLOT_3_VELOCITY_PID_API_CLASS 11
#define SLOT_3_TORQUE_PID_API_CLASS 12
#define READ_API_CLASS 13
#define SETPOINT_API_CLASS 14

// Broadcast Indexes
#define BROADCAST_DISABLE 0          // Disable index
#define BROADCAST_SYSTEM_HALT 1      // System Halt index
#define BROADCAST_SYSTEM_RESET 2     // System Reset index
#define BROADCAST_DEVICE_ASSIGN 3    // Device Assign index
#define BROADCAST_DEVICE_QUERY 4     // Device Query index
#define BROADCAST_HEARTBEAT 5        // Heartbeat index
#define BROADCAST_SYNC 6             // Sync index
#define BROADCAST_UPDATE 7           // Update index
#define BROADCAST_FIRMWARE_VERSION 8 // Firmware Version index
#define BROADCAST_ENUMERATE 9        // Enumerate index
#define BROADCAST_SYSTEM_RESUME 10   // System Resume index

// Sync settings
#define MAIN_PROCESSOR_PERIOD_MS 10
#define MAIN_PROCESSOR_PERIOD_mS 10000
#define EXPECTED_TIME_BETWEEN_SYNC 10000

void initCan();

void addCanRegister(uint8_t id,
                    std::string name,
                    std::function<std::vector<uint8_t>()> readAction,
                    std::function<bool(uint8_t[])> writeAction);
uint32_t getID();
void resetTimeout();
// void checkTimeout();
void canSniff(const CAN_message_t &);

void ProcessClassIndexMessage(const CAN_message_t &msg);
void ProcessBroadcastMessage(const CAN_message_t &msg);
void ProcessHeartbeatMessage(const CAN_message_t &msg);
void ProcessIdSpecificMessage(const CAN_message_t &msg);

void SendVersionMessage();

void SyncThreadToSyncPulse();

void CanUpdate();
void Status();

void SendMotorEncoderData();
uint32_t BuildId(int class_var, int index);

extern std::vector<uint8_t> CanEchoRead();
extern bool CanEchoWrite(uint8_t args[]);

#endif