#include <motor_driver.h>
#include <cstdint>

static unsigned int DRIVERS_CREATED = 0;
const float DeadTimeNs = 50.0;

motor_driver_t *NewMotorDriver()
{
    motor_driver_t *driver;
    switch (DRIVERS_CREATED)
    {
    case 0:
        driver = &*motor_driver::motor_1;
        break;
    case 1:
        driver = &*motor_driver::motor_2;
        break;
    default:
        Serial.println("MOTOR_NUM_LIMIT reached");
        break;
    }
    DRIVERS_CREATED++;

    eFlex::Config myConfig;
    myConfig.setReloadLogic(kPWM_ReloadPwmFullCycle);
    myConfig.setClockSource(kPWM_Submodule0Clock);
    myConfig.setPrescale(kPWM_Prescale_Divide_1);
    myConfig.setInitializationControl(kPWM_Initialize_MasterSync);
    myConfig.setPwmFreqHz(PWM_FREQ);

    // SubModule config
    if (driver->phaseA_submodules.size() == 1)
    {
        myConfig.setPairOperation(kPWM_ComplementaryPwmA);
        if (driver->phaseA_submodules[0].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
    }
    else
    {
        myConfig.setPairOperation(kPWM_Independent);
        driver->phaseA_submodules[0].setupLevel(kPWM_HighTrue);
        if (driver->phaseA_submodules[0].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
        driver->phaseA_submodules[1].setupLevel(kPWM_LowTrue);
        if (driver->phaseA_submodules[1].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
    }

    if (driver->phaseB_submodules.size() == 1)
    {
        myConfig.setPairOperation(kPWM_ComplementaryPwmA);
        if (driver->phaseB_submodules[0].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
    }
    else
    {
        myConfig.setPairOperation(kPWM_Independent);
        driver->phaseB_submodules[0].setupLevel(kPWM_HighTrue);
        if (driver->phaseB_submodules[0].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
        driver->phaseB_submodules[1].setupLevel(kPWM_LowTrue);
        if (driver->phaseB_submodules[1].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
    }

    if (driver->phaseC_submodules.size() == 1)
    {
        myConfig.setPairOperation(kPWM_ComplementaryPwmA);
        if (driver->phaseC_submodules[0].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
    }
    else
    {
        myConfig.setPairOperation(kPWM_Independent);
        driver->phaseC_submodules[0].setupLevel(kPWM_HighTrue);
        if (driver->phaseC_submodules[0].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
        driver->phaseC_submodules[1].setupLevel(kPWM_LowTrue);
        if (driver->phaseC_submodules[1].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
    }

    // Timer config
    uint16_t deadTimeVal = ((uint64_t)driver->timers[0].srcClockHz() * DeadTimeNs) / 1000000000;
    driver->timers[0].setupDeadtime(deadTimeVal);

    deadTimeVal = ((uint64_t)driver->timers[1].srcClockHz() * DeadTimeNs) / 1000000000;
    driver->timers[1].setupDeadtime(deadTimeVal);

    // Start timers
    bool timer1Success = driver->timers[0].begin();
    bool timer2Success = driver->timers[1].begin();

    if (timer1Success || timer2Success)
    {
        Serial.println("failed to start eFlexPwm submodules");
        exit(EXIT_FAILURE);
    }
    else
    {
        Serial.println("eFlexPwm submodules successfully started");
    }

    return driver;
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
