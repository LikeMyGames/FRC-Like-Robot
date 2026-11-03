#include <ADC.h>

namespace ADC_ns
{
    // do not use out side of package
    inline ADC *adc;

    void init();
    int readPin(int pin);
    float readVoltage(int pin);
    float readRadians(int pin);
    float readDegrees(int pin);
    inline float intToVoltage(int value);
    inline float intToRadians(int value);
    inline float intToDegrees(int value);
}
