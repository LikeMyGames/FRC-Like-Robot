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

namespace motor_driver_ns
{
    extern motor_driver_t motor_driver_1;
    extern motor_driver_t motor_driver_2;

    void InitMotorDriver(motor_driver_t *driver);
}

// Functions

// Duty times should be in percent of the cycle
void DriveMotorByPercent(motor_driver_t *driver, float dA, float dB, float dC);

// Duty times should be in nanoseconds
void DriveMotorByTime(motor_driver_t *driver, float dA, float dB, float dC);

#endif