#ifndef FOC_MATH_H_
#define FOC_MATH_H_

#include "datatypes.h"
#include <Arduino.h>
#include <pid.h>
#include <encoder.h>

// Types
typedef enum
{
    IQ_TARGET = 0U,
    VEL_TARGET = 1U,
    POS_TARGET = 2U,
} foc_mode;

typedef struct
{
    float va;
    float vb;
    float vc;
    float id_target;
    float iq_target;
    float max_duty;
    float duty_now;
    float phase;
    float phase_cos;
    float phase_sin;
    float i_a;
    float i_b;
    float i_c;
    float i_alpha;
    float i_beta;
    float i_abs;
    float i_bus;
    float v_bus;
    float v_alpha;
    float v_beta;
    float i_d;
    float i_q;
    Pid *pid_i_d;
    Pid *pid_i_q;
    float vd;
    float vq;
    uint32_t svm_sector;
    float dA;
    float dB;
    float dC;
    foc_mode mode;
    float targetAngle;
    float targetVel;
    float Ts;
} foc_state_t;

typedef struct
{
    float kp_d, ki_d, kd_d;
    float kp_q, ki_q, kd_q;
    float Ts;
} foc_config_t;

typedef struct
{
    float alpha, beta;
} foc_ab_frame_t;

typedef struct
{
    float d, q;
} foc_dq_frame_t;

typedef struct
{
    uint32_t tA, tB, tC;
} foc_phase_duty_timings_t;

// Functions
foc_state_t *init_foc(foc_config_t config);
void foc_drive(foc_state_t *, float theta);
void foc_svpwm(float alpha, float beta, float max_mod, uint32_t PWMFullDutyCycle,
               float *tAout, float *tBout, float *tCout, uint32_t *svm_sector);

class Foc
{
public:
    float va;
    float vb;
    float vc;
    float id_target;
    float iq_target;
    float max_duty;
    float duty_now;
    float phase;
    float phase_cos;
    float phase_sin;
    float i_a;
    float i_b;
    float i_c;
    float i_alpha;
    float i_beta;
    float i_abs;
    float i_bus;
    float v_bus;
    float v_alpha;
    float v_beta;
    float i_d;
    float i_q;
    Pid *pid_i_d;
    Pid *pid_i_q;
    float vd;
    float vq;
    uint32_t svm_sector;
    float dA;
    float dB;
    float dC;
    foc_mode mode;
    float targetAngle;
    float targetVel;
    float Ts;

    Foc(foc_config_t config);
    void Drive(float theta, motor_running_mode running_mode);

private:
    void Svpwm(float alpha, float beta, float max_mod, uint32_t PWMFullDutyCycle,
               float *tAout, float *tBout, float *tCout, uint32_t *svm_sector);
};

#endif /* FOC_MATH_H_ */