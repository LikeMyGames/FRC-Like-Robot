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

    Pid(float *process, float kp, float ki, float kd, float dt);
    void Update(float input);
};

#endif