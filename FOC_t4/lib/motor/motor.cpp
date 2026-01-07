#include <motor.h>
#include <Arduino.h>
#include <TeensyThreads.h>
#include <motor_driver.h>

static unsigned int NEXT_MOTOR_ID = 0;

motor_state_t *NewMotor(motor_config_t *config)
{
    motor_state_t *state = {};
    state->config = config;
    state->internal_encoder = new Encoder(config->internal_encoder_pin);
    state->external_encoder = new Encoder(config->external_encoder_pin);
    state->driver = NewMotorDriver();
    state->foc_state = init_foc(state->foc_config);
    MOTOR_STATE_MAP.insert_or_assign(NEXT_MOTOR_ID, state);
    NEXT_MOTOR_ID++;
    return state;
}

// disabling the motor
void disableMotor(motor_state_t *state)
{
    state->status = MOTOR_DISABLED;
    if (state->MOTOR_THREAD_STARTED)
    {
        threads.suspend(state->MOTOR_THREAD_ID);
    }
}

void disableMotor(motor_state_t *state, motor_error error)
{
    disableMotor(state);
    state->error = error;
}

void disableMotor(unsigned int id)
{
    disableMotor(MOTOR_STATE_MAP.at(id));
}

void disableMotor(unsigned int id, motor_error error)
{
    disableMotor(MOTOR_STATE_MAP.at(id), error);
}

void disableAllMotors()
{
    Serial.println("Disabling all motors on controller");
    for (auto pair : MOTOR_STATE_MAP)
    {
        disableMotor(pair.second);
    }
}

void disableAllMotors(motor_error error)
{
    for (auto pair : MOTOR_STATE_MAP)
    {
        disableMotor(pair.second);
    }
}

// enabling the motor
void enableMotor(motor_state_t *state)
{
    if (!state->error)
    {
        state->status = MOTOR_ENABLED;
        if (!state->MOTOR_THREAD_STARTED)
        {
            threads.restart(state->MOTOR_THREAD_ID);
            state->MOTOR_THREAD_STARTED = true;
        }
    }
}

void enableMotor(unsigned int id)
{
    enableMotor(MOTOR_STATE_MAP.at(id));
}

void enableAllMotors()
{
    for (auto pair : MOTOR_STATE_MAP)
    {
        enableMotor(pair.second);
    }
}

// Main motor loop
void MotorLoop(int MOTOR_ID)
{
    motor_state_t *state = MOTOR_STATE_MAP[MOTOR_ID];
    foc_drive(state->foc_state, state->internal_encoder->ReadRad());
    DriveMotorByPercent(&(state->driver), state->foc_state->dA, state->foc_state->dB, state->foc_state->dC);
}

void MotorLoop(motor_state_t *state)
{
    foc_drive(state->foc_state, state->internal_encoder->ReadRad());
    DriveMotorByPercent(&(state->driver), state->foc_state->dA, state->foc_state->dB, state->foc_state->dC);
}