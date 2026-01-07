#include <pid.h>

pid_controller_t *NewPidController(float *process, float kp, float ki, float kd, float dt)
{
    pid_controller_t *controller = {};
    controller->kp = kp;
    controller->ki = ki;
    controller->kd = kd;
    controller->dT = dt;
    controller->process = process;

    return controller;
}

void UpdatePid(pid_controller_t *pid, float input)
{
    float err = *(pid->process) - input;

    pid->integralTerm += (err * pid->dT);

    pid->derivativeTerm = (err - pid->lastErr) / pid->dT;
    pid->lastErr = err;

    pid->output = (pid->kp * err) + (pid->ki * pid->integralTerm) + (pid->kd * pid->derivativeTerm);
}