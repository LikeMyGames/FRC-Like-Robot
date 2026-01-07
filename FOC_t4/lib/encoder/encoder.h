#ifndef ENCODER_H_
#define ENCODER_H_

#include <cstdint>
#include <Arduino.h>

typedef enum
{
    NO_ENCODER_ERROR = 0,
} encoder_error;

class Encoder
{
public:
    uint16_t offset;
    uint16_t res = 1023;
    uint pin;

    Encoder(uint pin);

    uint16_t Read();
    float ReadRad();
    void SetOffset(uint16_t offset);
    void SetRadOffset(float radOffset);
};

#endif