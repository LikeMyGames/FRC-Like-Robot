#include <motor_driver.h>
#include <cstdint>

// static unsigned int DRIVERS_CREATED = 0;
const float DeadTimeNs = 50.0;

namespace motor_driver_ns
{
    motor_driver_t motor_driver_1 = {
        {13},
        {20},
        {22},
        {Sm13},
        {Sm20},
        {Sm22},
        {Tm1, Tm2},
    };

    motor_driver_t motor_driver_2 = {
        {31},
        {40, 41},
        {42},
        {Sm31},
        {Sm40, Sm41},
        {Sm42},
        {Tm3, Tm4},
    };
}

void motor_driver_ns::InitMotorDriver(motor_driver_t *driver, int PWM_FREQ)
{
    eFlex::Config phaseA_config;
    eFlex::Config phaseB_config;
    eFlex::Config phaseC_config;

    // basic config for all phases
    phaseA_config.setReloadLogic(kPWM_ReloadPwmFullCycle);
    phaseA_config.setPwmFreqHz(PWM_FREQ);
    phaseB_config.setReloadLogic(kPWM_ReloadPwmFullCycle);
    phaseB_config.setPwmFreqHz(PWM_FREQ);
    phaseC_config.setReloadLogic(kPWM_ReloadPwmFullCycle);
    phaseC_config.setPwmFreqHz(PWM_FREQ);
    // myConfig.setClockSource(kPWM_Submodule0Clock);
    // myConfig.setInitializationControl(kPWM_Initialize_MasterSync);

    Serial.println("created basic eFlexPwm submodule configs");

    Serial.print("Phase A pwm submodule:\t");
    Serial.println(driver->phaseA_submodules.size());
    Serial.print("Number of Phase B pwm submodules:\t");
    Serial.println(driver->phaseB_submodules.size());
    Serial.print("Number of Phase C pwm submodules:\t");
    Serial.println(driver->phaseC_submodules.size());

    // Phase A config
    phaseA_config.setPairOperation(kPWM_ComplementaryPwmA);
    if (driver->phaseA_submodules[0].configure(phaseA_config) != true)
    {
        Serial.println("Phase A eFlexPwm Submodule initialization failed");
        exit(EXIT_FAILURE);
    }

    // Phase B config
    // Different because Phase B could be driven by two seperate pwm modules instead of one
    if (driver->phaseB_submodules.size() == 1)
    {
        phaseB_config.setPairOperation(kPWM_ComplementaryPwmA);
        if (driver->phaseB_submodules[0].configure(phaseB_config) != true)
        {
            Serial.println("Phase B eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
    }
    else
    {
        // submodule B1 config
        if (driver->phaseB_submodules[0].index() != 0)
        {
            phaseB_config.setClockSource(kPWM_Submodule0Clock);
        }
        else
        {
            phaseB_config.setClockSource(kPWM_BusClock);
        }
        phaseB_config.setPairOperation(kPWM_Independent);
        driver->phaseB_submodules[0].setupLevel(kPWM_HighTrue);
        if (driver->phaseB_submodules[0].configure(phaseB_config) != true)
        {
            Serial.println("Phase B eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }

        // submodule B2 config
        if (driver->phaseB_submodules[1].index() != 0)
        {
            phaseB_config.setClockSource(kPWM_Submodule0Clock);
        }
        else
        {
            phaseB_config.setClockSource(kPWM_BusClock);
        }
        driver->phaseB_submodules[1].setupLevel(kPWM_LowTrue);
        if (driver->phaseB_submodules[1].configure(phaseB_config) != true)
        {
            Serial.println("Phase B eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
    }

    // Phase C config
    if (driver->phaseC_submodules[0].index() != 0)
    {
        phaseC_config.setClockSource(kPWM_Submodule0Clock);
    }
    else
    {
        phaseC_config.setClockSource(kPWM_BusClock);
    }
    phaseC_config.setPairOperation(kPWM_ComplementaryPwmA);
    if (driver->phaseC_submodules[0].configure(phaseC_config) != true)
    {
        Serial.println("Phase C eFlexPwm Submodule initialization failed");
        exit(EXIT_FAILURE);
    }

    // Timer config
    uint16_t deadTimeVal = ((uint64_t)driver->timers[0].srcClockHz() * DeadTimeNs) / 1000000000;
    driver->timers[0].setupDeadtime(deadTimeVal);

    deadTimeVal = ((uint64_t)driver->timers[1].srcClockHz() * DeadTimeNs) / 1000000000;
    driver->timers[1].setupDeadtime(deadTimeVal);

    Serial.print("calculated deadTimeVal: ");
    Serial.println(deadTimeVal);

    // Start timers
    bool timer1Success = driver->timers[0].begin();
    bool timer2Success = driver->timers[1].begin();
    Serial.println("Started timers");
    Serial.println(timer1Success);
    Serial.println(timer2Success);

    if (!timer1Success || !timer2Success)
    {
        Serial.println("failed to start eFlexPwm submodule timers");
        exit(EXIT_FAILURE);
    }
    else
    {
        Serial.println("eFlexPwm submodules successfully started");
    }
}

void DriveMotorByPercent(motor_driver_t *driver, float dA, float dB, float dC)
{
    // Set timing for all phase A submodules
    for (auto submodule : driver->phaseA_submodules)
    {
        submodule.updateDutyCyclePercent(dA, eFlex::ChanA);
    }

    // Set timing for all phase B submodules
    for (auto submodule : driver->phaseB_submodules)
    {
        submodule.updateDutyCyclePercent(dB, eFlex::ChanA);
    }

    // Set timing for all phase C submodules
    for (auto submodule : driver->phaseC_submodules)
    {
        submodule.updateDutyCyclePercent(dC, eFlex::ChanA);
    }
}

// Not implemented
void DriveMotorByTime(motor_driver_t *driver, float dA, float dB, float dC)
{
}
