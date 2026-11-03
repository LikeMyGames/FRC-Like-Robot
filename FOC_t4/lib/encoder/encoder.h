#ifndef ENCODER_H_
#define ENCODER_H_

#include <cstdint>
#include <Arduino.h>

#define COUNTS_PER_REV 42

typedef enum
{
    NO_ENCODER_ERROR = 0,
} encoder_error;

class Encoder
{
public:
    float offset;
    uint pin;
    uint16_t val;
    float angle;

    Encoder(uint pin);

    float ReadRadians();
    float ReadDegrees();
    void SetOffsetRadians(float radians);
    void SetOffsetDegrees(float degrees);
};

namespace Encoder_ns
{
    float countsToRad(int counts);
    float countsToDegrees(int counts);
}

#endif