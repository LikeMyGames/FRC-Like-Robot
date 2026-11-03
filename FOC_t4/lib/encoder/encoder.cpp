#include <encoder.h>
#include <util_math.h>
#include <analog.h>

Encoder::Encoder(uint pin)
{
    this->pin = pin;
    pinMode((uint8_t)pin, INPUT_PULLDOWN);
}

float Encoder::ReadRadians()
{
    // angle = utils_map_range((float)(this->Read()), 0.f, (float)res, 0.f, TWO_PI);
    // val = utils_map_range((float)angle, 0.f, TWO_PI, 0.f, (float)res);
    // return angle;
    return ADC_ns::readRadians(this->pin) + offset;
}

float Encoder::ReadDegrees()
{
    return RAD_TO_DEG * this->ReadRadians();
}

void Encoder::SetOffsetRadians(float radians)
{
    this->offset = radians;
}

void Encoder::SetOffsetDegrees(float degrees)
{
    this->offset = DEG_TO_RAD * degrees;
}

float Encoder_ns::countsToRad(int counts)
{
    return (float)(counts % COUNTS_PER_REV) * M_TWOPI;
}

float Encoder_ns::countsToDegrees(int counts)
{
    return (float)(counts % COUNTS_PER_REV) * 360.0f;
}
