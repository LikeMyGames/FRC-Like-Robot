#ifndef MOTOR_H_
#define MOTOR_H_

#include <motor_driver.h>
#include <foc_math.h>
#include <unordered_map>
#include <Arduino.h>
#include <encoder.h>

#define MOTOR_PWM_DUTY_CYCLE ((uint32_t)20) // try 50 //in microseconds
#define MOTOR_TORQUE_LOOP_TIME 20           // try 100 // in microseconds
#define MOTOR_VELOCITY_LOOP_TIME 100        // try 80, 250 // in microseconds
#define MOTOR_POSITION_LOOP_TIME 1000       // try 320, 1000 // in microseconds
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

typedef enum
{
    MOTOR_DISABLED = 0,
    MOTOR_ENABLED,
} motor_status;

typedef enum
{
    NO_MOTOR_ERROR = 0,
    CAN_TIMEOUT,
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

static Motor *ref;

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
    motor_driver_t *driver;
    Pid *pos_pid;
    Pid *vel_pid;
    Pid *torque_pid;
    motor_running_mode running_mode;
    float Kt;
    volatile float target_pos;
    float cur_pos;
    volatile float target_vel;
    float cur_vel;
    volatile float target_torque;
    float cur_torque;
    float cur_current;
    uint32_t torque_loop_time;
    uint32_t velocity_loop_time;
    uint32_t position_loop_time;

    Motor();

    void Init(motor_driver_t *driver);

    void Disable();
    void Disable(motor_error error);

    void Enable();

    void Update();
    void SetRunningMode(motor_running_mode newMode);
    void ClearFault();
    void Error(motor_error error);

    void ReadCurrents();

    void DrivePhasesByPercent(float dA, float dB, float dC);
    void DrivePhasesByPercentFOC();
    void LoadParameterChanges(std::unordered_map<int, std::vector<uint8_t>> changes);

private:
    void TorqueLoop();
    void VelocityLoop();
    void PositionLoop();
};

void threadDriveMotor();

// void disableMotor(unsigned int id);
// void disableMotor(unsigned int id, motor_error error);
void disableMotor();
void disableMotor(motor_error error);

// void enableMotor(unsigned int id);
// void enableAllMotors();

// void MotorLoop(int MOTOR_ID);

namespace Motor_ns
{
    static Motor *ref;
    bool changes_to_parameters;
    std::unordered_map<int, std::vector<uint8_t>> changes;
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