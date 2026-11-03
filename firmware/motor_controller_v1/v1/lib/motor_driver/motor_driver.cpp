#include "motor_driver.h"

namespace
{
    eFlex::SubModule Sm13(8, 7);  // A or U
    eFlex::SubModule Sm20(4, 33); // B or V
    eFlex::SubModule Sm22(6, 9);  // C or W
    eFlex::Timer &timer_A = Sm13.timer();
    eFlex::Timer &timer_B = Sm20.timer();
    eFlex::Timer &timer_C = Sm22.timer();
    int time_to_duty = 0;
} // namespace

void driver::Init(uint32_t pwm_freq_hz)
{
    eFlex::Config config;

    config.setReloadLogic(kPWM_ReloadPwmFullCycle);
    config.setClockSource(kPWM_BusClock);
    config.setPwmFreqHz(pwm_freq_hz);
    config.setPairOperation(kPWM_ComplementaryPwmA);
    if (Sm13.configure(config) != true)
    {
        Serial.println("Phase A eFlexPwm Submodule initialization failed");
        exit(EXIT_FAILURE);
    }
    Sm13.setupLevel(kPWM_HighTrue);
    Sm13.setPwmMode(kPWM_CenterAligned);

    if (Sm20.configure(config) != true)
    {
        Serial.println("Phase B eFlexPwm Submodule initialization failed");
        exit(EXIT_FAILURE);
    }
    Sm20.setupLevel(kPWM_HighTrue);
    Sm20.setPwmMode(kPWM_CenterAligned);

    if (Sm22.configure(config) != true)
    {
        Serial.println("Phase C eFlexPwm Submodule initialization failed");
        exit(EXIT_FAILURE);
    }
    Sm22.setupLevel(kPWM_HighTrue);
    Sm22.setPwmMode(kPWM_CenterAligned);

    timer_A = Sm13.timer();
    timer_B = Sm20.timer();
    timer_C = Sm22.timer();

    uint64_t deadtime_ns = 1000;
    uint64_t deadtime = 0;

    deadtime = ((uint64_t)timer_A.srcClockHz() * deadtime_ns) / 1000000000;
    Serial.printf("timer_A deadtime: %d", deadtime);
    timer_A.setupDeadtime(deadtime);
    timer_A.begin();
    timer_A.stop();

    deadtime = ((uint64_t)timer_B.srcClockHz() * deadtime_ns) / 1000000000;
    Serial.printf("timer_B deadtime: %d", deadtime);
    timer_B.setupDeadtime(deadtime);
    timer_B.begin();
    timer_B.stop();

    deadtime = ((uint64_t)timer_C.srcClockHz() * deadtime_ns) / 1000000000;
    Serial.printf("timer_C deadtime: %d", deadtime);
    timer_C.setupDeadtime(deadtime);
    timer_C.begin();
    timer_C.stop();
}

void driver::Enable()
{
    timer_A.start();
    timer_B.start();
    timer_C.start();
}

void driver::Disable()
{
    timer_A.stop();
    timer_B.stop();
    timer_C.stop();
}

void driver::DriveByTime(float time_A, float time_B, float time_C)
{
    Sm13.updateDutyCycle(time_A * time_to_duty, eFlex::ChanA);
    Sm20.updateDutyCycle(time_B * time_to_duty, eFlex::ChanA);
    Sm22.updateDutyCycle(time_C * time_to_duty, eFlex::ChanA);
}

void driver::DriveByPercent(float percent_A, float percent_B, float percent_C)
{
    Sm13.updateDutyCyclePercent(percent_A, eFlex::ChanA);
    Sm20.updateDutyCyclePercent(percent_B, eFlex::ChanA);
    Sm22.updateDutyCyclePercent(percent_C, eFlex::ChanA);
}
