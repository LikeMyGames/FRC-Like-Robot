#ifndef MOTOR_DRIVER_H_
#define MOTOR_DRIVER_H_

#include <eFlexPwm.h>

#define PWM_FREQ 20000

typedef struct
{
    // std::vector<int> submodule_numbers;
    std::vector<eFlex::SubModule> phaseA_submodules;
    std::vector<eFlex::SubModule> phaseB_submodules;
    std::vector<eFlex::SubModule> phaseC_submodules;
    eFlex::Timer timers[2];
} motor_driver_t;

// typedef struct {

// } module

// motor 1
static eFlex::SubModule Sm13(8, 7);
static eFlex::SubModule Sm20(4, 33);
static eFlex::SubModule Sm22(6, 9);

// motor 2
static eFlex::SubModule Sm31(29, 28);
static eFlex::SubModule Sm40(22);
static eFlex::SubModule Sm41(23);
static eFlex::SubModule Sm42(2, 3);

// motor 1 timers
static eFlex::Timer &Tm1 = Sm13.timer();
static eFlex::Timer &Tm2 = Sm20.timer();

// motor 2 timers
static eFlex::Timer &Tm3 = Sm31.timer();
static eFlex::Timer &Tm4 = Sm40.timer();

namespace motor_driver
{
    motor_driver_t motor1_instance = {
        // {13, 20, 22},
        {Sm13},
        {Sm20},
        {Sm22},
        {Tm1, Tm2},
    };

    motor_driver_t motor2_instance = {
        // {31, 40, 41, 42},
        {Sm31},
        {Sm40, Sm41},
        {Sm42},
        {Tm3, Tm4},
    };
    motor_driver_t *motor_1 = &motor1_instance;
    motor_driver_t *motor_2 = &motor2_instance;
}
// Functions
motor_driver_t NewMotorDriver();

// Duty times should be in percent of the cycle
void DriveMotorByPercent(motor_driver_t *driver, float dA, float dB, float dC);

// Duty times should be in nanoseconds
void DriveMotorByTime(motor_driver_t *driver, float dA, float dB, float dC);

#endif