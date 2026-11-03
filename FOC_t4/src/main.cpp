#include <can.h>
#include <motor.h>
#include <util_math.h>
#include <motor_driver.h>
#include <TeensyThreads.h>
#include <analog.h>

volatile Motor *motor;
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
    pinMode(LED_BUILTIN, OUTPUT);
    digitalWrite(LED_BUILTIN, HIGH);
    Serial.begin(115200);
    Serial.println("Starting setup");
    // analogReadResolution(12);
    // analogReadAveraging(1);
    // ADC_ns::init();
    initCan();
    motor = new Motor();
    Serial.println("Finished setup");
    digitalWrite(LED_BUILTIN, LOW);
}

void loop()
{
    // elapsedMicros timer;
    CanUpdate();
    delay(250);

    digitalWriteFast(LED_BUILTIN, !digitalReadFast(LED_BUILTIN));
}