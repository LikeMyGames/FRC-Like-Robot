#include <can.h>
#include <FlexCAN_T4.h>
#include <EEPROM.h>
#include <TeensyThreads.h>
#include <util_address.h>
#include <motor.h>
#include <string>
#include <util_math.h>
#include <cstring>

typedef struct
{
    std::string name;
    std::function<std::vector<uint8_t>()> read;
    std::function<bool(uint8_t[])> write;
} CAN_REGISTER_t;

FlexCAN_T4<CAN2, RX_SIZE_256, TX_SIZE_16> Can0;
std::unordered_map<uint8_t, CAN_REGISTER_t *> REG_MAP;
uint32_t CAN_ID = 10;
// int LOOP_THREAD_ID;
int TIMEOUT_THREAD_ID;
volatile uint16_t TIMEOUT_MS_COUNT = 0;
volatile static uint16_t MAX_TIMEOUT_MS = 100;
elapsedMicros elapsed_since_sync;
uint32_t low_id;
uint32_t high_id;
uint32_t heartbeat_id = 0x01011840;
uint32_t broadcast_low_id = 0b00000000000000000000000000000;
uint32_t broadcast_high_id = 0b00000000000000000001111000000;

void initCan()
{
    MAX_TIMEOUT_MS = EEPROM.read(CAN_MAX_TIMEOUT_MS_ADDRESS);
    int deviceType = 0b10;
    int manufacturer = 0b11111111;
    // int deviceId = EEPROM.read(CAN_ID_ADDRESS);
    // CAN_ID = deviceId > 63 ? 0 : deviceId;
    low_id = (((deviceType << 8 | manufacturer) << 6 | 0) << 4 | 0) << 6 | CAN_ID;
    high_id = (((deviceType << 8 | manufacturer) << 6 | 0b111111) << 4 | 0b1111) << 6 | CAN_ID;

    Serial.printf("Loaded with Device ID: %d\n", CAN_ID);
    Serial.printf("Loaded with Low Can ID: %d\n", low_id);
    Serial.printf("Loaded with High Can ID: %d\n", high_id);
    Can0.begin();
    Can0.setBaudRate(1000000);
    // Can0.setMaxMB(NUM_TX_MAILBOXES + NUM_RX_MAILBOXES);
    // for (int i = 0; i < NUM_RX_MAILBOXES; i++)
    // {
    //     Can0.setMB((FLEXCAN_MAILBOX)i, RX, EXT);
    // }
    // for (int i = NUM_RX_MAILBOXES; i < (NUM_TX_MAILBOXES + NUM_RX_MAILBOXES); i++)
    // {
    //     Can0.setMB((FLEXCAN_MAILBOX)i, TX, EXT);
    // }

    // Can0.setMBFilter(REJECT_ALL);
    // Can0.enableMBInterrupts();
    // Can0.onReceive(MB0, canSniff);
    // Can0.onReceive(MB1, ProcessClassIndexMessage);
    // Can0.onReceive(MB2, ProcessBroadcastMessage);
    // Can0.onReceive(MB3, ProcessHeartbeatMessage);
    // Can0.setMBUserFilter(MB0, 0x00, 0x00);
    // Can0.setMBFilterRange(MB1, low_id, high_id);
    // Can0.setMBFilterRange(MB2, broadcast_low_id, broadcast_high_id);
    // Can0.setMBUserFilter(MB3, heartbeat_id, 0xff);
    // Can0.enhanceFilter(MB0);
    // Can0.enhanceFilter(MB1);
    // Can0.enhanceFilter(MB2);
    // Can0.enhanceFilter(MB3);
    Can0.mailboxStatus();
    // Can0.distribute();
    // TIMEOUT_THREAD_ID = threads.addThread(checkTimeout);
    // Serial.println(LOOP_THREAD_ID);
    // Serial.println(TIMEOUT_THREAD_ID);
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
    // digitalToggle(LED_BUILTIN);
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
    if (Serial.available())
    {
        Serial.println("Recieved ClassIndex Formated Message");
    }
    Serial.printf("id: %b", msg.id);
    int apiClass = (msg.id & API_CLASS_MASK) >> 9;
    int apiIndex = (msg.id & API_INDEX_MASK) >> 5;

    Serial.printf("API Class: %v", apiClass);
    Serial.printf("API Index: %v", apiIndex);

    int slotNum = -1;
    uint8_t data[4];

    switch (apiClass)
    {
    case INFORMATION_API_CLASS:
        switch (apiIndex)
        {
        case 0:
            SendVersionMessage();
            break;
        case 1:
            Motor_ns::PersistParameters();
            break;
        case 2:
            Motor_ns::ref->LoadParameterChanges();
        }
        break;
    case SLOT_0_POSITION_PID_API_CLASS:
    case SLOT_1_POSITION_PID_API_CLASS:
    case SLOT_2_POSITION_PID_API_CLASS:
    case SLOT_3_POSITION_PID_API_CLASS:
        slotNum = apiClass == SLOT_0_POSITION_PID_API_CLASS   ? 0
                  : apiClass == SLOT_1_POSITION_PID_API_CLASS ? 1
                  : apiClass == SLOT_2_POSITION_PID_API_CLASS ? 2
                  : apiClass == SLOT_3_POSITION_PID_API_CLASS ? 3
                                                              : -1;

        switch (apiIndex)
        {
        case 0:

            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerP(slotNum, 0, utils_bytes_to_float(data));
            break;
        case 1:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerI(slotNum, 0, utils_bytes_to_float(data));
            break;
        case 2:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerD(slotNum, 0, utils_bytes_to_float(data));
            break;
        case 3:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerIZone(slotNum, 0, utils_bytes_to_float(data));
            break;
        case 4:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerFF(slotNum, 0, utils_bytes_to_float(data));
            break;
        }
        break;
    case SLOT_0_VELOCITY_PID_API_CLASS:
    case SLOT_1_VELOCITY_PID_API_CLASS:
    case SLOT_2_VELOCITY_PID_API_CLASS:
    case SLOT_3_VELOCITY_PID_API_CLASS:
        slotNum = apiClass == SLOT_0_VELOCITY_PID_API_CLASS   ? 0
                  : apiClass == SLOT_1_VELOCITY_PID_API_CLASS ? 1
                  : apiClass == SLOT_2_VELOCITY_PID_API_CLASS ? 2
                  : apiClass == SLOT_3_VELOCITY_PID_API_CLASS ? 3
                                                              : -1;

        switch (apiIndex)
        {
        case 0:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerP(slotNum, 1, utils_bytes_to_float(data));
            break;
        case 1:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerI(slotNum, 1, utils_bytes_to_float(data));
            break;
        case 2:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerD(slotNum, 1, utils_bytes_to_float(data));
            break;
        case 3:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerIZone(slotNum, 1, utils_bytes_to_float(data));
            break;
        case 4:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerFF(slotNum, 1, utils_bytes_to_float(data));
            break;
        }
        break;
    case SLOT_0_TORQUE_PID_API_CLASS:
    case SLOT_1_TORQUE_PID_API_CLASS:
    case SLOT_2_TORQUE_PID_API_CLASS:
    case SLOT_3_TORQUE_PID_API_CLASS:
        slotNum = apiClass == SLOT_0_TORQUE_PID_API_CLASS   ? 0
                  : apiClass == SLOT_1_TORQUE_PID_API_CLASS ? 1
                  : apiClass == SLOT_2_TORQUE_PID_API_CLASS ? 2
                  : apiClass == SLOT_3_TORQUE_PID_API_CLASS ? 3
                                                            : -1;

        switch (apiIndex)
        {
        case 0:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerP(slotNum, 2, utils_bytes_to_float(data));
            break;
        case 1:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerI(slotNum, 2, utils_bytes_to_float(data));
            break;
        case 2:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerD(slotNum, 2, utils_bytes_to_float(data));
            break;
        case 3:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerIZone(slotNum, 2, utils_bytes_to_float(data));
            break;
        case 4:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetSlotControllerFF(slotNum, 2, utils_bytes_to_float(data));
            break;
        }
        break;
    case READ_API_CLASS:
        switch (apiIndex)
        {
        case 0:
            // Motor_ns::ref->internal_encoder->angle;
            break;
        case 1:
            // Motor_ns::ref->external_encoder->angle;
            break;
        case 2:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetInternalEncoderOffset(utils_bytes_to_float(data));
            break;
        case 3:
            // uint8_t data[4];
            std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(data));
            Motor_ns::SetExternalEncoderOffset(utils_bytes_to_float(data));
            break;
        case 4:
            Motor_ns::SetExternalEncoderType(msg.buf[0]);
            break;
        }
        break;

    case SETPOINT_API_CLASS:
        // float value;
        // int control_type;
        // int slot;

        uint8_t value_data[4];
        uint8_t control_type_data[1];
        uint8_t slot_data[1];

        std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(value_data));
        std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(control_type_data));
        std::copy(std::begin(msg.buf), std::begin(msg.buf) + 4, std::begin(slot_data));

        Motor_ns::SetSetpoint(utils_bytes_to_float(value_data), *control_type_data, *slot_data);
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
    // uint64_t tournament_type_mask = 0b111 << 32;
    uint64_t system_watchdog_mask = 0b1ULL << 35;
    // uint64_t test_mode_mask = 0b1ULL << 36;
    // uint64_t autonomous_mode_mask = 0b1ULL << 37;
    uint64_t enabled_mask = 0b1ULL << 38;
    // uint64_t red_alliance_mask = 0b1ULL << 39;
    // uint64_t replay_number_mask = 0b111111ULL << 40;
    // uint64_t match_number_mask = 0b1111111111ULL << 46;
    // uint64_t match_time_mask = 0b11111111ULL << 56;

    bool enabled = (value & enabled_mask) >> 38;
    bool system_watchdog = (value & system_watchdog_mask) >> 35;
    if (enabled && system_watchdog)
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

void SendVersionMessage()
{
    CAN_message_t msg;
    msg.id = CAN_ID;
    msg.flags.extended = true;
    msg.len = 1;
    msg.buf[0] = VERSION_MAJOR;
    msg.buf[1] = VERSION_MINOR;
    msg.buf[2] = VERSION_PATCH;

    Can0.write(msg);
    Serial.println("Sending version message");
}

void CanUpdate()
{
    CAN_message_t msg;

    while (Can0.read(msg))
    {
        Serial.printf("Recieved CAN frame: id=%d class=%d index=%d \n", msg.id, msg.id & 0b1111110000000000, msg.id & 1111000000);
        digitalWriteFast(LED_BUILTIN, HIGH);
        if (low_id < msg.id && msg.id < high_id)
        {
            ProcessClassIndexMessage(msg);
        }
        else if (broadcast_low_id < msg.id && msg.id < broadcast_high_id)
        {
            ProcessBroadcastMessage(msg);
        }
        else if (msg.id == heartbeat_id)
        {
            ProcessHeartbeatMessage(msg);
        }
        digitalWriteFast(LED_BUILTIN, LOW);
    }

    SendMotorEncoderData();
}

uint64_t TimeSinceLastSync()
{
    return elapsed_since_sync;
}

void SyncThreadToSyncPulse()
{
    // int64_t timeDifference = elapsed_since_sync - EXPECTED_TIME_BETWEEN_SYNC;
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

void SendMotorEncoderData()
{
    // Serial.println(Motor_ns::MotorInternalPosition());

    CAN_message_t msg;

    msg.id = BuildId(13, 0);
    msg.flags.extended = true;
    int32_t pos = Motor_ns::MotorInternalPosition();
    std::memcpy(msg.buf, &pos, sizeof(int32_t));
    msg.len = sizeof(pos);
    Can0.write(msg);

    msg.id = BuildId(13, 1);
    int32_t vel = Motor_ns::MotorInternalVelocity();
    std::memcpy(msg.buf, &vel, sizeof(int32_t));
    msg.len = sizeof(vel);
    Can0.write(msg);

    msg.id = BuildId(13, 2);
    int32_t accel = Motor_ns::MotorInternalAcceleration();
    std::memcpy(msg.buf, &accel, sizeof(int32_t));
    msg.len = sizeof(accel);
    Can0.write(msg);

    msg.id = BuildId(13, 3);
    pos = Motor_ns::MotorExternalPosition();
    std::memcpy(msg.buf, &pos, sizeof(int32_t));
    msg.len = sizeof(pos);
    Can0.write(msg);

    msg.id = BuildId(13, 4);
    vel = Motor_ns::MotorExternalVelocity();
    std::memcpy(msg.buf, &vel, sizeof(int32_t));
    msg.len = sizeof(vel);
    Can0.write(msg);

    msg.id = BuildId(13, 5);
    accel = Motor_ns::MotorExternalAccerlation();
    std::memcpy(msg.buf, &accel, sizeof(int32_t));
    msg.len = sizeof(accel);
    Can0.write(msg);
}

uint32_t BuildId(int class_var, int index)
{
    return low_id | ((class_var << 4 | index) << 7);
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