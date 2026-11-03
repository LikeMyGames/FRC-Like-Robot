#include <Arduino.h>
#include <EEPROM.h>
#include <cstring>

#define CAN_DEVICE_ID 0

#define SLOT_0_POSITION_P_EEPROM_INDEX 1
#define SLOT_0_POSITION_I_EEPROM_INDEX 5
#define SLOT_0_POSITION_D_EEPROM_INDEX 9
#define SLOT_0_POSITION_IZONE_EEPROM_INDEX 13
#define SLOT_0_POSITION_FF_EEPROM_INDEX 17
#define SLOT_0_VELOCITY_P_EEPROM_INDEX 21
#define SLOT_0_VELOCITY_I_EEPROM_INDEX 25
#define SLOT_0_VELOCITY_D_EEPROM_INDEX 29
#define SLOT_0_VELOCITY_IZONE_EEPROM_INDEX 33
#define SLOT_0_VELOCITY_FF_EEPROM_INDEX 37
#define SLOT_0_TORQUE_P_EEPROM_INDEX 41
#define SLOT_0_TORQUE_I_EEPROM_INDEX 45
#define SLOT_0_TORQUE_D_EEPROM_INDEX 49
#define SLOT_0_TORQUE_IZONE_EEPROM_INDEX 53
#define SLOT_0_TORQUE_FF_EEPROM_INDEX 57

#define SLOT_1_POSITION_P_EEPROM_INDEX 1
#define SLOT_1_POSITION_I_EEPROM_INDEX 5
#define SLOT_1_POSITION_D_EEPROM_INDEX 9
#define SLOT_1_POSITION_IZONE_EEPROM_INDEX 13
#define SLOT_1_POSITION_FF_EEPROM_INDEX 17
#define SLOT_1_VELOCITY_P_EEPROM_INDEX 21
#define SLOT_1_VELOCITY_I_EEPROM_INDEX 25
#define SLOT_1_VELOCITY_D_EEPROM_INDEX 29
#define SLOT_1_VELOCITY_IZONE_EEPROM_INDEX 33
#define SLOT_1_VELOCITY_FF_EEPROM_INDEX 37
#define SLOT_1_TORQUE_P_EEPROM_INDEX 41
#define SLOT_1_TORQUE_I_EEPROM_INDEX 45
#define SLOT_1_TORQUE_D_EEPROM_INDEX 49
#define SLOT_1_TORQUE_IZONE_EEPROM_INDEX 53
#define SLOT_1_TORQUE_FF_EEPROM_INDEX 57

#define SLOT_2_POSITION_P_EEPROM_INDEX 61
#define SLOT_2_POSITION_I_EEPROM_INDEX 65
#define SLOT_2_POSITION_D_EEPROM_INDEX 69
#define SLOT_2_POSITION_IZONE_EEPROM_INDEX 73
#define SLOT_2_POSITION_FF_EEPROM_INDEX 77
#define SLOT_2_VELOCITY_P_EEPROM_INDEX 81
#define SLOT_2_VELOCITY_I_EEPROM_INDEX 85
#define SLOT_2_VELOCITY_D_EEPROM_INDEX 89
#define SLOT_2_VELOCITY_IZONE_EEPROM_INDEX 93
#define SLOT_2_VELOCITY_FF_EEPROM_INDEX 97
#define SLOT_2_TORQUE_P_EEPROM_INDEX 101
#define SLOT_2_TORQUE_I_EEPROM_INDEX 105
#define SLOT_2_TORQUE_D_EEPROM_INDEX 109
#define SLOT_2_TORQUE_IZONE_EEPROM_INDEX 113
#define SLOT_2_TORQUE_FF_EEPROM_INDEX 117

#define SLOT_3_POSITION_P_EEPROM_INDEX 121
#define SLOT_3_POSITION_I_EEPROM_INDEX 125
#define SLOT_3_POSITION_D_EEPROM_INDEX 129
#define SLOT_3_POSITION_IZONE_EEPROM_INDEX 133
#define SLOT_3_POSITION_FF_EEPROM_INDEX 137
#define SLOT_3_VELOCITY_P_EEPROM_INDEX 131
#define SLOT_3_VELOCITY_I_EEPROM_INDEX 145
#define SLOT_3_VELOCITY_D_EEPROM_INDEX 149
#define SLOT_3_VELOCITY_IZONE_EEPROM_INDEX 153
#define SLOT_3_VELOCITY_FF_EEPROM_INDEX 157
#define SLOT_3_TORQUE_P_EEPROM_INDEX 161
#define SLOT_3_TORQUE_I_EEPROM_INDEX 165
#define SLOT_3_TORQUE_D_EEPROM_INDEX 169
#define SLOT_3_TORQUE_IZONE_EEPROM_INDEX 173
#define SLOT_3_TORQUE_FF_EEPROM_INDEX 177

namespace Storage
{
    template <typename T>
    void Write(int address, const T value)
    {
        int size = sizeof(value);
        uint8_t bytes[size];
        std::memcpy(bytes, &value, size);

        for (int i = 0; i < size; i++)
        {
            EEPROM.write(address + i, bytes[i]);
        }
    }

    template <typename T>
    void Read(int address, const T &value)
    {
        int size = sizeof(*value);
        uint8_t bytes[size];

        for (int i = 0; i < size; i++)
        {
            bytes[i] = EEPROM.read(address + i);
        }

        std::memcpy(value, bytes, size);
    }

    // template <typename T>
    // void Put(int address, const T value);

    // template <typename T>
    // void Get(int address, const T &value);

    template <typename T>
    void Update(int address, const T value)
    {
        int size = sizeof(value);
        uint8_t bytes[size];
        std::memcpy(bytes, &value, size);

        for (int i = 0; i < size; i++)
        {
            EEPROM.update(address + i, bytes[i]);
        }
    }
}