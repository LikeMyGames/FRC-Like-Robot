#ifndef PID_H_
#define PID_H_

typedef struct
{
    float kp, ki, kd, izone, ff;
} PidConstants;

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
    volatile float output;
    float out_min;
    float out_max;
    float izone;
    float ff;

    Pid(float *process, float kp, float ki, float kd, float dt, float out_min, float out_max);
    volatile void Update(float input);
    void SetKp(float kp);
    void SetKi(float ki);
    void Setkd(float kd);
    void SetIzone(float izone);
    void SetFF(float ff);
};

#endif