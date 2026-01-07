#ifndef PID_H_
#define PID_H_

typedef struct
{
    float kp, ki, kd;
    float *process;
    float dT;
    float propTerm;
    float integralTerm;
    float derivativeTerm;
    float lastErr;
    float output;
} pid_controller_t;

pid_controller_t *NewPidController(float *process, float kp, float ki, float kd, float dt);
void UpdatePid(pid_controller_t *pid, float input);

#endif