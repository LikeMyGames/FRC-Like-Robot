#include <motor.h>
#include <Arduino.h>
#include <TeensyThreads.h>
#include <motor_driver.h>
#include <can.h>
#include <util_math.h>
#include <string>
#include <analog.h>
#include <storage.h>

// #define SLOT_POSITION_PID_P_INDEX 0
// #define SLOT_VELOCITY_PID_P_INDEX 1

#define SLOT_0_POSITION_PID_INDEX 0
#define SLOT_0_VELOCITY_PID_INDEX 1
#define SLOT_0_TORQUE_PID_INDEX 2
#define SLOT_1_POSITION_PID_INDEX 3
#define SLOT_1_VELOCITY_PID_INDEX 4
#define SLOT_1_TORQUE_PID_INDEX 5
#define SLOT_2_POSITION_PID_INDEX 6
#define SLOT_2_VELOCITY_PID_INDEX 7
#define SLOT_2_TORQUE_PID_INDEX 8
#define SLOT_3_POSITION_PID_INDEX 9
#define SLOT_3_VELOCITY_PID_INDEX 10
#define SLOT_3_TORQUE_PID_INDEX 11

// static unsigned int NEXT_MOTOR_ID = 0;

Motor::Motor()
{
    pinMode(MOTOR_INTERNAL_ENCODER_A, INPUT_PULLDOWN);
    pinMode(MOTOR_INTERNAL_ENCODER_B, INPUT_PULLDOWN);
    pinMode(MOTOR_INTERNAL_ENCODER_C, INPUT_PULLDOWN);
    pinMode(EXTERNAL_ENCODER_A, INPUT_PULLDOWN);
    pinMode(EXTERNAL_ENCODER_B, INPUT_PULLDOWN);
    pinMode(EXTERNAL_ENCODER_INDEX, INPUT_PULLDOWN);
    pinMode(EXTERNAL_ENCODER_ABSOLUTE, INPUT_PULLDOWN);

    this->position_loop_time = MOTOR_POSITION_LOOP_TIME;
    this->velocity_loop_time = MOTOR_VELOCITY_LOOP_TIME;
    this->torque_loop_time = MOTOR_TORQUE_LOOP_TIME;
    this->foc_config.Ts = MOTOR_PWM_DUTY_CYCLE / 1000000;

    // this->internal_encoder = new QuadEncoder(1, MOTOR_INTERNAL_ENCODER_A, MOTOR_INTERNAL_ENCODER_B, 0, MOTOR_INTERNAL_ENCODER_C, 255, 255);
    // this->internal_encoder = new QuadEncoder(1, MOTOR_INTERNAL_ENCODER_A, MOTOR_INTERNAL_ENCODER_B, 0);
    // this->external_encoder = new QuadEncoder(2, EXTERNAL_ENCODER_A, EXTERNAL_ENCODER_B, 0, EXTERNAL_ENCODER_INDEX, 255, 255);

    this->internal_encoder = new Encoder(MOTOR_INTERNAL_ENCODER_A, MOTOR_INTERNAL_ENCODER_B);
    this->external_encoder = new Encoder(EXTERNAL_ENCODER_A, EXTERNAL_ENCODER_B);

    // this->internal_encoder->setInitConfig();
    // this->internal_encoder->EncConfig.enableReverseDirection = 1;
    // this->internal_encoder->EncConfig.revolutionCountCondition = 1;
    // this->internal_encoder->EncConfig.enableModuloCountMode = 1;
    // this->internal_encoder->EncConfig.positionModulusValue = 42;
    // this->internal_encoder->EncConfig.IndexTrigger = 1;
    // this->internal_encoder->EncConfig.INDEXTriggerMode = 1;
    // this->internal_encoder->init();

    // this->external_encoder->setInitConfig();
    // this->external_encoder->EncConfig.enableReverseDirection = 0;
    // this->external_encoder->EncConfig.revolutionCountCondition = 0;
    // this->external_encoder->EncConfig.enableModuloCountMode = 0;
    // this->external_encoder->EncConfig.IndexTrigger = 1;
    // this->external_encoder->EncConfig.INDEXTriggerMode = 1;
    // this->external_encoder->init();

    this->foc = new Foc(this->foc_config);
    this->pos_pid = new Pid(&(this->cur_pos_internal), 0.f, 0.f, 0.f, MOTOR_POSITION_LOOP_TIME / 1000000.f, -1.f, 1.f);
    this->vel_pid = new Pid(&(this->cur_vel_internal), 0.f, 0.f, 0.f, MOTOR_VELOCITY_LOOP_TIME / 1000000.f, -1.f, 1.f);

    this->driver = &motor_driver_ns::motor_driver_1;
    Serial.println("created motor driver");
    motor_driver_ns::InitMotorDriver(this->driver, this->config.PWM_FREQ);
    Motor_ns::ref = this;

    this->Disable();
}

// disabling the motor
void Motor::Disable()
{
    status = motor_status::MOTOR_DISABLED;
    this->StopHardwareTimers();
}

void Motor::Disable(motor_error error)
{
    this->Disable();
    this->error = error;
}

// enabling the motor
void Motor::Enable()
{
    if (error == motor_error::NO_ERROR)
    {
        status = motor_status::MOTOR_ENABLED;
        StartHardwareTimers();
    }
}

void Motor::SetRunningMode(motor_running_mode newMode)
{
    switch (running_mode)
    {
    case CONTROL_MODE_POS:
    case CONTROL_MODE_POS_BRAKE:
        pos_pid->integralTerm = 0;
        pos_pid->derivativeTerm = 0;
        pos_pid->lastErr = 0;
        // break;
    case CONTROL_MODE_SPEED:
    case CONTROL_MODE_SPEED_BRAKE:
        vel_pid->integralTerm = 0;
        vel_pid->derivativeTerm = 0;
        vel_pid->lastErr = 0;
        // break;
    case CONTROL_MODE_TORQUE:
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
    // elapsedMicros elapsedLoopTime;
    // if (this->error != NO_MOTOR_ERROR || this->status != MOTOR_ENABLED)
    // {
    //     threads.delay_us(max(0, int(elapsedLoopTime - this->torque_loop_time)));
    //     return;
    // }
    // switch (running_mode)
    // {
    // case CONTROL_MODE_POS:
    //     this->PositionLoop(this);
    //     break;
    // case CONTROL_MODE_POS_BRAKE:
    //     this->PositionLoop(this);
    //     break;
    // case CONTROL_MODE_SPEED:
    //     this->VelocityLoop(this);
    //     break;
    // case CONTROL_MODE_SPEED_BRAKE:
    //     this->VelocityLoop(this);
    //     break;
    // case CONTROL_MODE_TORQUE:
    //     this->TorqueLoop(this);
    //     break;
    // case CONTROL_MODE_TORQUE_BRAKE:
    //     this->TorqueLoop(this);
    //     break;
    // default:
    //     return;
    // }

    // foc->Drive(internal_encoder->ReadRad(), running_mode);
    // DriveMotorByPercent(driver, foc->dA, foc->dB, foc->dC);
}

void Motor::ClearFault()
{
    error = motor_error::NO_ERROR;
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

    // V=IR
    // voltage_drop / MOTOR_CURRENT_SENSE_GAIN = current * MOTOR_CURRENT_SENSE_RESISTANCE

    // this->foc->i_a = (analogRead(MOTOR_PHASE_A_CURRENT_READ) / MOTOR_CURRENT_SENSE_GAIN) / MOTOR_CURRENT_SENSE_RESISTANCE;
    // this->foc->i_b = (analogRead(MOTOR_PHASE_B_CURRENT_READ) / MOTOR_CURRENT_SENSE_GAIN) / MOTOR_CURRENT_SENSE_RESISTANCE;
    // this->foc->i_c = (analogRead(MOTOR_PHASE_C_CURRENT_READ) / MOTOR_CURRENT_SENSE_GAIN) / MOTOR_CURRENT_SENSE_RESISTANCE;

    this->foc->i_a = (ADC_ns::readPin(MOTOR_PHASE_A_CURRENT_READ) / MOTOR_CURRENT_SENSE_GAIN) / MOTOR_CURRENT_SENSE_RESISTANCE;
    this->foc->i_a = (ADC_ns::readPin(MOTOR_PHASE_B_CURRENT_READ) / MOTOR_CURRENT_SENSE_GAIN) / MOTOR_CURRENT_SENSE_RESISTANCE;
    this->foc->i_a = (ADC_ns::readPin(MOTOR_PHASE_C_CURRENT_READ) / MOTOR_CURRENT_SENSE_GAIN) / MOTOR_CURRENT_SENSE_RESISTANCE;
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
    Motor *inst = Motor_ns::ref;
    // inst->UpdateEncoderValues();
    if (inst == nullptr)
        return;
    if (inst->running_mode == CONTROL_MODE_TORQUE || inst->running_mode == CONTROL_MODE_SPEED || inst->running_mode == CONTROL_MODE_POS)
    {

        inst->ReadCurrents();
        inst->torque_pid->Update(inst->target_torque);
        inst->foc->iq_target = inst->torque_pid->output;
        inst->foc->Drive(inst->internal_encoder->read());
    }
}

void Motor::VelocityLoop()
{
    Motor *inst = Motor_ns::ref;
    if (inst == nullptr)
        return;
    if (inst->running_mode == CONTROL_MODE_SPEED || inst->running_mode == CONTROL_MODE_POS)
    {
        inst->vel_pid->Update(inst->target_vel);
        inst->target_torque = inst->vel_pid->output;
    }
}

void Motor::PositionLoop()
{
    Motor *inst = Motor_ns::ref;
    if (inst == nullptr)
        return;
    if (inst->running_mode == CONTROL_MODE_POS)
    {
        inst->pos_pid->Update(inst->target_pos);
        inst->target_vel = inst->pos_pid->output;
    }
}

void Motor::StartHardwareTimers()
{
    this->posLoop.begin(PositionLoop, MOTOR_POSITION_LOOP_TIME);
    this->velLoop.begin(VelocityLoop, MOTOR_VELOCITY_LOOP_TIME);
    this->torqueLoop.begin(TorqueLoop, MOTOR_TORQUE_LOOP_TIME);
}

void Motor::StopHardwareTimers()
{
    this->posLoop.end();
    this->velLoop.end();
    this->torqueLoop.end();
}

void Motor::UpdateEncoderValues()
{
    // int32_t old_pos_internal = this->cur_pos_internal;
    // int32_t old_vel_internal = this->cur_vel_internal;
    // int32_t old_pos_external = this->cur_pos_external;
    // int32_t old_vel_external = this->cur_vel_external;

    // this->cur_pos_internal = this->internal_encoder->read();
    // this->cur_vel_internal;
    // this->cur_pos_internal;
    // this->cur_torque_external;
    // this->cur_vel_external;
    // this->cur_torque_external;
}

void Motor::LoadParameterChanges()
{
    this->StopHardwareTimers();

    for (auto pair : Motor_ns::slot_pid_changes)
    {
        int var = pair.first / 12; // 0: p, 1: i, 2: d, 3: izone, 4: ff;
        int mod = pair.first % 12;
        // int type = mod % 3;
        // int slot = mod - type;

        switch (var)
        {
        case 0:
            this->pid_slots.at(mod).kp = pair.second;
            break;
        case 1:
            this->pid_slots.at(mod).ki = pair.second;
            break;
        case 2:
            this->pid_slots.at(mod).kd = pair.second;
            break;
        case 3:
            this->pid_slots.at(mod).izone = pair.second;
            break;
        case 4:
            this->pid_slots.at(mod).ff = pair.second;
            break;
        }
    }

    this->StartHardwareTimers();
}

// void threadDriveMotor()
// {
//     while (true)
//     {
//         Motor_ns::ref->Update();
//         if (Motor_ns::changes_to_parameters)
//         {
//             Motor_ns::ref->LoadParameterChanges();
//         }
//     }
// }

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

void enableMotor()
{
    Motor_ns::ref->Enable();
}

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

int32_t Motor_ns::MotorInternalPosition()
{
    return Motor_ns::ref->internal_encoder->read();
}

float Motor_ns::MotorInternalVelocity()
{
    return Motor_ns::ref->cur_vel_internal;
}

float Motor_ns::MotorInternalAcceleration()
{
    return Motor_ns::ref->cur_torque_internal;
}

float Motor_ns::MotorExternalPosition()
{
    return Motor_ns::ref->cur_pos_external;
}

float Motor_ns::MotorExternalVelocity()
{
    return Motor_ns::ref->cur_vel_external;
}

float Motor_ns::MotorExternalAccerlation()
{
    return Motor_ns::ref->cur_torque_external;
}

void Motor_ns::PersistParameters()
{
    Motor *inst = Motor_ns::ref;
    if (inst == nullptr)
    {
        return;
    }
    if (inst->status != motor_status::MOTOR_DISABLED)
    {
        return;
    }
    Storage::Update(SLOT_0_POSITION_P_EEPROM_INDEX, inst->pid_slots[0].kp);
    Storage::Update(SLOT_0_POSITION_I_EEPROM_INDEX, inst->pid_slots[0].ki);
    Storage::Update(SLOT_0_POSITION_D_EEPROM_INDEX, inst->pid_slots[0].kd);
    Storage::Update(SLOT_0_POSITION_IZONE_EEPROM_INDEX, inst->pid_slots[0].izone);
    Storage::Update(SLOT_0_POSITION_FF_EEPROM_INDEX, inst->pid_slots[0].ff);
    Storage::Update(SLOT_0_VELOCITY_P_EEPROM_INDEX, inst->pid_slots[1].kp);
    Storage::Update(SLOT_0_VELOCITY_I_EEPROM_INDEX, inst->pid_slots[1].ki);
    Storage::Update(SLOT_0_VELOCITY_D_EEPROM_INDEX, inst->pid_slots[1].kd);
    Storage::Update(SLOT_0_VELOCITY_IZONE_EEPROM_INDEX, inst->pid_slots[1].izone);
    Storage::Update(SLOT_0_TORQUE_FF_EEPROM_INDEX, inst->pid_slots[1].ff);
    Storage::Update(SLOT_0_TORQUE_P_EEPROM_INDEX, inst->pid_slots[2].kp);
    Storage::Update(SLOT_0_TORQUE_I_EEPROM_INDEX, inst->pid_slots[2].ki);
    Storage::Update(SLOT_0_TORQUE_D_EEPROM_INDEX, inst->pid_slots[2].kd);
    Storage::Update(SLOT_0_TORQUE_IZONE_EEPROM_INDEX, inst->pid_slots[2].izone);
    Storage::Update(SLOT_0_TORQUE_FF_EEPROM_INDEX, inst->pid_slots[2].ff);

    Storage::Update(SLOT_1_POSITION_P_EEPROM_INDEX, inst->pid_slots[3].kp);
    Storage::Update(SLOT_1_POSITION_I_EEPROM_INDEX, inst->pid_slots[3].ki);
    Storage::Update(SLOT_1_POSITION_D_EEPROM_INDEX, inst->pid_slots[3].kd);
    Storage::Update(SLOT_1_POSITION_IZONE_EEPROM_INDEX, inst->pid_slots[3].izone);
    Storage::Update(SLOT_1_POSITION_FF_EEPROM_INDEX, inst->pid_slots[3].ff);
    Storage::Update(SLOT_1_VELOCITY_P_EEPROM_INDEX, inst->pid_slots[4].kp);
    Storage::Update(SLOT_1_VELOCITY_I_EEPROM_INDEX, inst->pid_slots[4].ki);
    Storage::Update(SLOT_1_VELOCITY_D_EEPROM_INDEX, inst->pid_slots[4].kd);
    Storage::Update(SLOT_1_VELOCITY_IZONE_EEPROM_INDEX, inst->pid_slots[4].izone);
    Storage::Update(SLOT_1_TORQUE_FF_EEPROM_INDEX, inst->pid_slots[4].ff);
    Storage::Update(SLOT_1_TORQUE_P_EEPROM_INDEX, inst->pid_slots[5].kp);
    Storage::Update(SLOT_1_TORQUE_I_EEPROM_INDEX, inst->pid_slots[5].ki);
    Storage::Update(SLOT_1_TORQUE_D_EEPROM_INDEX, inst->pid_slots[5].kd);
    Storage::Update(SLOT_1_TORQUE_IZONE_EEPROM_INDEX, inst->pid_slots[5].izone);
    Storage::Update(SLOT_1_TORQUE_FF_EEPROM_INDEX, inst->pid_slots[5].ff);

    Storage::Update(SLOT_2_POSITION_P_EEPROM_INDEX, inst->pid_slots[6].kp);
    Storage::Update(SLOT_2_POSITION_I_EEPROM_INDEX, inst->pid_slots[6].ki);
    Storage::Update(SLOT_2_POSITION_D_EEPROM_INDEX, inst->pid_slots[6].kd);
    Storage::Update(SLOT_2_POSITION_IZONE_EEPROM_INDEX, inst->pid_slots[6].izone);
    Storage::Update(SLOT_2_POSITION_FF_EEPROM_INDEX, inst->pid_slots[6].ff);
    Storage::Update(SLOT_2_VELOCITY_P_EEPROM_INDEX, inst->pid_slots[7].kp);
    Storage::Update(SLOT_2_VELOCITY_I_EEPROM_INDEX, inst->pid_slots[7].ki);
    Storage::Update(SLOT_2_VELOCITY_D_EEPROM_INDEX, inst->pid_slots[7].kd);
    Storage::Update(SLOT_2_VELOCITY_IZONE_EEPROM_INDEX, inst->pid_slots[7].izone);
    Storage::Update(SLOT_2_TORQUE_FF_EEPROM_INDEX, inst->pid_slots[7].ff);
    Storage::Update(SLOT_2_TORQUE_P_EEPROM_INDEX, inst->pid_slots[8].kp);
    Storage::Update(SLOT_2_TORQUE_I_EEPROM_INDEX, inst->pid_slots[8].ki);
    Storage::Update(SLOT_2_TORQUE_D_EEPROM_INDEX, inst->pid_slots[8].kd);
    Storage::Update(SLOT_2_TORQUE_IZONE_EEPROM_INDEX, inst->pid_slots[8].izone);
    Storage::Update(SLOT_2_TORQUE_FF_EEPROM_INDEX, inst->pid_slots[8].ff);

    Storage::Update(SLOT_3_POSITION_P_EEPROM_INDEX, inst->pid_slots[9].kp);
    Storage::Update(SLOT_3_POSITION_I_EEPROM_INDEX, inst->pid_slots[9].ki);
    Storage::Update(SLOT_3_POSITION_D_EEPROM_INDEX, inst->pid_slots[9].kd);
    Storage::Update(SLOT_3_POSITION_IZONE_EEPROM_INDEX, inst->pid_slots[9].izone);
    Storage::Update(SLOT_3_POSITION_FF_EEPROM_INDEX, inst->pid_slots[9].ff);
    Storage::Update(SLOT_3_VELOCITY_P_EEPROM_INDEX, inst->pid_slots[10].kp);
    Storage::Update(SLOT_3_VELOCITY_I_EEPROM_INDEX, inst->pid_slots[10].ki);
    Storage::Update(SLOT_3_VELOCITY_D_EEPROM_INDEX, inst->pid_slots[10].kd);
    Storage::Update(SLOT_3_VELOCITY_IZONE_EEPROM_INDEX, inst->pid_slots[10].izone);
    Storage::Update(SLOT_3_TORQUE_FF_EEPROM_INDEX, inst->pid_slots[10].ff);
    Storage::Update(SLOT_3_TORQUE_P_EEPROM_INDEX, inst->pid_slots[11].kp);
    Storage::Update(SLOT_3_TORQUE_I_EEPROM_INDEX, inst->pid_slots[11].ki);
    Storage::Update(SLOT_3_TORQUE_D_EEPROM_INDEX, inst->pid_slots[11].kd);
    Storage::Update(SLOT_3_TORQUE_IZONE_EEPROM_INDEX, inst->pid_slots[11].izone);
    Storage::Update(SLOT_3_TORQUE_FF_EEPROM_INDEX, inst->pid_slots[11].ff);
}

void Motor_ns::SetSetpoint(float value, int control_type, int slot)
{
    Motor_ns::ref->StopHardwareTimers();
    Motor_ns::ref->activeSlot = slot;
    Motor_ns::ref->SetRunningMode(static_cast<motor_running_mode>(control_type));
    // Motor_ns::ref->running_mode = static_cast<motor_running_mode>(control_type);
    switch (control_type)
    {
    case CONTROL_MODE_POS:
    case CONTROL_MODE_POS_BRAKE:
        Motor_ns::ref->target_pos = value;
        break;
    case CONTROL_MODE_SPEED:
    case CONTROL_MODE_SPEED_BRAKE:
        Motor_ns::ref->target_vel = value;
        break;
    case CONTROL_MODE_TORQUE:
    case CONTROL_MODE_TORQUE_BRAKE:
        Motor_ns::ref->target_torque = value;
        break;
    }
    Motor_ns::ref->StartHardwareTimers();
}

void Motor_ns::SetSlotControllerP(int slotNum, int type, float p)
{
    int index = -1;

    switch (type)
    {
    case 0:
        index = slotNum == 0   ? SLOT_0_POSITION_PID_INDEX
                : slotNum == 1 ? SLOT_1_POSITION_PID_INDEX
                : slotNum == 2 ? SLOT_2_POSITION_PID_INDEX
                : slotNum == 3 ? SLOT_3_POSITION_PID_INDEX
                               : -1;
    case 1:
        index = slotNum == 0   ? SLOT_0_VELOCITY_PID_INDEX
                : slotNum == 1 ? SLOT_1_VELOCITY_PID_INDEX
                : slotNum == 2 ? SLOT_2_VELOCITY_PID_INDEX
                : slotNum == 3 ? SLOT_3_VELOCITY_PID_INDEX
                               : -1;
    case 2:
        index = slotNum == 0   ? SLOT_0_TORQUE_PID_INDEX
                : slotNum == 1 ? SLOT_1_TORQUE_PID_INDEX
                : slotNum == 2 ? SLOT_2_TORQUE_PID_INDEX
                : slotNum == 3 ? SLOT_3_TORQUE_PID_INDEX
                               : -1;
    }

    if (index == -1)
    {
        return;
    }

    Motor_ns::slot_pid_changes[index] = p;
}

void Motor_ns::SetSlotControllerI(int slotNum, int type, float i)
{
    int index = -1;

    switch (type)
    {
    case 0:
        index = slotNum == 0   ? SLOT_0_POSITION_PID_INDEX
                : slotNum == 1 ? SLOT_1_POSITION_PID_INDEX
                : slotNum == 2 ? SLOT_2_POSITION_PID_INDEX
                : slotNum == 3 ? SLOT_3_POSITION_PID_INDEX
                               : -1;
    case 1:
        index = slotNum == 0   ? SLOT_0_VELOCITY_PID_INDEX
                : slotNum == 1 ? SLOT_1_VELOCITY_PID_INDEX
                : slotNum == 2 ? SLOT_2_VELOCITY_PID_INDEX
                : slotNum == 3 ? SLOT_3_VELOCITY_PID_INDEX
                               : -1;
    case 2:
        index = slotNum == 0   ? SLOT_0_TORQUE_PID_INDEX
                : slotNum == 1 ? SLOT_1_TORQUE_PID_INDEX
                : slotNum == 2 ? SLOT_2_TORQUE_PID_INDEX
                : slotNum == 3 ? SLOT_3_TORQUE_PID_INDEX
                               : -1;
    }

    if (index == -1)
    {
        return;
    }

    Motor_ns::slot_pid_changes[index + 12] = i;
}

void Motor_ns::SetSlotControllerD(int slotNum, int type, float d)
{
    int index = -1;

    switch (type)
    {
    case 0:
        index = slotNum == 0   ? SLOT_0_POSITION_PID_INDEX
                : slotNum == 1 ? SLOT_1_POSITION_PID_INDEX
                : slotNum == 2 ? SLOT_2_POSITION_PID_INDEX
                : slotNum == 3 ? SLOT_3_POSITION_PID_INDEX
                               : -1;
    case 1:
        index = slotNum == 0   ? SLOT_0_VELOCITY_PID_INDEX
                : slotNum == 1 ? SLOT_1_VELOCITY_PID_INDEX
                : slotNum == 2 ? SLOT_2_VELOCITY_PID_INDEX
                : slotNum == 3 ? SLOT_3_VELOCITY_PID_INDEX
                               : -1;
    case 2:
        index = slotNum == 0   ? SLOT_0_TORQUE_PID_INDEX
                : slotNum == 1 ? SLOT_1_TORQUE_PID_INDEX
                : slotNum == 2 ? SLOT_2_TORQUE_PID_INDEX
                : slotNum == 3 ? SLOT_3_TORQUE_PID_INDEX
                               : -1;
    }

    if (index == -1)
    {
        return;
    }

    Motor_ns::slot_pid_changes[index + 24] = d;
}

void Motor_ns::SetSlotControllerIZone(int slotNum, int type, float izone)
{
    int index = -1;

    switch (type)
    {
    case 0:
        index = slotNum == 0   ? SLOT_0_POSITION_PID_INDEX
                : slotNum == 1 ? SLOT_1_POSITION_PID_INDEX
                : slotNum == 2 ? SLOT_2_POSITION_PID_INDEX
                : slotNum == 3 ? SLOT_3_POSITION_PID_INDEX
                               : -1;
    case 1:
        index = slotNum == 0   ? SLOT_0_VELOCITY_PID_INDEX
                : slotNum == 1 ? SLOT_1_VELOCITY_PID_INDEX
                : slotNum == 2 ? SLOT_2_VELOCITY_PID_INDEX
                : slotNum == 3 ? SLOT_3_VELOCITY_PID_INDEX
                               : -1;
    case 2:
        index = slotNum == 0   ? SLOT_0_TORQUE_PID_INDEX
                : slotNum == 1 ? SLOT_1_TORQUE_PID_INDEX
                : slotNum == 2 ? SLOT_2_TORQUE_PID_INDEX
                : slotNum == 3 ? SLOT_3_TORQUE_PID_INDEX
                               : -1;
    }

    if (index == -1)
    {
        return;
    }

    Motor_ns::slot_pid_changes[index + 36] = izone;
}

void Motor_ns::SetSlotControllerFF(int slotNum, int type, float ff)
{
    int index = -1;

    switch (type)
    {
    case 0:
        index = slotNum == 0   ? SLOT_0_POSITION_PID_INDEX
                : slotNum == 1 ? SLOT_1_POSITION_PID_INDEX
                : slotNum == 2 ? SLOT_2_POSITION_PID_INDEX
                : slotNum == 3 ? SLOT_3_POSITION_PID_INDEX
                               : -1;
    case 1:
        index = slotNum == 0   ? SLOT_0_VELOCITY_PID_INDEX
                : slotNum == 1 ? SLOT_1_VELOCITY_PID_INDEX
                : slotNum == 2 ? SLOT_2_VELOCITY_PID_INDEX
                : slotNum == 3 ? SLOT_3_VELOCITY_PID_INDEX
                               : -1;
    case 2:
        index = slotNum == 0   ? SLOT_0_TORQUE_PID_INDEX
                : slotNum == 1 ? SLOT_1_TORQUE_PID_INDEX
                : slotNum == 2 ? SLOT_2_TORQUE_PID_INDEX
                : slotNum == 3 ? SLOT_3_TORQUE_PID_INDEX
                               : -1;
    }

    if (index == -1)
    {
        return;
    }

    Motor_ns::slot_pid_changes[index + 48] = ff;
}

float Motor_ns::GetInternalEncoderAngle()
{
    return (float)(ref->internal_encoder->read()) / MOTOR_INTERNAL_COUNTS_PER_REV;
    // return Encoder_ns::countsToRad(ref->internal_encoder->read());
}

float Motor_ns::GetExternalEncoderAngle()
{
    return (float)(ref->external_encoder->read());
    // return Encoder_ns::countsToRad(ref->external_encoder->read());
}

void Motor_ns::SetInternalEncoderOffset(float offset)
{
}

void Motor_ns::SetExternalEncoderOffset(float offset)
{
}

void Motor_ns::SetExternalEncoderType(int type)
{
}
