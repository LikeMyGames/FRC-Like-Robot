#include <foc_math.h>
#include <util_math.h>
#include <math.h>
#include <analog.h>

void clarke_transform(float ia, float ib, float *i_alpha, float *i_beta)
{
    // Clarke transform
    *i_alpha = ia;
    *i_beta = ONE_BY_SQRT3 * ia + TWO_BY_SQRT3 * ib;
}

void park_transform(float i_alpha, float i_beta, float sin_theta, float cos_theta, float *id, float *iq)
{
    // Park transform
    *id = i_alpha * cos_theta + i_beta * sin_theta;
    *iq = i_beta * cos_theta - i_alpha * sin_theta;
}

// Fast combined Clarke + Park Transform
void clarke_park_transform(float ia, float ib, float sin_theta, float cos_theta, float *id, float *iq)
{
    // Clarke transform
    float i_alpha = ia;
    float i_beta = ONE_BY_SQRT3 * ia + TWO_BY_SQRT3 * ib;

    // Park transform
    *id = i_alpha * cos_theta + i_beta * sin_theta;
    *iq = i_beta * cos_theta - i_alpha * sin_theta;
}

// Inverse Park Transform
void inverse_park_transform(float vd, float vq, float sin_theta, float cos_theta, float *valpha, float *vbeta)
{
    *valpha = vd * cos_theta - vq * sin_theta;
    *vbeta = vd * sin_theta + vq * cos_theta;
}

// Inverse Clarke Transform
void inverse_clarke_transform(float valpha, float vbeta, float *va, float *vb, float *vc)
{
    *va = valpha;
    *vb = -0.5f * valpha + SQRT3_BY_2 * vbeta; // cos(120°), sin(120°)
    *vc = -0.5f * valpha - SQRT3_BY_2 * vbeta;
}

Foc::Foc(foc_config_t config)
{
    Ts = config.Ts;
    Serial.println("added foc config");

    pid_i_d = new Pid(&(id_target), config.kp_d, config.ki_d, config.kd_d, config.Ts, -MAX_CURRENT, MAX_CURRENT);
    Serial.println("created i_d pid controller");
    pid_i_q = new Pid(&(iq_target), config.kp_q, config.ki_q, config.kd_q, config.Ts, -MAX_CURRENT, MAX_CURRENT);
    Serial.println("created i_q pid controller");

    bandwidth = config.bandwidth;

    // this->driver = &motor_driver_ns::motor_driver_1;
    Serial.println("created motor driver");
    // motor_driver_ns::InitMotorDriver(this->driver, BLDC_PWM_FREQ);
    driver::Init(BLDC_PWM_FREQ);

    Calibrate();

    driver::DriveByTime(0, 0, 0);
    driver::Disable();
}

void Foc::Drive(float theta)
{
    this->read_currents();
    float sin_theta, cos_theta = 0.0f;
    utils_fast_sincos(theta, &sin_theta, &cos_theta);
    clarke_park_transform(i_a, i_b, sin_theta, cos_theta, &i_d, &i_q);

    // d and q pid updates
    pid_i_d->Update(i_d);
    pid_i_q->Update(i_q);

    Serial.printf("PID ERR i_d=%f i_q=%f\n", pid_i_d->lastErr, pid_i_q->lastErr);

    Serial.printf("v_d=%f v_q=%f\n", pid_i_d->output, pid_i_q->output);

    inverse_park_transform(pid_i_d->output, pid_i_q->output, sin_theta, cos_theta, &v_alpha, &v_beta);

    Serial.printf("v_alpha=%f v_beta=%f\n", v_alpha, v_beta);

    // svpwm(v_alpha, v_beta, 12, Ts, &dA, &dB, &dC);
    svpwm(v_alpha, v_beta, 5, Ts, &dA, &dB, &dC);
    // this->Svpwm(v_alpha, v_beta, float max_mod, int *tAout, int *tBout, int *tCout);
    driver::DriveByTime(this->dA, this->dB, this->dC);
}

void Foc::Calibrate()
{
    // if (hfoc.foc_mode == FOC_MODE_SENSORED || hfoc.foc_mode == FOC_MODE_HYBRID)
    // {
    //     foc_cal_encoder(&hfoc);
    // }

    IntervalTimer meas_inj_timer;
    Serial.println("starting meas_inj_timer");
    meas_inj_timer.begin([this]() -> void
                         { this->meas_inj_dq_process(FOC_TS); }, 50);

    Serial.println("starting resistance measurement");
    measure_R(1.0f);
    Serial.println("finished resistance measurement");
    Serial.println("starting inductance measurement");
    measure_L(1000.0f, 1.2f);
    Serial.println("finished inductance measurement");

    Serial.println("stopping meas_inj_timer");
    meas_inj_timer.end();

    // smo_update_R_L(&hfoc.smo, m_config.Rs, (m_config.Ld + m_config.Lq) * 0.5);

    // flash_auto_tuning_torque_control(&m_config);
    this->pid_i_d->kp = bandwidth * Ld;
    this->pid_i_d->ki = bandwidth * R;
    this->pid_i_q->kp = bandwidth * Lq;
    this->pid_i_q->kp = bandwidth * R;

    Serial.printf("Calculated resistance: %f\n", R);
    Serial.printf("Calculated inductance: Ld= %f Lq=%f\n", Ld, Lq);
    Serial.printf("New pid gains: pid_i_d->kp=%f pid_i_d->ki=%f pid_i_q->kp=%f pid_i_q->ki=%f\n", pid_i_d->kp, pid_i_d->ki, pid_i_q->kp, pid_i_q->ki);
}

void Foc::open_loop_voltage_control(float vd_ref, float vq_ref,
                                    float angle_rad)
{
    uint32_t da, db, dc;
    float sin_theta, cos_theta;
    utils_fast_sincos(angle_rad, &sin_theta, &cos_theta);
    inverse_park_transform(vd_ref, vq_ref, sin_theta, cos_theta, &this->v_alpha, &this->v_beta);
    // Svpwm(this->v_alpha, this->v_beta, (float)1, this->Ts, &da, &db, &dc);
    svpwm(this->v_alpha, this->v_beta, 12, this->Ts, &da, &db, &dc);
    driver::DriveByTime(da, db, dc);
    // DRV8302_set_pwm(&hfoc->drv8302, da, db, dc);
}

/**
 * @brief svm Space vector modulation. Magnitude must not be larger than
 * sqrt(3)/2, or 0.866 to avoid overmodulation. See
 * https://github.com/vedderb/bldc/pull/372#issuecomment-962499623 for a full
 * description.
 * @param alpha voltage
 * @param beta Park transformed and normalized voltage
 * @param PWMFullDutyCycle is the peak value of the PWM counter.
 * @param tAout PWM duty cycle phase A (0 = off all of the time,
 * PWMFullDutyCycle = on all of the time)
 * @param tBout PWM duty cycle phase B
 * @param tCout PWM duty cycle phase C
 */
void Foc::Svpwm(float alpha, float beta, float max_mod, uint32_t PWMFullDutyCycle,
                uint32_t *tAout, uint32_t *tBout, uint32_t *tCout)
{
    uint32_t sector;

    if (beta >= 0.0f)
    {
        if (alpha >= 0.0f)
        {
            // quadrant I
            if (ONE_BY_SQRT3 * beta > alpha)
            {
                sector = 2;
            }
            else
            {
                sector = 1;
            }
        }
        else
        {
            // quadrant II
            if (-ONE_BY_SQRT3 * beta > alpha)
            {
                sector = 3;
            }
            else
            {
                sector = 2;
            }
        }
    }
    else
    {
        if (alpha >= 0.0f)
        {
            // quadrant IV5
            if (-ONE_BY_SQRT3 * beta > alpha)
            {
                sector = 5;
            }
            else
            {
                sector = 6;
            }
        }
        else
        {
            // quadrant III
            if (ONE_BY_SQRT3 * beta > alpha)
            {
                sector = 4;
            }
            else
            {
                sector = 5;
            }
        }
    }

    // PWM timings
    int tA, tB, tC;

    switch (sector)
    {

    // sector 1-2
    case 1:
    {
        // Vector on-times
        int t1 = (alpha - ONE_BY_SQRT3 * beta) * PWMFullDutyCycle;
        int t2 = (TWO_BY_SQRT3 * beta) * PWMFullDutyCycle;

        // PWM timings
        tA = (PWMFullDutyCycle + t1 + t2) / 2;
        tB = tA - t1;
        tC = tB - t2;

        break;
    }

    // sector 2-3
    case 2:
    {
        // Vector on-times
        int t2 = (alpha + ONE_BY_SQRT3 * beta) * PWMFullDutyCycle;
        int t3 = (-alpha + ONE_BY_SQRT3 * beta) * PWMFullDutyCycle;

        // PWM timings
        tB = (PWMFullDutyCycle + t2 + t3) / 2;
        tA = tB - t3;
        tC = tA - t2;

        break;
    }

    // sector 3-4
    case 3:
    {
        // Vector on-times
        int t3 = (TWO_BY_SQRT3 * beta) * PWMFullDutyCycle;
        int t4 = (-alpha - ONE_BY_SQRT3 * beta) * PWMFullDutyCycle;

        // PWM timings
        tB = (PWMFullDutyCycle + t3 + t4) / 2;
        tC = tB - t3;
        tA = tC - t4;

        break;
    }

    // sector 4-5
    case 4:
    {
        // Vector on-times
        int t4 = (-alpha + ONE_BY_SQRT3 * beta) * PWMFullDutyCycle;
        int t5 = (-TWO_BY_SQRT3 * beta) * PWMFullDutyCycle;

        // PWM timings
        tC = (PWMFullDutyCycle + t4 + t5) / 2;
        tB = tC - t5;
        tA = tB - t4;

        break;
    }

    // sector 5-6
    case 5:
    {
        // Vector on-times
        int t5 = (-alpha - ONE_BY_SQRT3 * beta) * PWMFullDutyCycle;
        int t6 = (alpha - ONE_BY_SQRT3 * beta) * PWMFullDutyCycle;

        // PWM timings
        tC = (PWMFullDutyCycle + t5 + t6) / 2;
        tA = tC - t5;
        tB = tA - t6;

        break;
    }

    // sector 6-1
    case 6:
    {
        // Vector on-times
        int t6 = (-TWO_BY_SQRT3 * beta) * PWMFullDutyCycle;
        int t1 = (alpha + ONE_BY_SQRT3 * beta) * PWMFullDutyCycle;

        // PWM timings
        tA = (PWMFullDutyCycle + t6 + t1) / 2;
        tC = tA - t1;
        tB = tC - t6;

        break;
    }
    }

    int t_max = PWMFullDutyCycle * (1.0 - (1.0 - max_mod) * 0.5);
    utils_truncate_number_int(&tA, 0, t_max);
    utils_truncate_number_int(&tB, 0, t_max);
    utils_truncate_number_int(&tC, 0, t_max);

    *tAout = tA;
    *tBout = tB;
    *tCout = tC;
    Serial.printf("tA=%d tB=%d tC=%d\n", tA, tB, tC);
    this->svm_sector = sector;
}

void svpwm(float valpha, float vbeta, float vbus, uint32_t pwm_period,
           uint32_t *pwm_u, uint32_t *pwm_v, uint32_t *pwm_w)
{
    // 1. Normalize voltages by vbus
    float alpha = valpha / vbus;
    float beta = vbeta / vbus;

    // 2. Sector determination
    uint8_t sector;
    if (beta >= 0.0f)
    {
        if (alpha >= 0.0f)
        {
            sector = (ONE_BY_SQRT3 * beta > alpha) ? 2 : 1; // 1/sqrt(3) ≈ 0.577
        }
        else
        {
            sector = (-ONE_BY_SQRT3 * beta > alpha) ? 3 : 2;
        }
    }
    else
    {
        if (alpha >= 0.0f)
        {
            sector = (-ONE_BY_SQRT3 * beta > alpha) ? 5 : 6;
        }
        else
        {
            sector = (ONE_BY_SQRT3 * beta > alpha) ? 4 : 5;
        }
    }

    // 3. Calculate active vector times
    int32_t t1, t2;
    switch (sector)
    {
    case 1:
        t1 = (int32_t)((alpha - ONE_BY_SQRT3 * beta) * pwm_period);
        t2 = (int32_t)(TWO_BY_SQRT3 * beta * pwm_period);
        *pwm_u = (pwm_period + t1 + t2) / 2;
        *pwm_v = *pwm_u - t1;
        *pwm_w = *pwm_v - t2;
        break;

    case 2:
        t1 = (int32_t)((alpha + ONE_BY_SQRT3 * beta) * pwm_period);
        t2 = (int32_t)((-alpha + ONE_BY_SQRT3 * beta) * pwm_period);
        *pwm_v = (pwm_period + t1 + t2) / 2;
        *pwm_u = *pwm_v - t2;
        *pwm_w = *pwm_u - t1;
        break;

    case 3:
        t1 = (int32_t)(TWO_BY_SQRT3 * beta * pwm_period);
        t2 = (int32_t)((-alpha - ONE_BY_SQRT3 * beta) * pwm_period);
        *pwm_v = (pwm_period + t1 + t2) / 2;
        *pwm_w = *pwm_v - t1;
        *pwm_u = *pwm_w - t2;
        break;

    case 4:
        t1 = (int32_t)((-alpha + ONE_BY_SQRT3 * beta) * pwm_period);
        t2 = (int32_t)(-TWO_BY_SQRT3 * beta * pwm_period);
        *pwm_w = (pwm_period + t1 + t2) / 2;
        *pwm_v = *pwm_w - t2;
        *pwm_u = *pwm_v - t1;
        break;

    case 5:
        t1 = (int32_t)((-alpha - ONE_BY_SQRT3 * beta) * pwm_period);
        t2 = (int32_t)((alpha - ONE_BY_SQRT3 * beta) * pwm_period);
        *pwm_w = (pwm_period + t1 + t2) / 2;
        *pwm_u = *pwm_w - t1;
        *pwm_v = *pwm_u - t2;
        break;

    case 6:
        t1 = (int32_t)(-TWO_BY_SQRT3 * beta * pwm_period);
        t2 = (int32_t)((alpha + ONE_BY_SQRT3 * beta) * pwm_period);
        *pwm_u = (pwm_period + t1 + t2) / 2;
        *pwm_w = *pwm_u - t2;
        *pwm_v = *pwm_w - t1;
        break;
    }

    // 4. Clamp outputs to valid range
    *pwm_u = (*pwm_u > pwm_period) ? pwm_period : *pwm_u;
    *pwm_v = (*pwm_v > pwm_period) ? pwm_period : *pwm_v;
    *pwm_w = (*pwm_w > pwm_period) ? pwm_period : *pwm_w;

    Serial.printf("pwm_u=%d pwm_v=%d pwm_w=%d\n", *pwm_u, *pwm_v, *pwm_w);
}

void Foc::start_measure(inject_target_t target)
{
    meas_inj_target = target;
    delay(100);
    meas_inj_start_flag = 1;
    while (meas_inj_start_flag)
    {
        // Serial.println("waiting for meas_inj_start_flag==0");
        delay(1);
    }
    // Serial.println("got meas_inj_start_flag");
}

int Foc::measure_R(float vdc)
{
    // future voltage precheck, need to implement voltage measurement
    /*
    if (v_bus > 12) {
        return -1;
    }
    */

    meas_inj_amp = vdc;

    memset(Vd_buff, 0, sizeof(Vd_buff));
    memset(Id_buff, 0, sizeof(Id_buff));

    start_measure(RS);
    estimate_resistance();

    return 0;
}

void Foc::estimate_resistance()
{
    float mean_vd = 0, mean_id = 0;

    for (int i = 0; i < MAX_I_SAMPLE; i++)
    {
        mean_vd += Vd_buff[i];
        mean_id += Id_buff[i];
    }

    Serial.printf("mean_vd=%f mean_id=%f\n", mean_vd, mean_id);
    mean_vd /= MAX_I_SAMPLE;
    mean_id /= MAX_I_SAMPLE;

    Serial.printf("averaged mean_vd=%f mean_id=%f\n", mean_vd, mean_id);

    this->R = mean_vd / mean_id;
}

int Foc::measure_L(float f, float amp)
{
    // future voltage precheck, need to implement voltage measurement
    /*
    if (v_bus > 12) {
        return -1;
    }
    */

    this->meas_inj_freq = f;
    this->meas_inj_amp = amp;
    this->meas_inj_omega = M_2_PI * f;

    memset(Vd_buff, 0, sizeof(Vd_buff));
    memset(Id_buff, 0, sizeof(Id_buff));
    memset(Vq_buff, 0, sizeof(Vq_buff));
    memset(Iq_buff, 0, sizeof(Iq_buff));

    start_measure(LD);
    start_measure(LQ);
    estimate_inductance(FOC_TS);

    //   m_config.Ld = hfoc.Ld;
    //   m_config.Lq = hfoc.Lq;
    return 0;
}

void Foc::estimate_inductance(float ts)
{
    float Vc_d = 0, Vs_d = 0, Ic_d = 0, Is_d = 0;
    float Vc_q = 0, Vs_q = 0, Ic_q = 0, Is_q = 0;
    float mean_vd = 0, mean_vq = 0, mean_id = 0, mean_iq = 0;

    // Remove DC offset
    for (int i = 0; i < MAX_I_SAMPLE; i++)
    {
        mean_vd += Vd_buff[i];
        mean_vq += Vq_buff[i];
        mean_id += Id_buff[i];
        mean_iq += Iq_buff[i];
    }
    mean_vd /= MAX_I_SAMPLE;
    mean_vq /= MAX_I_SAMPLE;
    mean_id /= MAX_I_SAMPLE;
    mean_iq /= MAX_I_SAMPLE;

    // Single-frequency DFT
    for (int i = 0; i < MAX_I_SAMPLE; i++)
    {
        float angle = meas_inj_omega * i * ts;

        float vd = Vd_buff[i] - mean_vd;
        float vq = Vq_buff[i] - mean_vq;
        float id = Id_buff[i] - mean_id;
        float iq = Iq_buff[i] - mean_iq;

        float cos, sin = 0.0;
        utils_fast_sincos(angle, &cos, &sin);

        Vc_d += vd * cos;
        Vs_d += vd * sin;
        Ic_d += id * cos;
        Is_d += id * sin;

        Vc_q += vq * cos;
        Vs_q += vq * sin;
        Ic_q += iq * cos;
        Is_q += iq * sin;
    }

    float norm = 2.0f / MAX_I_SAMPLE;
    Vc_d *= norm;
    Vs_d *= norm;
    Ic_d *= norm;
    Is_d *= norm;
    Vc_q *= norm;
    Vs_q *= norm;
    Ic_q *= norm;
    Is_q *= norm;

    // V & I amplitude
    float Vd_mag = sqrtf(Vc_d * Vc_d + Vs_d * Vs_d);
    float Id_mag = sqrtf(Ic_d * Ic_d + Is_d * Is_d);
    float Vq_mag = sqrtf(Vc_q * Vc_q + Vs_q * Vs_q);
    float Iq_mag = sqrtf(Ic_q * Ic_q + Is_q * Is_q);

    // (phi = arctan(Vs/Vc) - arctan(Is/Ic))
    float phi_d = atan2f(Vs_d, Vc_d) - atan2f(Is_d, Ic_d);
    float phi_q = atan2f(Vs_q, Vc_q) - atan2f(Is_q, Ic_q);

    // Impedansi & parameters
    float Zd_mag = Vd_mag / Id_mag;
    float Zq_mag = Vq_mag / Iq_mag;

    // float Rs_d = Zd_mag * cosf(phi_d);
    // float Rs_q = Zq_mag * cosf(phi_q);
    float Ld_est = (Zd_mag * sinf(phi_d)) / this->meas_inj_omega;
    float Lq_est = (Zq_mag * sinf(phi_q)) / this->meas_inj_omega;

    // hfoc->Rs = (Rs_d + Rs_q) * 0.5;
    this->Ld = fabs(Ld_est);
    this->Lq = fabs(Lq_est);
}

void Foc::meas_inj_dq_process(float ts)
{
    const uint32_t wt = 256; // waiting time

    if (this->meas_inj_start_flag)
    {
        Serial.println("running meas_inj process");
        Serial.println(this->meas_inj_n);
        Serial.println(this->meas_inj_omega);
        float vd = 0.0f;
        float vq = 0.0f;
        if (this->meas_inj_target == RS)
        {
            vd = this->meas_inj_amp;
            vq = 0.0f;
        }
        else
        {
            float angle = this->meas_inj_omega * this->meas_inj_n * ts;
            float v_inj = this->meas_inj_amp * utils_fast_sin(angle);
            if (this->meas_inj_target == LD)
            {
                vd = v_inj;
                vq = 0.0f;
            }
            else if (this->meas_inj_target == LQ)
            {
                vd = 0.0f;
                vq = v_inj;
            }
        }

        // float theta_e = hfoc->e_angle_rad_comp;
        float theta_e = 0.0f;

        open_loop_voltage_control(vd, vq, theta_e);

        // DRV8302_get_current(&hfoc->drv8302, &hfoc->ia, &hfoc->ib);
        this->read_currents();

        float sin_theta, cos_theta;
        float id, iq;
        utils_fast_sincos(theta_e, &sin_theta, &cos_theta);
        clarke_park_transform(this->i_a, this->i_b, sin_theta, cos_theta, &id, &iq);

        if (this->meas_inj_n >= (int)wt)
        {
            if (this->meas_inj_target == RS || this->meas_inj_target == LD)
            {
                Vd_buff[this->meas_inj_n - wt] = vd;
                Id_buff[this->meas_inj_n - wt] = id;
                Serial.println("setting vd_buff and id_buff");
            }
            else if (this->meas_inj_target == LQ)
            {
                Vq_buff[this->meas_inj_n - wt] = vq;
                Iq_buff[this->meas_inj_n - wt] = iq;
            }
        }

        this->meas_inj_n++;
        if (this->meas_inj_n >= (int)(MAX_I_SAMPLE + wt))
        {
            this->meas_inj_n = 0;
            this->meas_inj_start_flag = 0;
            this->open_loop_voltage_control(0, 0, 0);
        }
    }
}

void Foc::read_currents()
{
    float v_a = ADC_ns::readVoltage(PHASE_A_CURRENT_READ);
    float gained_v_a = v_a / CURRENT_READ_GAIN;
    this->i_a = gained_v_a / CURRENT_READ_RESISTANCE;

    // this->i_a = (ADC_ns::readVoltage(PHASE_A_CURRENT_READ) / CURRENT_READ_GAIN) / CURRENT_READ_RESISTANCE;
    this->i_b = (ADC_ns::readVoltage(PHASE_B_CURRENT_READ) / CURRENT_READ_GAIN) / CURRENT_READ_RESISTANCE;
    this->i_c = (ADC_ns::readVoltage(PHASE_C_CURRENT_READ) / CURRENT_READ_GAIN) / CURRENT_READ_RESISTANCE;

    Serial.printf("v_a=%f gained_v_a=%f i_a=%f i_b=%f, i_c=%f\n", v_a, gained_v_a, i_a, i_b, i_c);
}