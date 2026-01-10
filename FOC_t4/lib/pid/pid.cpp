#include <pid.h>

Pid::Pid(float *process, float kp, float ki, float kd, float dt)
{
    this->kp = kp;
    this->ki = ki;
    this->kd = kd;
    this->dT = dt;
    this->process = process;
}

void Pid::Update(float input)
{
    float err = *(process)-input;

    integralTerm += (err * dT);

    derivativeTerm = (err - lastErr) / dT;
    lastErr = err;

    output = (kp * err) + (ki * integralTerm) + (kd * derivativeTerm);
}