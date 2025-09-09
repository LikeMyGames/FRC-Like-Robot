#include "pico/stdlib.h"
#include "hardware/pwm.h"
#include <cmath>

struct abc_t
{
    float a, b, c;
};

struct ab_t
{
    float alpha, beta;
};

// Global state for slice numbers and wrap value
static uint slice_a, slice_b, slice_c;
static uint16_t pwm_wrap;
static bool enabled = false;
static float angle = 0.0f;

// Init function
void pwm_init_3phase(uint pin_a, uint pin_b, uint pin_c, uint16_t wrap)
{
    pwm_wrap = wrap;

    // Configure pins
    gpio_set_function(pin_a, GPIO_FUNC_PWM);
    gpio_set_function(pin_b, GPIO_FUNC_PWM);
    gpio_set_function(pin_c, GPIO_FUNC_PWM);

    slice_a = pwm_gpio_to_slice_num(pin_a);
    slice_b = pwm_gpio_to_slice_num(pin_b);
    slice_c = pwm_gpio_to_slice_num(pin_c);

    pwm_set_wrap(slice_a, wrap);
    pwm_set_wrap(slice_b, wrap);
    pwm_set_wrap(slice_c, wrap);

    // Start with 50% duty (neutral)
    pwm_set_chan_level(slice_a, pwm_gpio_to_channel(pin_a), wrap / 2);
    pwm_set_chan_level(slice_b, pwm_gpio_to_channel(pin_b), wrap / 2);
    pwm_set_chan_level(slice_c, pwm_gpio_to_channel(pin_c), wrap / 2);

    pwm_set_enabled(slice_a, true);
    pwm_set_enabled(slice_b, true);
    pwm_set_enabled(slice_c, true);
}

// Update function (atomic-ish update)
void pwm_update_3phase(abc_t duty)
{
    // Clamp values
    uint16_t da = duty.a < 0 ? 0 : (duty.a > pwm_wrap ? pwm_wrap : (uint16_t)duty.a);
    uint16_t db = duty.b < 0 ? 0 : (duty.b > pwm_wrap ? pwm_wrap : (uint16_t)duty.b);
    uint16_t dc = duty.c < 0 ? 0 : (duty.c > pwm_wrap ? pwm_wrap : (uint16_t)duty.c);

    // Apply new levels
    pwm_set_chan_level(slice_a, 0, da);
    pwm_set_chan_level(slice_b, 0, db);
    pwm_set_chan_level(slice_c, 0, dc);
}

abc_t svpwm(ab_t v_ab, float vbus, uint16_t pwm_max)
{
    abc_t duty;

    // Normalize by DC bus
    float v_alpha = v_ab.alpha / vbus;
    float v_beta = v_ab.beta / vbus;

    // Inverse Clarke: get three-phase voltages
    float v_a = v_alpha;
    float v_b = -0.5f * v_alpha + (sqrtf(3) / 2.0f) * v_beta;
    float v_c = -0.5f * v_alpha - (sqrtf(3) / 2.0f) * v_beta;

    // Find min/max for re-centering
    float v_min = fminf(fminf(v_a, v_b), v_c);
    float v_max = fmaxf(fmaxf(v_a, v_b), v_c);

    // Normalize to 0..1 range
    duty.a = (v_a - v_min) / (v_max - v_min);
    duty.b = (v_b - v_min) / (v_max - v_min);
    duty.c = (v_c - v_min) / (v_max - v_min);

    // Scale to PWM wrap
    duty.a *= pwm_max;
    duty.b *= pwm_max;
    duty.c *= pwm_max;

    return duty;
}

int main()
{
    stdio_init_all();

    // Example pins for motor phases
    const uint PHASE_A = 0;
    const uint PHASE_B = 1;
    const uint PHASE_C = 2;

    const uint16_t PWM_WRAP = 4095; // 12-bit resolution
    const float VBUS = 12.0f;       // Example DC bus voltage

    pwm_init_3phase(PHASE_A, PHASE_B, PHASE_C, PWM_WRAP);

    while (true)
    {
        if (!enabled)
        {
            continue;
        }
        else
        {
            angle += 0.01f;
            if (angle > 2 * M_PI)
                angle = 0;
        }
        // Example rotating vector (simulate FOC voltage command)

        ab_t v_ab;
        v_ab.alpha = cosf(angle) * 6.0f; // request 6V amplitude
        v_ab.beta = sinf(angle) * 6.0f;

        // Compute SVPWM duties
        abc_t duties = svpwm(v_ab, VBUS, PWM_WRAP);

        // Update hardware
        pwm_update_3phase(duties);

        sleep_us(100); // ~10 kHz update rate
    }
}
