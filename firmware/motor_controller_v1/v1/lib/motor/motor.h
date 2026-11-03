#ifndef MOTOR_H_
#define MOTOR_H_

#include "foc_math.h"
#include <string>
#include <map>
#include <array>
#include <Arduino.h>
#include <Encoder.h>

#define MOTOR_PWM_DUTY_CYCLE ((uint32_t)50) // try 50 //in microseconds (us)
#define MOTOR_TORQUE_LOOP_TIME 100          // try 100 // in microseconds (us) // 20
#define MOTOR_VELOCITY_LOOP_TIME 500        // try 80, 250 // in microseconds (us) // 100
#define MOTOR_POSITION_LOOP_TIME 1000       // try 320, 1000 // in microseconds (us) // 1000
// #define MOTOR_INTERNAL_ENCODER 21
// #define MOTOR_EXTERNAL_ENCODER 20
#define MOTOR_PHASE_A_CURRENT_READ 23
#define MOTOR_PHASE_B_CURRENT_READ 22
#define MOTOR_PHASE_C_CURRENT_READ 21
#define MOTOR_INTERNAL_ENCODER_A 20
#define MOTOR_INTERNAL_ENCODER_B 19
#define MOTOR_INTERNAL_ENCODER_C 18
#define MOTOR_TEMP_READ 17
#define EXTERNAL_ENCODER_A 16
#define EXTERNAL_ENCODER_B 15
#define EXTERNAL_ENCODER_INDEX 14
#define EXTERNAL_ENCODER_ABSOLUTE 24

#define MOTOR_CURRENT_SENSE_RESISTANCE 0.001 // ohms
#define MOTOR_CURRENT_SENSE_GAIN 100         // V/V

#define MOTOR_INTERNAL_COUNTS_PER_REV 42

typedef enum
{
    MOTOR_DISABLED = 0,
    MOTOR_ENABLED,
} motor_status;

typedef enum
{
    NO_ERROR = 0,
    CAN_TIMEOUT,
    SETTING_SETPOINT,
} motor_error;

typedef struct
{
    float internal_by_external_ratio;
    uint internal_encoder_pin;
    uint external_encoder_pin;
    uint phaseA_current_read_pin;
    uint phaseB_current_read_pin;
    int PWM_FREQ;
} motor_config_t;

typedef struct
{
    float a, b, c;
} motor_phase_currents_t;

typedef struct
{
    uint activeSlot = 0;
    float Slot0_P, Slot1_P, Slot2_P, Slot3_P;
    float Slot0_I, Slot1_I, Slot2_I, Slot3_I;
    float Slot0_D, Slot1_D, Slot2_D, Slot3_D;
    float Slot0_IZone, Slot1_IZone, Slot2_IZone, Slot3_IZone;
    float Slot0_FF, Slot1_FF, Slot2_FF, Slot3_FF;

} motor_parameters_t;

// typedef struct
// {
//     motor_status status;
//     motor_error error;
//     foc_state_t foc_state;
//     foc_config_t foc_config;
//     motor_config_t config;
//     Encoder *internal_encoder;
//     Encoder *external_encoder;
//     int MOTOR_THREAD_ID = -1;
//     bool MOTOR_THREAD_STARTED = false;
//     motor_driver_t driver;
// } motor_state_t;

// static Motor *ref;

class Motor
{
public:
    motor_status status;
    motor_error error;
    Foc *foc;
    foc_config_t foc_config;
    motor_config_t config;
    Encoder *internal_encoder;
    Encoder *external_encoder;
    int MOTOR_THREAD_ID = -1;
    bool MOTOR_THREAD_STARTED = false;
    Pid *pos_pid;
    Pid *vel_pid;
    Pid *torque_pid;
    std::array<PidConstants, 12> pid_slots; // 0-2: slot 0, 3-5: slot 1, 6-8: slot 2, 9-11: slot 3;
    motor_running_mode running_mode;
    float Kt;
    volatile float target_pos;
    float cur_pos_internal;
    float cur_pos_external;
    volatile float target_vel;
    float cur_vel_internal;
    float cur_vel_external;
    volatile float target_torque;
    float cur_torque_internal;
    float cur_torque_external;
    float cur_current;
    uint32_t torque_loop_time;
    uint32_t velocity_loop_time;
    uint32_t position_loop_time;
    uint activeSlot = 0;
    int internal_encoder_a, internal_encoder_b, internal_encoder_c;

    Motor();
    // ~Motor();

    // void Init(motor_driver_t *driver);

    void Disable();
    void Disable(motor_error error);

    void Enable();

    void SetRunningMode(motor_running_mode newMode);
    void ClearFault();
    void Error(motor_error error);

    void ReadCurrents();

    // void DrivePhasesByPercent(float dA, float dB, float dC);
    // void DrivePhasesByPercentFOC();
    void LoadParameterChanges();

    void StartHardwareTimers();
    void StopHardwareTimers();

    // static void vSensorUpdateTimerCallback(TimerHandle_t xTimer);

    inline void TorqueLoop();
    inline void VelocityLoop();
    inline void PositionLoop();

    static void Update();

private:
    IntervalTimer update_loop;
    int update_loop_period = MOTOR_TORQUE_LOOP_TIME;

    // IntervalTimer posLoop;
    // IntervalTimer velLoop;
    // IntervalTimer torqueLoop;

    // void TorqueLoop();
    // void VelocityLoop();
    // void PositionLoop();
};

// void threadDriveMotor();

// void disableMotor(unsigned int id);
// void disableMotor(unsigned int id, motor_error error);
void disableMotor();
void disableMotor(motor_error error);

void enableMotor();
// void enableAllMotors();

// void MotorLoop(int MOTOR_ID);

namespace Motor_ns
{
    inline Motor *ref;

    inline int newRunningStatus;

    // bool changes_to_parameters;
    // motor_parameters_t changes;
    // std::unordered_map<std::string, std::vector<uint8_t>> changes;

    inline std::map<int, float> slot_pid_changes;
    // SemaphoreHandle_t xMotorRefSemaphore;

    // void SetSlotControllerControlType(int slotNum, int type, uint8_t type);

    // void TorqueLoop();
    // void VelocityLoop();
    // void PositionLoop();

    void Init();

    bool IsEnabled();

    int32_t MotorInternalPosition();
    float MotorInternalVelocity();
    float MotorInternalAcceleration();

    int32_t MotorExternalPosition();
    float MotorExternalVelocity();
    float MotorExternalAccerlation();

    void PersistParameters();

    void SetSetpoint(float value, int control_type, int slot);

    void SetSlotControllerP(int slotNum, int type, float p);
    void SetSlotControllerI(int slotNum, int type, float i);
    void SetSlotControllerD(int slotNum, int type, float d);
    void SetSlotControllerIZone(int slotNum, int type, float izone);
    void SetSlotControllerFF(int slotNum, int type, float ff);
    float GetInternalEncoderAngle();
    float GetExternalEncoderAngle();
    void SetInternalEncoderOffset(float offset);
    void SetExternalEncoderOffset(float offset);
    void SetExternalEncoderType(int type);

    // static std::unordered_map<uint, Motor *> map;
}

// static Motor *ref;
// static std::unordered_map<uint, Motor *> MOTOR_MAP;

// class Motor
// {
// public:
//     motor_state_t state;

//     Motor();

//     void disableMotor();
//     void disableMotor(motor_error error);
//     void enableMotor();
// };

#endif