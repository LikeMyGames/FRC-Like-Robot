#include <motor.h>
#include <Arduino.h>
#include <TeensyThreads.h>
#include <motor_driver.h>

static unsigned int NEXT_MOTOR_ID = 0;

Motor::Motor(motor_config_t config)
{
    Serial.println("created new motor state");
    this->config = config;
    foc_config.Ts = 1000000 / config.PWM_FREQ;
    Serial.println("assigned config to new state");
    internal_encoder = new Encoder(config.internal_encoder_pin);
    Serial.println("created internal encoder object");
    external_encoder = new Encoder(config.external_encoder_pin);
    Serial.println("created external encoder object");
    foc = new Foc(foc_config);
    Serial.println("init foc_state");
    MOTOR_STATE_MAP.insert_or_assign(NEXT_MOTOR_ID, this);
    Serial.println("added motor to motor map");
    pos_pid = new Pid(&(internal_encoder->angle), 0, 0, 0, 0);
    vel_pid = new Pid(&(internal_encoder->angle), 0, 0, 0, 0);
    NEXT_MOTOR_ID++;
    Serial.println("incremented NEXT_MOTOR_ID");
}

void Motor::Init(motor_driver_t *driver)
{
    this->driver = driver;
    Serial.println("created motor driver");
    motor_driver_ns::InitMotorDriver(driver, this->config.PWM_FREQ);
}

// disabling the motor
void Motor::Disable()
{
    status = MOTOR_DISABLED;
    if (MOTOR_THREAD_STARTED)
    {
        threads.suspend(MOTOR_THREAD_ID);
    }
}

void Motor::Disable(motor_error error)
{
    this->Disable();
    this->error = error;
}

// enabling the motor
void Motor::Enable()
{
    if (!error)
    {
        status = MOTOR_ENABLED;
        if (!MOTOR_THREAD_STARTED)
        {
            threads.restart(MOTOR_THREAD_ID);
            MOTOR_THREAD_STARTED = true;
        }
    }
}

void Motor::Update()
{
    foc->Drive(internal_encoder->ReadRad());
    DriveMotorByPercent(driver, foc->dA, foc->dB, foc->dC);
}

void disableMotor(unsigned int id)
{
    MOTOR_STATE_MAP.at(id)->Disable();
}

void disableMotor(unsigned int id, motor_error error)
{
    MOTOR_STATE_MAP.at(id)->Disable(error);
}

void disableAllMotors()
{
    Serial.println("Disabling all motors on controller");
    for (auto pair : MOTOR_STATE_MAP)
    {
        pair.second->Disable();
    }
}

void disableAllMotors(motor_error error)
{
    for (auto pair : MOTOR_STATE_MAP)
    {
        pair.second->Disable(error);
    }
}

void enableMotor(unsigned int id)
{
    MOTOR_STATE_MAP.at(id)->Enable();
}

void enableAllMotors()
{
    for (auto pair : MOTOR_STATE_MAP)
    {
        pair.second->Enable();
    }
}

// Main motor loop
void MotorLoop(int MOTOR_ID)
{
    MOTOR_STATE_MAP[MOTOR_ID]->Update();
}
