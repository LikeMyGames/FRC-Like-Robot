#ifndef FOC_MATH_H_
#define FOC_MATH_H_

#include "datatypes.h"
#include <Arduino.h>
#include <pid.h>
#include <encoder.h>
#include "motor_driver.h"

#define MOTOR_PWM_DUTY_CYCLE ((uint32_t)50)
#define BLDC_PWM_FREQ (1.f / ((float)MOTOR_PWM_DUTY_CYCLE / 1000000.f))
// #define FOC_TS ((float)MOTOR_PWM_DUTY_CYCLE / 1000000.f)
#define FOC_TS ((float)MOTOR_PWM_DUTY_CYCLE)

#define MAX_I_SAMPLE 128

#define MAX_CURRENT 0.25

#define PHASE_A_CURRENT_READ 23
#define PHASE_B_CURRENT_READ 22
#define PHASE_C_CURRENT_READ 21
#define CURRENT_READ_GAIN 50 // 100
#define CURRENT_READ_RESISTANCE 0.001

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
    volatile float iq_target;
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
    float bandwidth;
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

typedef enum
{
    RS,
    LD,
    LQ
} inject_target_t;

// Functions
foc_state_t *init_foc(foc_config_t config);
void foc_drive(foc_state_t *, float theta);
void foc_svpwm(float alpha, float beta, float max_mod, uint32_t PWMFullDutyCycle,
               float *tAout, float *tBout, float *tCout, uint32_t *svm_sector);
void clarke_transform(float ia, float ib, float *i_alpha, float *i_beta);
void park_transform(float i_alpha, float i_beta, float sin_theta, float cos_theta, float *id, float *iq);
void clarke_park_transform(float ia, float ib, float sin_theta, float cos_theta, float *id, float *iq);
void inverse_park_transform(float vd, float vq, float sin_theta, float cos_theta, float *valpha, float *vbeta);
void inverse_clarke_transform(float valpha, float vbeta, float *va, float *vb, float *vc);
void svpwm(float valpha, float vbeta, float vbus, uint32_t pwm_period,
           uint32_t *pwm_u, uint32_t *pwm_v, uint32_t *pwm_w);

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
    volatile float i_a;
    volatile float i_b;
    volatile float i_c;
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
    uint32_t dA;
    uint32_t dB;
    uint32_t dC;
    foc_mode mode;
    // float targetAngle;
    // float targetVel;
    float Ts;
    float R;
    float Ld;
    float Lq;
    float bandwidth;
    inject_target_t meas_inj_target;
    bool meas_inj_start_flag;
    float meas_inj_amp;
    float meas_inj_omega;
    float meas_inj_freq;
    int meas_inj_n;
    float Vd_buff[MAX_I_SAMPLE];
    float Id_buff[MAX_I_SAMPLE];
    float Vq_buff[MAX_I_SAMPLE];
    float Iq_buff[MAX_I_SAMPLE];

    Foc(foc_config_t config);
    void Drive(float theta);
    void Calibrate();
    void read_currents();

private:
    void open_loop_voltage_control(float vd_ref, float vq_ref, float angle_rad);
    void Svpwm(float alpha, float beta, float max_mod, uint32_t PWMFullDutyCycle,
               uint32_t *tAout, uint32_t *tBout, uint32_t *tCout);
    void start_measure(inject_target_t target);
    int measure_R(float vdc);
    void estimate_resistance();
    int measure_L(float f, float amp);
    void estimate_inductance(float ts);
    void meas_inj_dq_process(float ts);
};

#endif /* FOC_MATH_H_ */