#include <encoder.h>
#include <util_math.h>
#include <Arduino.h>

Encoder::Encoder(uint pin)
{
    this->pin = pin;
    pinMode((uint8_t)pin, INPUT_PULLDOWN);
}

uint16_t Encoder::Read()
{
    return analogRead(this->pin) + this->offset;
}

float Encoder::ReadRad()
{
    return utils_map_range((float)(this->Read()), 0.f, (float)res, 0.f, (float)TWO_PI);
}

void Encoder::SetOffset(uint16_t offset)
{
    this->offset = offset;
}

void Encoder::SetRadOffset(float radOffset)
{
    this->offset = (uint16_t)utils_map_range(radOffset, 0.f, (float)TWO_PI, 0.f, (float)res);
}