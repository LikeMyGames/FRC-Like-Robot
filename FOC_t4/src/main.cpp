#include <can.h>
#include <motor.h>
#include <util_math.h>
#include <motor_driver.h>
#include <TeensyThreads.h>

Motor *motor;
// motor_config_t motor1_config = {
//     1.f,
//     21,
//     20,
//     17,
//     16,
//     50,
// };

void setup()
{
    digitalWrite(LED_BUILTIN, HIGH);
    Serial.begin(115200);
    Serial.println("Starting setup");
    initCan();
    motor = new Motor();
    Serial.println("Finished setup");
    digitalWrite(LED_BUILTIN, LOW);
}

void loop()
{
    motor->Update();
}