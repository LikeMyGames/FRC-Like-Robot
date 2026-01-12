#include <foc_math.h>
#include <util_math.h>
#include <math.h>

Foc::Foc(foc_config_t config)
{
    Ts = config.Ts;
    Serial.println("added foc config");

    pid_i_d = new Pid(&(id_target), config.kp_d, config.ki_d, config.kd_d, config.Ts, 0, 0);
    Serial.println("created i_d pid controller");
    pid_i_q = new Pid(&(iq_target), config.kp_q, config.ki_q, config.kd_q, config.Ts, 0, 0);
    Serial.println("created i_q pid controller");
}

void Foc::Drive(float theta, motor_running_mode running_mode)
{
    // motor_phase_currents_t abc = state->readPhaseCurrents();

    // clark transform
    float i_alpha = ((2.f * i_a) - (i_b - i_c)) / 3.f;
    float i_beta = (TWO_BY_SQRT3 * (i_b - i_c));

    i_alpha = i_alpha;
    i_beta = i_beta;

    float cos_theta = cosf(theta);
    float sin_theta = sinf(theta);

    phase = theta;
    phase_cos = cos_theta;
    phase_sin = sin_theta;

    // park transform
    float i_d = (i_alpha * cos_theta) + (i_beta * sin_theta);
    float i_q = (i_beta * cos_theta) - (i_alpha * sin_theta);

    // d and q pid updates
    pid_i_d->Update(i_d);
    pid_i_q->Update(i_q);

    // inv_park transform
    float v_alpha = (i_d * cos_theta) - (i_q * sin_theta);
    float v_beta = (i_q * cos_theta) + (i_d * sin_theta);
    this->Svpwm(v_alpha, v_beta, (float)1, Ts, &(dA), &(dB), &(dA), &(svm_sector));
}

/**
 * @brief svm Space vector modulation. Magnitude must not be larger than sqrt(3)/2, or 0.866 to avoid overmodulation.
 *        See https://github.com/vedderb/bldc/pull/372#issuecomment-962499623 for a full description.
 * @param alpha voltage
 * @param beta Park transformed and normalized voltage
 * @param PWMFullDutyCycle is the peak value of the PWM counter.
 * @param tAout PWM duty cycle phase A (0 = off all of the time, PWMFullDutyCycle = on all of the time)
 * @param tBout PWM duty cycle phase B
 * @param tCout PWM duty cycle phase C
 */
void Foc::Svpwm(float alpha, float beta, float max_mod, uint32_t PWMFullDutyCycle,
                float *tAout, float *tBout, float *tCout, uint32_t *svm_sector)
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
    *svm_sector = sector;
}