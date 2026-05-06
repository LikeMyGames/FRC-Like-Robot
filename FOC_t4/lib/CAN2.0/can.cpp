#include <can.h>
#include <FlexCAN_T4.h>
#include <EEPROM.h>
#include <TeensyThreads.h>
#include <util_address.h>
#include <motor.h>
#include <string>
#include <util_math.h>
#include <cstring>

// Can Controller Config
#define NUM_TX_MAILBOXES 2
#define NUM_RX_MAILBOXES 6

// Can Id masks
#define API_CLASS_MASK 0b0000000000000111111000000000
#define API_INDEX_MASK 0b0000000000000000000111100000

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

typedef struct
{
    std::string name;
    std::function<std::vector<uint8_t>()> read;
    std::function<bool(uint8_t[])> write;
} CAN_REGISTER_t;

FlexCAN_T4<CAN3, RX_SIZE_256, TX_SIZE_16> Can0;
std::unordered_map<uint8_t, CAN_REGISTER_t *> REG_MAP;
uint32_t CAN_ID = 4;
// int LOOP_THREAD_ID;
int TIMEOUT_THREAD_ID;
volatile uint16_t TIMEOUT_MS_COUNT = 0;
volatile static uint16_t MAX_TIMEOUT_MS = 100;
elapsedMicros elapsed_since_sync;

/// @brief
void initCan()
{
    CAN_ID = EEPROM.read(CAN_ID_ADDRESS);
    MAX_TIMEOUT_MS = EEPROM.read(CAN_MAX_TIMEOUT_MS_ADDRESS);
    int deviceType = 0b10;
    int manufacturer = 0b11111111;
    int deviceId = CAN_ID;
    int canRangeLow = ((deviceType << 8 | manufacturer) << 15) | (deviceId & 0b11111);
    int canRangeHigh = ((deviceType << 8 | manufacturer) << 15) | 0b1111111111 << 5 | (deviceId & 0b11111);
    // delay(400);
    Can0.begin();
    Can0.setBaudRate(1 * 1000 * 1000);
    Can0.setMaxMB(NUM_TX_MAILBOXES + NUM_RX_MAILBOXES);
    for (int i = 0; i < NUM_RX_MAILBOXES; i++)
    {
        Can0.setMB((FLEXCAN_MAILBOX)i, RX, EXT);
    }
    for (int i = NUM_RX_MAILBOXES; i < (NUM_TX_MAILBOXES + NUM_RX_MAILBOXES); i++)
    {
        Can0.setMB((FLEXCAN_MAILBOX)i, TX, EXT);
    }

    Can0.setMBFilter(REJECT_ALL);
    Can0.enableMBInterrupts();
    Can0.onReceive(MB0, canSniff);
    Can0.onReceive(MB1, ProcessClassIndexMessage);
    Can0.onReceive(MB2, ProcessBroadcastMessage);
    Can0.onReceive(MB3, ProcessHeartbeatMessage);
    // Can0.setMBUserFilter(MB0, 0x00, 0x00);
    Can0.setMBFilterRange(MB1, canRangeLow, canRangeHigh);
    Can0.setMBFilterRange(MB2, 0b0, 0b111100000);
    Can0.setMBUserFilter(MB3, 0x01011840, 0xff);
    Can0.enhanceFilter(MB0);
    Can0.enhanceFilter(MB1);
    Can0.enhanceFilter(MB2);
    Can0.enhanceFilter(MB3);
    Can0.mailboxStatus();
    Can0.distribute();
    Serial.println(getID());
    // TIMEOUT_THREAD_ID = threads.addThread(checkTimeout);
    // Serial.println(LOOP_THREAD_ID);
    Serial.println(TIMEOUT_THREAD_ID);
}

std::vector<uint8_t> CanEchoRead()
{
    return (std::vector<uint8_t>){0, 0, 0, 0, 0, 0, 0};
}

bool CanEchoWrite(uint8_t args[])
{
    return true;
}

/// @brief
/// @param msg
void canSniff(const CAN_message_t &msg)
{
    digitalToggle(LED_BUILTIN);
    Serial.print("MB ");
    Serial.print(msg.mb);
    Serial.print("  OVERRUN: ");
    Serial.print(msg.flags.overrun);
    Serial.print("  LEN: ");
    Serial.print(msg.len);
    Serial.print(" EXT: ");
    Serial.print(msg.flags.extended);
    Serial.print(" TS: ");
    Serial.print(msg.timestamp);
    Serial.print(" ID: ");
    Serial.print(msg.id, HEX);
    Serial.print(" Buffer: ");
    for (uint8_t i = 0; i < msg.len; i++)
    {
        Serial.print(msg.buf[i], HEX);
        Serial.print(" ");
    }
    Serial.println();
}

void ProcessClassIndexMessage(const CAN_message_t &msg)
{
    Serial.printf("id: %b", msg.id);
    int apiClass = (msg.id & API_CLASS_MASK) >> 9;
    int apiIndex = (msg.id & API_INDEX_MASK) >> 5;

    Serial.printf("API Class: %v", apiClass);
    Serial.printf("API Index: %v", apiIndex);

    switch (apiClass)
    {
    case INFORMATION_API_CLASS:
        switch (apiIndex)
        {
        case 0:
            SendVersionMessage(msg.id);
            break;
        }
        break;
    case SLOT_0_POSITION_PID_API_CLASS:
    case SLOT_1_POSITION_PID_API_CLASS:
    case SLOT_2_POSITION_PID_API_CLASS:
    case SLOT_3_POSITION_PID_API_CLASS:
        int slotNum = apiClass == SLOT_0_POSITION_PID_API_CLASS   ? 0
                      : apiClass == SLOT_1_POSITION_PID_API_CLASS ? 1
                      : apiClass == SLOT_2_POSITION_PID_API_CLASS ? 2
                      : apiClass == SLOT_3_POSITION_PID_API_CLASS ? 3
                                                                  : -1;
        if (slotNum == -1)
        {
            return;
        }

        switch (apiIndex)
        {
        case 0:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerP(slotNum, 0, utils_bytes_to_float(data));
            break;
        case 1:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerI(slotNum, 0, utils_bytes_to_float(data));
            break;
        case 2:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerD(slotNum, 0, utils_bytes_to_float(data));
            break;
        case 3:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerIZone(slotNum, 0, utils_bytes_to_float(data));
            break;
        case 4:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerFF(slotNum, 0, utils_bytes_to_float(data));
            break;
        }
        break;
    case SLOT_0_VELOCITY_PID_API_CLASS:
    case SLOT_1_VELOCITY_PID_API_CLASS:
    case SLOT_2_VELOCITY_PID_API_CLASS:
    case SLOT_3_VELOCITY_PID_API_CLASS:
        int slotNum = apiClass == SLOT_0_VELOCITY_PID_API_CLASS   ? 0
                      : apiClass == SLOT_1_VELOCITY_PID_API_CLASS ? 1
                      : apiClass == SLOT_2_VELOCITY_PID_API_CLASS ? 2
                      : apiClass == SLOT_3_VELOCITY_PID_API_CLASS ? 3
                                                                  : -1;

        if (slotNum == -1)
        {
            return;
        }

        switch (apiIndex)
        {
        case 0:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerP(slotNum, 1, utils_bytes_to_float(data));
            break;
        case 1:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerI(slotNum, 1, utils_bytes_to_float(data));
            break;
        case 2:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerD(slotNum, 1, utils_bytes_to_float(data));
            break;
        case 3:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerIZone(slotNum, 1, utils_bytes_to_float(data));
            break;
        case 4:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerFF(slotNum, 1, utils_bytes_to_float(data));
            break;
        }
        break;
    case SLOT_0_TORQUE_PID_API_CLASS:
    case SLOT_1_TORQUE_PID_API_CLASS:
    case SLOT_2_TORQUE_PID_API_CLASS:
    case SLOT_3_TORQUE_PID_API_CLASS:
        int slotNum = apiClass == SLOT_0_TORQUE_PID_API_CLASS   ? 0
                      : apiClass == SLOT_1_TORQUE_PID_API_CLASS ? 1
                      : apiClass == SLOT_2_TORQUE_PID_API_CLASS ? 2
                      : apiClass == SLOT_3_TORQUE_PID_API_CLASS ? 3
                                                                : -1;

        if (slotNum == -1)
        {
            return;
        }

        switch (apiIndex)
        {
        case 0:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerP(slotNum, 2, utils_bytes_to_float(data));
            break;
        case 1:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerI(slotNum, 2, utils_bytes_to_float(data));
            break;
        case 2:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerD(slotNum, 2, utils_bytes_to_float(data));
            break;
        case 3:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerIZone(slotNum, 2, utils_bytes_to_float(data));
            break;
        case 4:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerFF(slotNum, 2, utils_bytes_to_float(data));
            break;
        }
        break;
    case READ_API_CLASS:
        switch (apiIndex)
        {
        case 0:
            Motor_ns::ref->internal_encoder->angle;
            break;
        case 1:
            Motor_ns::ref->external_encoder->angle;
            break;
        case 2:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetInternalEncoderOffset(utils_bytes_to_float(data));
            break;
        case 3:
            uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetExternalEncoderOffset(utils_bytes_to_float(data));
            break;
        case 4:
            Motor_ns::SetExternalEncoderType(msg.buf[0]);
            break;
        }
        break;
    }
}

void ProcessBroadcastMessage(const CAN_message_t &msg)
{
    int apiIndex = msg.id & API_INDEX_MASK;
    switch (apiIndex)
    {
    case BROADCAST_DISABLE:
        disableMotor();
        break;
    case BROADCAST_SYSTEM_HALT:
        // haltMotor();
        break;
    case BROADCAST_SYSTEM_RESET:
        break;
    case BROADCAST_DEVICE_ASSIGN:
        break;
    case BROADCAST_DEVICE_QUERY:
        break;
    case BROADCAST_HEARTBEAT:
        break;
    case BROADCAST_SYNC:
        break;
    case BROADCAST_UPDATE:
        break;
    case BROADCAST_FIRMWARE_VERSION:
        break;
    case BROADCAST_ENUMERATE:
        break;
    case BROADCAST_SYSTEM_RESUME:
        // unHaltMotor()
        break;
    }
}

void ProcessHeartbeatMessage(const CAN_message_t &msg)
{
    uint64_t value;
    std::memcpy(&value, msg.buf, sizeof(uint64_t));

    // int hours_mask = 0b11111;
    // int minutes_mask = 0b111111 << 5;
    // int seconds_mask = 0b111111 << 11;
    // int day_mask = 0b11111 << 17;
    // int month_mask = 0b1111 << 22;
    // int year_mask = 0b111111 << 26;
    // int tournament_type_mask = 0b111 << 32;
    // int system_watchdog_mask = 0b1 << 35;
    // int test_mode_mask = 0b1 << 36;
    // int autonomous_mode_mask = 0b1 << 37;
    int enabled_mask = 0b1 << 38;
    // int red_alliance_mask = 0b1 << 39;
    // int replay_number_mask = 0b111111 << 40;
    // int match_number_mask = 0b1111111111 << 46;
    // int match_time_mask = 0b11111111 << 56;

    bool enabled = (value & enabled_mask) >> 38;
    if (enabled)
    {
        enableMotor();
    }
    else
    {
        disableMotor();
    }

    SyncThreadToSyncPulse();
    elapsed_since_sync = 0;
}

void ProcessIdSpecificMessage(const CAN_message_t &msg)
{
    resetTimeout();
    Serial.println("Got ID specific messasge");

    if (msg.len == 0)
    {
        return;
    }

    CAN_message_t rmsg;
    rmsg.id = 0;
    rmsg.buf[0] = msg.buf[0];
    rmsg.flags.extended = true;
    Serial.println(rmsg.idhit);

    Serial.print("manipulating register: ");
    Serial.println(rmsg.buf[0]);
    Serial.print("message length: ");
    Serial.println(msg.len);

    if (msg.len == 1)
    {

        if (REG_MAP.count(msg.buf[0]) == 0)
        {
            Serial.print("No register exists with id ");
            Serial.println(msg.buf[0]);
        }

        auto reg = REG_MAP[msg.buf[0]];

        std::vector<uint8_t> buf = reg->read();
        if (buf.size() > 7)
        {
            rmsg.seq = true;
            uint size = buf.size();
            size %= 7;
            uint bytesRead = 0;
            rmsg.len = 8;
            while (bytesRead <= buf.size() - size)
            {
                for (int i = 0; i < 7; i++)
                {
                    rmsg.buf[i + 1] = buf[bytesRead + i];
                    bytesRead++;
                }

                Can0.write(rmsg);
            }

            if (size > 0)
            {
                rmsg.len = size + 1;
                for (uint i = 0; i < size; i++)
                {
                    rmsg.buf[i + 1] = buf.at(bytesRead + i);
                }
                Can0.write(rmsg);
            }
        }
        else
        {
            rmsg.len = buf.size() + 1;
            for (uint i = 1; i < buf.size(); i++)
            {
                rmsg.buf[i] = buf[i - 1];
            }
            Can0.write(rmsg);
        }
        return;
    }

    rmsg.len = 2;
    uint8_t funcBuf[7];
    for (uint i = 0; i < sizeof(funcBuf); i++)
    {
        funcBuf[i] = msg.buf[i + 1];
    }
    auto reg = REG_MAP[msg.buf[0]];
    rmsg.buf[1] = reg->write(funcBuf);

    Can0.write(rmsg);
}

void addCanRegister(uint8_t id, std::string name, std::function<std::vector<uint8_t>()> readAction, std::function<bool(uint8_t[])> writeAction)
{
    if (REG_MAP.count(id) != 0)
    {
        Serial.println("Could not add register, already exists");
    }

    CAN_REGISTER_t *reg = new CAN_REGISTER_t();
    Serial.println("created new can register");
    reg->name = name;
    reg->read = readAction;
    reg->write = writeAction;
    Serial.println("configured new can register");

    REG_MAP.insert_or_assign(id, reg);
    Serial.println("added new can register to map");
}

void SendVersionMessage(int id)
{
    CAN_message_t msg;
    msg.id = id;
    msg.flags.extended = true;
    msg.len = 1;
    msg.buf[0] = VERSION_MAJOR;
    msg.buf[1] = VERSION_MINOR;
    msg.buf[2] = VERSION_PATCH;

    Can0.write(msg);
}

void CanUpdate()
{
    Can0.events();
}

uint64_t TimeSinceLastSync()
{
    return elapsed_since_sync;
}

void SyncThreadToSyncPulse()
{
    int64_t timeDifference = elapsed_since_sync - EXPECTED_TIME_BETWEEN_SYNC;
    if (elapsed_since_sync > EXPECTED_TIME_BETWEEN_SYNC)
    {
        threads.delay_us(EXPECTED_TIME_BETWEEN_SYNC - elapsed_since_sync);
    }
    else if (elapsed_since_sync < EXPECTED_TIME_BETWEEN_SYNC)
    {
        threads.delay_us(elapsed_since_sync % EXPECTED_TIME_BETWEEN_SYNC);
    }
}

void Status()
{
    Can0.mailboxStatus();
}

uint32_t getID()
{
    return CAN_ID;
}

void setID(uint32_t newID)
{
    Can0.FLEXCAN_EnterFreezeMode();
    Can0.setMBUserFilter(MB3, newID, 0xff);
    Can0.FLEXCAN_ExitFreezeMode();
    CAN_ID = newID;
}

void resetTimeout()
{
    TIMEOUT_MS_COUNT = 0;
}

// void checkTimeout()
// {
//     while (1)
//     {
//         if (TIMEOUT_MS_COUNT >= MAX_TIMEOUT_MS)
//         {
//             disableMotor();
//             threads.delay(500);
//         }
//         TIMEOUT_MS_COUNT++;
//         threads.delay(1);
//     }
// }

uint16_t getMaxCanTimeout()
{
    return TIMEOUT_MS_COUNT;
}

void setMaxCanTimeout(uint16_t newTimeout)
{
    MAX_TIMEOUT_MS = newTimeout;
}