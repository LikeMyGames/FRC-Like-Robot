#include <motor.h>
#include <Arduino.h>
#include <TeensyThreads.h>
#include <motor_driver.h>
#include <can.h>

static unsigned int NEXT_MOTOR_ID = 0;

Motor::Motor()
{
    this->position_loop_time = MOTOR_POSITION_LOOP_TIME;
    this->velocity_loop_time = MOTOR_VELOCITY_LOOP_TIME;
    this->torque_loop_time = MOTOR_TORQUE_LOOP_TIME;
    this->foc_config.Ts = MOTOR_PWM_DUTY_CYCLE / 1000000;
    // this->internal_encoder = new Encoder(MOTOR_INTERNAL_ENCODER);
    // this->external_encoder = new Encoder(MOTOR_EXTERNAL_ENCODER);
    this->foc = new Foc(this->foc_config);
    this->pos_pid = new Pid(&(this->cur_pos), 0.f, 0.f, 0.f, MOTOR_POSITION_LOOP_TIME / 1000000.f, -1.f, 1.f);
    this->vel_pid = new Pid(&(this->cur_vel), 0.f, 0.f, 0.f, MOTOR_VELOCITY_LOOP_TIME / 1000000.f, -1.f, 1.f);
    this->Init(&motor_driver_ns::motor_driver_1);
    Motor_ns::ref = this;
    // Motor_ns::map.insert_or_assign(0, this);
    // threads.addThread(threadDriveMotor);
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
        torque_pid->integralTerm = 0;
        torque_pid->derivativeTerm = 0;
        torque_pid->lastErr = 0;
        break;
    case CONTROL_MODE_TORQUE_BRAKE:
        torque_pid->integralTerm = 0;
        torque_pid->derivativeTerm = 0;
        torque_pid->lastErr = 0;
        break;
    default:
        return;
    }
    this->running_mode = newMode;
}

void Motor::Update()
{
    elapsedMicros elapsedLoopTime;
    if (this->error != NO_MOTOR_ERROR || this->status != MOTOR_ENABLED)
    {
        threads.delay_us(max(0, int(elapsedLoopTime - this->torque_loop_time)));
        return;
    }
    switch (running_mode)
    {
    case CONTROL_MODE_POS:
        this->PositionLoop();
        break;
    case CONTROL_MODE_POS_BRAKE:
        this->PositionLoop();
        break;
    case CONTROL_MODE_SPEED:
        this->VelocityLoop();
        break;
    case CONTROL_MODE_SPEED_BRAKE:
        this->VelocityLoop();
        break;
    case CONTROL_MODE_TORQUE:
        this->TorqueLoop();
        break;
    case CONTROL_MODE_TORQUE_BRAKE:
        this->TorqueLoop();
        break;
    default:
        return;
    }

    // foc->Drive(internal_encoder->ReadRad(), running_mode);
    // DriveMotorByPercent(driver, foc->dA, foc->dB, foc->dC);
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

void Motor::DrivePhasesByPercent(float dA, float dB, float dC)
{
    DriveMotorByPercent(driver, dA, dB, dC);
}

void Motor::DrivePhasesByPercentFOC()
{
    DriveMotorByPercent(driver, foc->dA, foc->dB, foc->dC);
}

void Motor::TorqueLoop()
{
    elapsedMicros elapsedTime;
    this->ReadCurrents();
    torque_pid->Update(target_torque);
    foc->iq_target = torque_pid->output;
    threads.delay_us(max(0, int(elapsedTime - this->torque_loop_time)));
    CanUpdate();
}

void Motor::VelocityLoop()
{
    elapsedMicros elapsedTime;
    vel_pid->Update(target_vel);
    while (elapsedTime <= this->velocity_loop_time)
    {
        TorqueLoop();
    }
}

void Motor::PositionLoop()
{
    elapsedMicros elapsedTime;
    pos_pid->Update(target_pos);
    while (elapsedTime <= this->position_loop_time)
    {
        VelocityLoop();
    }
}

void Motor::LoadParameterChanges(std::unordered_map<int, std::vector<uint8_t>> changes)
{
}

void threadDriveMotor()
{
    while (true)
    {
        Motor_ns::ref->Update();
        if (Motor_ns::changes_to_parameters)
        {
            Motor_ns::ref->LoadParameterChanges(Motor_ns::changes);
        }
    }
}

// void disableMotor(unsigned int id)
// {
//     MOTOR_STATE_MAP.at(id)->Disable();
// }

// void disableMotor(unsigned int id, motor_error error)
// {
//     MOTOR_STATE_MAP.at(id)->Disable(error);
// }

void disableMotor()
{
    Motor_ns::ref->Disable();
}

void disableMotor(motor_error error)
{
    Motor_ns::ref->Disable(error);
}

// void enableMotor(unsigned int id)
// {
//     MOTOR_STATE_MAP.at(id)->Enable();
// }

// void enableAllMotors()
// {
//     for (auto pair : MOTOR_STATE_MAP)
//     {
//         pair.second->Enable();
//     }
// }

// // Main motor loop
// void MotorLoop(int MOTOR_ID)
// {
//     MOTOR_STATE_MAP[MOTOR_ID]->Update();
// }
