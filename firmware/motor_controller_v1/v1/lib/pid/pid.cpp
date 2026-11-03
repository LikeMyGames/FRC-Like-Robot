#include <pid.h>
#include <stdlib.h>

Pid::Pid(float *process, float kp, float ki, float kd, float dt, float out_min, float out_max)
{
    this->kp = kp;
    this->ki = ki;
    this->kd = kd;
    this->dT = dt;
    this->process = process;
    this->out_min = out_min;
    this->out_max = out_max;
    this->izone = 0.f;
    this->ff = 0.f;
}

volatile void Pid::Update(float input)
{
    float err = *(process)-input;

    this->derivativeTerm = (err - lastErr) / dT;
    this->lastErr = err;

    float integrator = (ki * integralTerm);

    this->integralTerm += (err * dT);

    if (izone > 0.0f && abs(err) > izone)
        integrator = 0;

    output = (kp * err) + integrator + (kd * derivativeTerm) + ff;

    if (output > out_max)
        output = out_max;
    if (output < out_min)
        output = out_min;
}