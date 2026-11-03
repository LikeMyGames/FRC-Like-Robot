#ifndef MOTOR_DRIVER_H_
#define MOTOR_DRIVER_H_

#include <eFlexPwm.h>

namespace driver
{
    typedef enum
    {
        A = 0,
        B = 1,
        C = 2,
        U = 0,
        V = 1,
        W = 2
    } Phase;

    void Init(uint32_t pwm_freq_hz);

    void Enable();
    void Disable();

    void DriveByTime(float time_A, float time_B, float time_C);
    void DriveByPercent(float percent_A, float percent_B, float percent_C);

}

#endif /* MOTOR_DRIVER_H_ */