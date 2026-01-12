#ifndef PID_H_
#define PID_H_

class Pid
{
public:
    volatile float kp, ki, kd;
    float *process;
    float dT;
    float propTerm;
    float integralTerm;
    float derivativeTerm;
    float lastErr;
    float output;
    float out_min;
    float out_max;

    Pid(float *process, float kp, float ki, float kd, float dt, float out_min, float out_max);
    void Update(float input);
};

#endif