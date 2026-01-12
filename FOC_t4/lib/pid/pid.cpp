#include <pid.h>

Pid::Pid(float *process, float kp, float ki, float kd, float dt, float out_min, float out_max)
{
    this->kp = kp;
    this->ki = ki;
    this->kd = kd;
    this->dT = dt;
    this->process = process;
    this->out_min = out_min;
    this->out_max = out_max;
}

void Pid::Update(float input)
{
    float err = *(process)-input;

    integralTerm += (err * dT);

    derivativeTerm = (err - lastErr) / dT;
    lastErr = err;

    float integrator = (ki * integralTerm);

    if (integrator > out_max)
        integrator = out_max;
    if (integrator < out_min)
        integrator = out_min;

    output = (kp * err) + integrator + (kd * derivativeTerm);

    if (output > out_max)
        output = out_max;
    if (output < out_min)
        output = out_min;
}