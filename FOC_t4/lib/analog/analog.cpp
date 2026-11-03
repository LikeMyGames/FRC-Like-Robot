#include <analog.h>
#include <ADC.h>
#include <util_math.h>

void ADC_ns::init()
{
    adc = new ADC();

    ///// ADC0 ////
    adc->adc0->setAveraging(4);   // set number of averages
    adc->adc0->setResolution(12); // set bits of resolution
    adc->adc0->setConversionSpeed(
        ADC_CONVERSION_SPEED::HIGH_SPEED); // change the conversion speed
    adc->adc0->setSamplingSpeed(
        ADC_SAMPLING_SPEED::HIGH_VERY_HIGH_SPEED); // change the sampling speed

////// ADC1 /////
#ifdef ADC_DUAL_ADCS
    adc->adc1->setAveraging(4);   // set number of averages
    adc->adc1->setResolution(12); // set bits of resolution
    adc->adc1->setConversionSpeed(
        ADC_CONVERSION_SPEED::HIGH_SPEED); // change the conversion speed
    adc->adc1->setSamplingSpeed(
        ADC_SAMPLING_SPEED::HIGH_VERY_HIGH_SPEED); // change the sampling speed
#endif

    delay(500);
}

int ADC_ns::readPin(int pin)
{
    if (pin == 24 || pin == 25)
    {
        return adc->adc0->analogRead(pin);
    }
    if (pin == 26 || pin == 27)
    {
        return adc->adc1->analogRead(pin);
    }
    if (pin <= 23 && 14 <= pin)
    {
        return adc->adc0->analogRead(pin);
    }
    return -1;
}

float ADC_ns::readVoltage(int pin)
{
    return ADC_ns::intToVoltage(ADC_ns::readPin(pin));
}

float ADC_ns::readRadians(int pin)
{
    return ADC_ns::intToRadians(ADC_ns::readPin(pin));
}

float ADC_ns::readDegrees(int pin)
{
    return ADC_ns::intToDegrees(ADC_ns::readPin(pin));
}

float ADC_ns::intToVoltage(int value)
{
    return (value / 4095.0f) * 3.3f;
}

float ADC_ns::intToRadians(int value)
{
    return (value / 4095.0f) * M_TWOPI;
}

float ADC_ns::intToDegrees(int value)
{
    return intToRadians(value) * RAD_TO_DEG;
}