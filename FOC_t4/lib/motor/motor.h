#ifndef MOTOR_H_
#define MOTOR_H_

#include <motor_driver.h>
#include <foc_math.h>
#include <unordered_map>
#include <Arduino.h>
#include <encoder.h>

#define MOTOR_PWM_DUTY_CYCLE ((uint32_t)50)

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
} motor_config_t;

typedef struct
{
    float a, b, c;
} motor_phase_currents_t;

typedef struct
{
    motor_status status;
    motor_error error;
    foc_state_t *foc_state;
    foc_config_t *foc_config;
    motor_config_t *config;
    Encoder *internal_encoder;
    Encoder *external_encoder;
    int MOTOR_THREAD_ID = -1;
    bool MOTOR_THREAD_STARTED = false;
    motor_driver_t *driver;
} motor_state_t;

static std::unordered_map<unsigned int, motor_state_t *> MOTOR_STATE_MAP;

motor_state_t NewMotor(motor_config_t *config);

void disableMotor(motor_state_t *motor);
void disableMotor(motor_state_t *motor, motor_error error);
void disableMotor(unsigned int id);
void disableMotor(unsigned int id, motor_error error);
void disableAllMotors();
void disableAllMotors(motor_error error);
void enableMotor(motor_state_t *motor);
void enableMotor(unsigned int id);
void enableAllMotors();

void MotorLoop(int MOTOR_ID);
void MotorLoop(motor_state_t *state);

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