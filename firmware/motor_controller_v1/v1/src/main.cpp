#include <can.h>
#include <motor.h>
#include <pins_arduino.h>
#include <analog.h>
// #include <Arduino.h>

IntervalTimer blinkTimer;

void blink()
{
    if (Motor_ns::IsEnabled())
    {
        digitalWriteFast(LED_BUILTIN, !digitalReadFast(LED_BUILTIN));
    }
    else
    {
        digitalWriteFast(LED_BUILTIN, HIGH);
    }
}

void setup()
{
    pinMode(LED_BUILTIN, OUTPUT);
    digitalWriteFast(LED_BUILTIN, HIGH);
    Serial.begin(115200);
    Serial.println("Starting setup");
    ADC_ns::init();
    initCan();
    Motor_ns::Init();
    Serial.println("Finished setup");

    blinkTimer.begin(blink, 250000);

    digitalWriteFast(LED_BUILTIN, LOW);
}

void loop()
{
    CanUpdate();
    delay(1);
}
