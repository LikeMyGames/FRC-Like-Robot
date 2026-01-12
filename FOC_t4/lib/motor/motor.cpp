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
    pos_pid = new Pid(&(cur_pos), 0, 0, 0, 1000000 / config.PWM_FREQ, 0, 0);
    vel_pid = new Pid(&(cur_vel), 0, 0, 0, 1000000 / config.PWM_FREQ, 0, 0);
    torque_pid = new Pid(&(cur_torque), 0, 0, 0, 1000000 / config.PWM_FREQ, 0, 0);
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
    // if (MOTOR_THREAD_STARTED)
    // {
    //     threads.suspend(MOTOR_THREAD_ID);
    // }
}

void Motor::Disable(motor_error error)
{
    this->Disable();
    this->error = error;
}

// enabling the motor
void Motor::Enable()
{
    if (error == NO_MOTOR_ERROR)
    {
        status = MOTOR_ENABLED;
        // if (!MOTOR_THREAD_STARTED)
        // {
        //     threads.restart(MOTOR_THREAD_ID);
        //     MOTOR_THREAD_STARTED = true;
        // }
    }
}

void Motor::SetRunningMode(motor_running_mode newMode)
{
    switch (running_mode)
    {
    case CONTROL_MODE_POS:
        pos_pid->integralTerm = 0;
        pos_pid->derivativeTerm = 0;
        pos_pid->lastErr = 0;
        break;
    case CONTROL_MODE_POS_BRAKE:
        pos_pid->integralTerm = 0;
        pos_pid->derivativeTerm = 0;
        pos_pid->lastErr = 0;
        break;
    case CONTROL_MODE_SPEED:
        vel_pid->integralTerm = 0;
        vel_pid->derivativeTerm = 0;
        vel_pid->lastErr = 0;
        break;
    case CONTROL_MODE_SPEED_BRAKE:
        vel_pid->integralTerm = 0;
        vel_pid->derivativeTerm = 0;
        vel_pid->lastErr = 0;
        break;
    case CONTROL_MODE_TORQUE:
        break;
    case CONTROL_MODE_TORQUE_BRAKE:
        break;
    default:
        return;
    }
    this->running_mode = newMode;
}

void Motor::Update()
{
    this->ReadCurrents();
    switch (running_mode)
    {
    case CONTROL_MODE_POS:
        pos_pid->Update(target_pos);
        vel_pid->Update(pos_pid->output);
        cur_torque = Kt * cur_current;
        torque_pid->Update(vel_pid->output);
        foc->iq_target = torque_pid->output;
        break;
    case CONTROL_MODE_POS_BRAKE:
        pos_pid->Update(target_pos);
        vel_pid->Update(pos_pid->output);
        torque_pid->Update(vel_pid->output);
        foc->iq_target = torque_pid->output;
        break;
    case CONTROL_MODE_SPEED:
        vel_pid->Update(target_vel);
        torque_pid->Update(vel_pid->output);
        foc->iq_target = torque_pid->output;
        break;
    case CONTROL_MODE_SPEED_BRAKE:
        vel_pid->Update(target_vel);
        torque_pid->Update(vel_pid->output);
        foc->iq_target = torque_pid->output;
        break;
    case CONTROL_MODE_TORQUE:
        torque_pid->Update(target_torque);
        foc->iq_target = torque_pid->output;
        break;
    case CONTROL_MODE_TORQUE_BRAKE:
        torque_pid->Update(target_torque);
        foc->iq_target = torque_pid->output;
        break;
    default:
        return;
    }

    foc->Drive(internal_encoder->ReadRad(), running_mode);
    DriveMotorByPercent(driver, foc->dA, foc->dB, foc->dC);
}

void Motor::ClearFault()
{
    error = NO_MOTOR_ERROR;
}

void Motor::Error(motor_error error)
{
    this->error = error;
}

void Motor::ReadCurrents()
{
    // need to implement current sensing
    // might be able to use hall sensors on rev motors to sense phase current
    // most likely going to use voltage divider on main incoming power to detect current
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
    // Serial.println("Disabling all motors on controller");
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
