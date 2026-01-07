#include <can.h>
#include <motor.h>
#include <util_math.h>

std::vector<uint8_t> Motor1_Iq_read();
bool Motor1_Iq_write(uint8_t buf[7]);
std::vector<uint8_t> Motor2_Iq_read();
bool Motor2_Iq_write(uint8_t buf[7]);
std::vector<uint8_t> Motor1_mode_read();
bool Motor1_mode_write(uint8_t buf[7]);
std::vector<uint8_t> Motor2_mode_read();
bool Motor2_mode_write(uint8_t buf[7]);
std::vector<uint8_t> Motor1_target_pos_read();
bool Motor1_target_pos_write(uint8_t buf[7]);
std::vector<uint8_t> Motor2_target_pos_read();
bool Motor2_target_pos_write(uint8_t buf[7]);
std::vector<uint8_t> Motor1_target_vel_read();
bool Motor1_target_vel_write(uint8_t buf[7]);
std::vector<uint8_t> Motor2_target_vel_read();
bool Motor2_target_vel_write(uint8_t buf[7]);

motor_state_t *Motor1;
static motor_config_t motor1_config = {
    1.f,
    21U,
    20U,
};

motor_state_t *Motor2;
static motor_config_t motor2_config = {
    1.f,
    19,
    18,
};

void setup()
{
    Serial.begin(115200);
    // Serial.setTimeout(0);
    Serial.println("Starting setup");
    initCan();
    Serial.println("finished motor controller setup");
    Motor1 = NewMotor(&motor1_config);
    Motor2 = NewMotor(&motor2_config);
    addCanRegister(0, "echo", CanEchoRead, CanEchoWrite);
    addCanRegister(1, "motor1 mode", Motor1_mode_read, Motor1_mode_write);
    addCanRegister(2, "motor2 mode", Motor2_mode_read, Motor2_mode_write);
    addCanRegister(3, "motor1 target iq", Motor1_Iq_read, Motor1_Iq_write);
    addCanRegister(4, "motor2 target iq", Motor2_Iq_read, Motor2_Iq_write);
    addCanRegister(5, "motor1 target pos", Motor1_target_pos_read, Motor1_target_pos_write);
    addCanRegister(6, "motor1 target vel", Motor1_target_vel_read, Motor1_target_vel_write);
    addCanRegister(7, "motor2 target pos", Motor2_target_pos_read, Motor2_target_pos_write);
    addCanRegister(8, "motor2 target vel", Motor2_target_vel_read, Motor2_target_vel_write);
}

std::vector<uint8_t> Motor1_Iq_read()
{
    return {(uint8_t)(Motor1->foc_state->iq_target)};
}

bool Motor1_Iq_write(uint8_t buf[7])
{
    Motor1->foc_state->iq_target = (float)buf[0];
    return true;
}

std::vector<uint8_t> Motor2_Iq_read()
{
    return {(uint8_t)(Motor2->foc_state->iq_target)};
}

bool Motor2_Iq_write(uint8_t buf[7])
{
    Motor2->foc_state->iq_target = (float)buf[0];
    return true;
}

std::vector<uint8_t> Motor1_mode_read()
{
    return {(uint8_t)(Motor1->foc_state->mode)};
}

bool Motor1_mode_write(uint8_t buf[7])
{
    Motor1->foc_state->mode = (foc_mode)buf[0];
    return true;
}

std::vector<uint8_t> Motor2_mode_read()
{
    return {(u_int8_t)(Motor2->foc_state->mode)};
}

bool Motor2_mode_write(uint8_t buf[7])
{
    Motor2->foc_state->mode = (foc_mode)buf[0];
    return true;
}

std::vector<uint8_t> Motor1_target_pos_read()
{
    return utils_float_to_bytes(Motor1->foc_state->targetAngle);
}

bool Motor1_target_pos_write(uint8_t buf[7])
{
    uint8_t newBuf[4];
    for (int i = 0; i < (int)sizeof(newBuf); i++)
    {
        newBuf[i] = buf[i];
    }
    Motor1->foc_state->targetAngle = utils_bytes_to_float(newBuf);
    return true;
}

std::vector<uint8_t> Motor2_target_pos_read()
{
    return utils_float_to_bytes(Motor2->foc_state->targetAngle);
}

bool Motor2_target_pos_write(uint8_t buf[7])
{
    uint8_t newBuf[4];
    for (int i = 0; i < (int)sizeof(newBuf); i++)
    {
        newBuf[i] = buf[i];
    }
    Motor2->foc_state->targetAngle = utils_bytes_to_float(newBuf);
    return true;
}

std::vector<uint8_t> Motor1_target_vel_read()
{
    return utils_float_to_bytes(Motor1->foc_state->targetVel);
}

bool Motor1_target_vel_write(uint8_t buf[7])
{
    uint8_t newBuf[4];
    for (int i = 0; i < (int)sizeof(newBuf); i++)
    {
        newBuf[i] = buf[i];
    }
    Motor1->foc_state->targetVel = utils_bytes_to_float(newBuf);
    return true;
}

std::vector<uint8_t> Motor2_target_vel_read()
{
    return utils_float_to_bytes(Motor2->foc_state->targetVel);
}

bool Motor2_target_vel_write(uint8_t buf[7])
{
    uint8_t newBuf[4];
    for (int i = 0; i < (int)sizeof(newBuf); i++)
    {
        newBuf[i] = buf[i];
    }
    Motor2->foc_state->targetVel = utils_bytes_to_float(newBuf);
    return true;
}

void loop()
{
    MotorLoop(Motor1);
    MotorLoop(Motor2);
}