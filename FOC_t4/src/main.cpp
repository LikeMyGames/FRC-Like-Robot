#include <can.h>
#include <motor.h>
#include <util_math.h>
#include <motor_driver.h>

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

Motor *Motor1;
static motor_config_t motor1_config = {
    1.f,
    21U,
    20U,
    50,
};

Motor *Motor2;
static motor_config_t motor2_config = {
    1.f,
    19,
    18,
    50,
};

void setup()
{
    digitalWrite(LED_BUILTIN, HIGH);
    Serial.begin(115200);
    // Serial.setTimeout(0);
    Serial.println("Starting setup");
    initCan();
    Serial.println("finished motor controller setup");
    Motor1 = new Motor(motor1_config);
    Serial.println("created motor 1");
    Motor2 = new Motor(motor2_config);
    Serial.println("created motor 2");
    Motor1->Init(&motor_driver_ns::motor_driver_1);
    Serial.println("started motor 1 driver");
    Motor2->Init(&motor_driver_ns::motor_driver_1);
    // (&motor_driver_ns::motor_driver_2);
    Serial.println("started motor 2 driver");
    addCanRegister(0, "echo", CanEchoRead, CanEchoWrite);
    Serial.println("created can register 0");
    addCanRegister(1, "motor1 mode", Motor1_mode_read, Motor1_mode_write);
    Serial.println("created can register 1");
    addCanRegister(2, "motor2 mode", Motor2_mode_read, Motor2_mode_write);
    Serial.println("created can register 2");
    addCanRegister(3, "motor1 target iq", Motor1_Iq_read, Motor1_Iq_write);
    Serial.println("created can register 3");
    addCanRegister(4, "motor2 target iq", Motor2_Iq_read, Motor2_Iq_write);
    Serial.println("created can register 4");
    addCanRegister(5, "motor1 target pos", Motor1_target_pos_read, Motor1_target_pos_write);
    Serial.println("created can register 5");
    addCanRegister(6, "motor1 target vel", Motor1_target_vel_read, Motor1_target_vel_write);
    Serial.println("created can register 6");
    addCanRegister(7, "motor2 target pos", Motor2_target_pos_read, Motor2_target_pos_write);
    Serial.println("created can register 7");
    addCanRegister(8, "motor2 target vel", Motor2_target_vel_read, Motor2_target_vel_write);
    Serial.println("created can register 8");
    addCanRegister(9, "motor1 pos pid kp", []() -> std::vector<uint8_t>
                   { return utils_float_to_bytes(Motor1->pos_pid->kp); }, [](uint8_t buf[7]) -> bool
                   { uint8_t newBuf[4];
                        for (int i = 0; i < (int)sizeof(newBuf); i++)
                        {
                            newBuf[i] = buf[i];
                        }
                        Motor1->pos_pid->kp = utils_bytes_to_float(newBuf);
                        return true; });
    addCanRegister(9, "motor1 pos pid ki", []() -> std::vector<uint8_t>
                   { return utils_float_to_bytes(Motor1->pos_pid->ki); }, [](uint8_t buf[7]) -> bool
                   { uint8_t newBuf[4];
                        for (int i = 0; i < (int)sizeof(newBuf); i++)
                        {
                            newBuf[i] = buf[i];
                        }
                        Motor1->pos_pid->ki = utils_bytes_to_float(newBuf);
                        return true; });
    addCanRegister(9, "motor1 pos pid kd", []() -> std::vector<uint8_t>
                   { return utils_float_to_bytes(Motor1->pos_pid->kd); }, [](uint8_t buf[7]) -> bool
                   { uint8_t newBuf[4];
                        for (int i = 0; i < (int)sizeof(newBuf); i++)
                        {
                            newBuf[i] = buf[i];
                        }
                        Motor1->pos_pid->kd = utils_bytes_to_float(newBuf);
                        return true; });
    addCanRegister(9, "motor1 vel pid kp", []() -> std::vector<uint8_t>
                   { return utils_float_to_bytes(Motor1->vel_pid->kp); }, [](uint8_t buf[7]) -> bool
                   { uint8_t newBuf[4];
                        for (int i = 0; i < (int)sizeof(newBuf); i++)
                        {
                            newBuf[i] = buf[i];
                        }
                        Motor1->vel_pid->kp = utils_bytes_to_float(newBuf);
                        return true; });
    addCanRegister(9, "motor1 vel pid ki", []() -> std::vector<uint8_t>
                   { return utils_float_to_bytes(Motor1->vel_pid->ki); }, [](uint8_t buf[7]) -> bool
                   { uint8_t newBuf[4];
                        for (int i = 0; i < (int)sizeof(newBuf); i++)
                        {
                            newBuf[i] = buf[i];
                        }
                        Motor1->vel_pid->ki = utils_bytes_to_float(newBuf);
                        return true; });
    addCanRegister(9, "motor1 vel pid kd", []() -> std::vector<uint8_t>
                   { return utils_float_to_bytes(Motor1->vel_pid->kd); }, [](uint8_t buf[7]) -> bool
                   { uint8_t newBuf[4];
                        for (int i = 0; i < (int)sizeof(newBuf); i++)
                        {
                            newBuf[i] = buf[i];
                        }
                        Motor1->vel_pid->kd = utils_bytes_to_float(newBuf);
                        return true; });
    digitalWrite(LED_BUILTIN, LOW);
    Motor1->foc->iq_target = 1;
}

std::vector<uint8_t> Motor1_Iq_read()
{
    return utils_float_to_bytes(Motor1->foc->iq_target);
}

bool Motor1_Iq_write(uint8_t buf[7])
{
    uint8_t newBuf[4];
    for (int i = 0; i < (int)sizeof(newBuf); i++)
    {
        newBuf[i] = buf[i];
    }
    Motor1->foc->iq_target = utils_bytes_to_float(newBuf);
    return true;
}

std::vector<uint8_t> Motor2_Iq_read()
{
    return utils_float_to_bytes(Motor2->foc->iq_target);
}

bool Motor2_Iq_write(uint8_t buf[7])
{
    uint8_t newBuf[4];
    for (int i = 0; i < (int)sizeof(newBuf); i++)
    {
        newBuf[i] = buf[i];
    }
    Motor2->foc->iq_target = utils_bytes_to_float(newBuf);
    return true;
}

std::vector<uint8_t> Motor1_mode_read()
{
    return {(uint8_t)(Motor1->foc->mode)};
}

bool Motor1_mode_write(uint8_t buf[7])
{
    Motor1->foc->mode = (foc_mode)buf[0];
    return true;
}

std::vector<uint8_t> Motor2_mode_read()
{
    return {(u_int8_t)(Motor2->foc->mode)};
}

bool Motor2_mode_write(uint8_t buf[7])
{
    Motor2->foc->mode = (foc_mode)buf[0];
    return true;
}

std::vector<uint8_t> Motor1_target_pos_read()
{
    return utils_float_to_bytes(Motor1->foc->targetAngle);
}

bool Motor1_target_pos_write(uint8_t buf[7])
{
    uint8_t newBuf[4];
    for (int i = 0; i < (int)sizeof(newBuf); i++)
    {
        newBuf[i] = buf[i];
    }
    Motor1->foc->targetAngle = utils_bytes_to_float(newBuf);
    return true;
}

std::vector<uint8_t> Motor2_target_pos_read()
{
    return utils_float_to_bytes(Motor2->foc->targetAngle);
}

bool Motor2_target_pos_write(uint8_t buf[7])
{
    uint8_t newBuf[4];
    for (int i = 0; i < (int)sizeof(newBuf); i++)
    {
        newBuf[i] = buf[i];
    }
    Motor2->foc->targetAngle = utils_bytes_to_float(newBuf);
    return true;
}

std::vector<uint8_t> Motor1_target_vel_read()
{
    return utils_float_to_bytes(Motor1->foc->targetVel);
}

bool Motor1_target_vel_write(uint8_t buf[7])
{
    uint8_t newBuf[4];
    for (int i = 0; i < (int)sizeof(newBuf); i++)
    {
        newBuf[i] = buf[i];
    }
    Motor1->foc->targetVel = utils_bytes_to_float(newBuf);
    return true;
}

std::vector<uint8_t> Motor2_target_vel_read()
{
    return utils_float_to_bytes(Motor2->foc->targetVel);
}

bool Motor2_target_vel_write(uint8_t buf[7])
{
    uint8_t newBuf[4];
    for (int i = 0; i < (int)sizeof(newBuf); i++)
    {
        newBuf[i] = buf[i];
    }
    Motor2->foc->targetVel = utils_bytes_to_float(newBuf);
    return true;
}

void loop()
{
    Motor1->Update();
    Motor2->Update();
    delayMicroseconds(50);
}