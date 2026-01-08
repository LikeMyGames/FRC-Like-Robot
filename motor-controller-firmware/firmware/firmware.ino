

#include <FlexCAN_T4.h>
#include <EEPROM.h>
#include <TeensyThreads.h>
#include <string>

// std::vector<uint8_t> Motor1_Iq_read();
// bool Motor1_Iq_write(uint8_t buf[7]);
// std::vector<uint8_t> Motor2_Iq_read();
// bool Motor2_Iq_write(uint8_t buf[7]);
// std::vector<uint8_t> Motor1_mode_read();
// bool Motor1_mode_write(uint8_t buf[7]);
// std::vector<uint8_t> Motor2_mode_read();
// bool Motor2_mode_write(uint8_t buf[7]);
// std::vector<uint8_t> Motor1_target_pos_read();
// bool Motor1_target_pos_write(uint8_t buf[7]);
// std::vector<uint8_t> Motor2_target_pos_read();
// bool Motor2_target_pos_write(uint8_t buf[7]);
// std::vector<uint8_t> Motor1_target_vel_read();
// bool Motor1_target_vel_write(uint8_t buf[7]);
// std::vector<uint8_t> Motor2_target_vel_read();
// bool Motor2_target_vel_write(uint8_t buf[7]);

#define NUM_TX_MAILBOXES 2
#define NUM_RX_MAILBOXES 6

typedef struct
{
    std::string name;
    std::function<std::vector<uint8_t>()> read;
    std::function<bool(uint8_t[])> write;
} CAN_REGISTER_t;

FlexCAN_T4<CAN2, RX_SIZE_256, TX_SIZE_16> Can0;
std::unordered_map<uint8_t, CAN_REGISTER_t *> REG_MAP;
uint32_t CAN_ID = 4;
int LOOP_THREAD_ID;
int TIMEOUT_THREAD_ID;
volatile uint16_t TIMEOUT_MS_COUNT = 0;
volatile static uint16_t MAX_TIMEOUT_MS = 2000;

/// @brief
void initCan()
{
    CAN_ID = EEPROM.read(CAN_ID_ADDRESS);
    MAX_TIMEOUT_MS = EEPROM.read(CAN_MAX_TIMEOUT_MS_ADDRESS);
    delay(400);
    Can0.begin();
    Can0.setBaudRate(1 * 1000 * 1000);
    Can0.setMaxMB(NUM_TX_MAILBOXES + NUM_RX_MAILBOXES);
    for (int i = 0; i < NUM_RX_MAILBOXES; i++)
    {
        Can0.setMB((FLEXCAN_MAILBOX)i, RX, EXT);
    }
    for (int i = NUM_RX_MAILBOXES; i < (NUM_TX_MAILBOXES + NUM_RX_MAILBOXES); i++)
    {
        Can0.setMB((FLEXCAN_MAILBOX)i, TX, EXT);
    }
    Can0.setMBFilter(REJECT_ALL);
    Can0.enableMBInterrupts();
    Can0.onReceive(MB0, canSniff);
    Can0.onReceive(MB3, ProcessIdSpecificMessage);
    Can0.setMBUserFilter(MB3, CAN_ID, 0xff);
    Can0.setMBUserFilter(MB0, 0x00, 0x00);
    Can0.enhanceFilter(MB0);
    Can0.enhanceFilter(MB3);
    Can0.mailboxStatus();
    Can0.distribute();
    Serial.println(getID());
    LOOP_THREAD_ID = threads.addThread(canLoop);
    TIMEOUT_THREAD_ID = threads.addThread(checkTimeout);
    Serial.println(LOOP_THREAD_ID);
    Serial.println(TIMEOUT_THREAD_ID);
}

std::vector<uint8_t> CanEchoRead()
{
    return (std::vector<uint8_t>){0};
}

bool CanEchoWrite(uint8_t args[])
{
    return true;
}

/// @brief
/// @param msg
void canSniff(const CAN_message_t &msg)
{
    digitalToggle(LED_BUILTIN);
    Serial.print("MB ");
    Serial.print(msg.mb);
    Serial.print("  OVERRUN: ");
    Serial.print(msg.flags.overrun);
    Serial.print("  LEN: ");
    Serial.print(msg.len);
    Serial.print(" EXT: ");
    Serial.print(msg.flags.extended);
    Serial.print(" TS: ");
    Serial.print(msg.timestamp);
    Serial.print(" ID: ");
    Serial.print(msg.id, HEX);
    Serial.print(" Buffer: ");
    for (uint8_t i = 0; i < msg.len; i++)
    {
        Serial.print(msg.buf[i], HEX);
        Serial.print(" ");
    }
    Serial.println();
}

void ProcessIdSpecificMessage(const CAN_message_t &msg)
{
    resetTimeout();
    Serial.println("Got ID specific messasge");

    if (msg.len == 0)
    {
        return;
    }

    CAN_message_t rmsg;
    rmsg.id = msg.id;
    rmsg.buf[0] = msg.buf[0];
    rmsg.flags.extended = true;

    if (msg.len == 1)
    {

        if (REG_MAP.count(msg.buf[0]) == 0)
        {
            Serial.print("No register exists with id ");
            Serial.println(msg.buf[0]);
        }

        auto reg = REG_MAP[msg.buf[0]];

        std::vector<uint8_t> buf = reg->read();
        if (buf.size() > 7)
        {
            uint size = buf.size();
            size %= 7;
            uint bytesRead = 0;
            rmsg.len = 8;
            while (buf.size() - bytesRead > size)
            {
                for (int i = 0; i < 7; i++)
                {
                    rmsg.buf[i + 1] = buf.at(bytesRead + i);
                    bytesRead++;
                }

                Can0.write(rmsg);
            }

            if (size > 0)
            {
                rmsg.len = size + 1;
                for (uint i = 0; i < size; i++)
                {
                    rmsg.buf[i + 1] = buf.at(bytesRead + i);
                }
                Can0.write(rmsg);
            }
        }
        else
        {
            rmsg.len = buf.size() + 1;
            for (uint i = 1; i < buf.size(); i++)
            {
                rmsg.buf[i] = buf.at(i - 1);
            }
            Can0.write(rmsg);
        }
        return;
    }

    rmsg.len = 2;
    uint8_t funcBuf[7];
    for (uint i = 0; i < sizeof(funcBuf); i++)
    {
        funcBuf[i] = msg.buf[i + 1];
    }
    auto reg = REG_MAP[msg.buf[0]];
    rmsg.buf[1] = reg->write(funcBuf);

    Can0.write(rmsg);
}

void addCanRegister(uint8_t id, std::string name, std::function<std::vector<uint8_t>()> readAction, std::function<bool(uint8_t[])> writeAction)
{
    if (REG_MAP.count(id) != 0)
    {
        Serial.println("Could not add register, already exists");
    }

    CAN_REGISTER_t *reg = {};
    reg->name = name;
    reg->read = readAction;
    reg->write = writeAction;

    REG_MAP.insert_or_assign(id, reg);
}

void canLoop()
{
    while (1)
    {
        Can0.events();
        if (Can0.getRXQueueCount() != 0)
        {
            Serial.println(Can0.getRXQueueCount());
        }
    }
}

void Status()
{
    Can0.mailboxStatus();
}

uint32_t getID()
{
    return CAN_ID;
}

void setID(uint32_t newID)
{
    Can0.FLEXCAN_EnterFreezeMode();
    Can0.setMBUserFilter(MB3, newID, 0xff);
    Can0.FLEXCAN_ExitFreezeMode();
    CAN_ID = newID;
}

void resetTimeout()
{
    TIMEOUT_MS_COUNT = 0;
}

void checkTimeout()
{
    while (1)
    {
        if (TIMEOUT_MS_COUNT >= MAX_TIMEOUT_MS)
        {
            disableAllMotors();
            threads.delay(500);
        }
        TIMEOUT_MS_COUNT++;
        threads.delay(1);
    }
}

uint16_t getMaxCanTimeout()
{
    return TIMEOUT_MS_COUNT;
}

void setMaxCanTimeout(uint16_t newTimeout)
{
    MAX_TIMEOUT_MS = newTimeout;
}

foc_state_t *init_foc(foc_config_t *config)
{
    foc_state_t *state = {};
    state->Ts = config->Ts;

    state->pid_i_d = NewPidController(&(state->id_target), config->kp_d, config->ki_d, config->kd_d, config->Ts);
    state->pid_i_q = NewPidController(&(state->iq_target), config->kp_q, config->ki_q, config->kd_q, config->Ts);

    return state;
}

foc_phase_duty_timings_t foc_drive(foc_state_t *state, float theta)
{
    // motor_phase_currents_t abc = state->readPhaseCurrents();

    // clark transform
    float i_alpha = ((2.f * state->i_a) - (state->i_b - state->i_c)) / 3.f;
    float i_beta = (TWO_BY_SQRT3 * (state->i_b - state->i_c));

    state->i_alpha = i_alpha;
    state->i_beta = i_beta;

    float cos_theta = cosf(theta);
    float sin_theta = sinf(theta);

    state->phase = theta;
    state->phase_cos = cos_theta;
    state->phase_sin = sin_theta;

    // park transform
    float i_d = (i_alpha * cos_theta) + (i_beta * sin_theta);
    float i_q = (i_beta * cos_theta) - (i_alpha * sin_theta);

    // d and q pid updates
    UpdatePid(state->pid_i_d, i_d);
    UpdatePid(state->pid_i_q, i_q);

    // inv_park transform
    float v_alpha = (i_d * cos_theta) - (i_q * sin_theta);
    float v_beta = (i_q * cos_theta) + (i_d * sin_theta);
    foc_svpwm(v_alpha, v_beta, (float)1, state->Ts, &(state->dA), &(state->dB), &(state->dA), &(state->svm_sector));
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
void foc_svpwm(float alpha, float beta, float max_mod, uint32_t PWMFullDutyCycle,
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

std::vector<uint8_t> utils_float_to_bytes(float f)
{
    // Ensure the float size matches the byte array size
    static_assert(sizeof(float) == 4, "float is not 4 bytes");
    static_assert(sizeof(uint8_t) == 1, "uint8_t is not 1 byte");
    std::vector<uint8_t> bytes;
    std::memcpy(&bytes, &f, sizeof(float));
    return bytes;
}

float utils_bytes_to_float(uint8_t bytes[4])
{
    float f;
    static_assert(sizeof(float) == 4, "float is not 4 bytes");
    static_assert(sizeof(uint8_t) == 1, "uint8_t is not 1 byte");
    std::memcpy(&f, bytes, sizeof(float));
    return f;
}

float utils_map_range(float val, float inMin, float inMax, float outMin, float outMax)
{
    return outMin + (val - inMin) * ((outMax - outMin) / (inMax - inMin));
}

uint16_t utils_map_range(uint16_t val, uint16_t inMin, uint16_t inMax, uint16_t outMin, uint16_t outMax)
{
    return outMin + (val - inMin) * ((outMax - outMin) / (inMax - inMin));
}

/*
 * Map angle from 0 to 1 in the range min to max. If angle is
 * outside of the range it will be less truncated to the closest
 * angle. Angle units: Degrees
 */
float utils_map_angle(float angle, float min, float max)
{
    if (max == min)
    {
        return -1;
    }

    float range_pos = max - min;
    utils_norm_angle(&range_pos);
    float range_neg = min - max;
    utils_norm_angle(&range_neg);
    float margin = range_neg / 2.0;

    angle -= min;
    utils_norm_angle(&angle);
    if (angle > (360 - margin))
    {
        angle -= 360.0;
    }

    float res = angle / range_pos;
    utils_truncate_number(&res, 0.0, 1.0);

    return res;
}

/**
 * Truncate absolute values less than tres to zero. The value
 * tres will be mapped to 0 and the value max to max.
 */
void utils_deadband(float *value, float tres, float max)
{
    if (fabsf(*value) < tres)
    {
        *value = 0.0;
    }
    else
    {
        float k = max / (max - tres);
        if (*value > 0.0)
        {
            *value = k * *value + max * (1.0 - k);
        }
        else
        {
            *value = -(k * -*value + max * (1.0 - k));
        }
    }
}

/**
 * Takes the average of a number of angles.
 *
 * @param angles
 * The angles in radians.
 *
 * @param angles_num
 * The number of angles.
 *
 * @param weights
 * The weight of the summarized angles
 *
 * @return
 * The average angle.
 */
float utils_avg_angles_rad_fast(float *angles, float *weights, int angles_num)
{
    float s_sum = 0.0;
    float c_sum = 0.0;

    for (int i = 0; i < angles_num; i++)
    {
        float s, c;
        utils_fast_sincos_better(angles[i], &s, &c);
        s_sum += s * weights[i];
        c_sum += c * weights[i];
    }

    return utils_fast_atan2(s_sum, c_sum);
}

/**
 * Interpolate two angles in radians and normalize the result to
 * -pi to pi.
 *
 * @param a1
 * The first angle
 *
 * @param a2
 * The second angle
 *
 * @param weight_a1
 * The weight of the first angle. If this is 1.0 the result will
 * be a1 and if it is 0.0 the result will be a2.
 *
 */
float utils_interpolate_angles_rad(float a1, float a2, float weight_a1)
{
    while ((a1 - a2) > M_PI)
        a2 += 2.0 * M_PI;
    while ((a2 - a1) > M_PI)
        a1 += 2.0 * M_PI;

    float res = a1 * weight_a1 + a2 * (1.0 - weight_a1);
    utils_norm_angle_rad(&res);
    return res;
}

/**
 * Get the middle value of three values
 *
 * @param a
 * First value
 *
 * @param b
 * Second value
 *
 * @param c
 * Third value
 *
 * @return
 * The middle value
 */
float utils_middle_of_3(float a, float b, float c)
{
    float middle;

    if ((a <= b) && (a <= c))
    {
        middle = (b <= c) ? b : c;
    }
    else if ((b <= a) && (b <= c))
    {
        middle = (a <= c) ? a : c;
    }
    else
    {
        middle = (a <= b) ? a : b;
    }
    return middle;
}

/**
 * Get the middle value of three values
 *
 * @param a
 * First value
 *
 * @param b
 * Second value
 *
 * @param c
 * Third value
 *
 * @return
 * The middle value
 */
int utils_middle_of_3_int(int a, int b, int c)
{
    int middle;

    if ((a <= b) && (a <= c))
    {
        middle = (b <= c) ? b : c;
    }
    else if ((b <= a) && (b <= c))
    {
        middle = (a <= c) ? a : c;
    }
    else
    {
        middle = (a <= b) ? a : b;
    }
    return middle;
}

/**
 * Fast atan2
 *
 * See http://www.dspguru.com/dsp/tricks/fixed-point-atan2-with-self-normalization
 *
 * @param y
 * y
 *
 * @param x
 * x
 *
 * @return
 * The angle in radians
 */
float utils_fast_atan2(float y, float x)
{
    float abs_y = fabsf(y) + 1e-20; // kludge to prevent 0/0 condition
    float angle;

    if (x >= 0)
    {
        float r = (x - abs_y) / (x + abs_y);
        float rsq = r * r;
        angle = ((0.1963 * rsq) - 0.9817) * r + (M_PI / 4.0);
    }
    else
    {
        float r = (x + abs_y) / (abs_y - x);
        float rsq = r * r;
        angle = ((0.1963 * rsq) - 0.9817) * r + (3.0 * M_PI / 4.0);
    }

    UTILS_NAN_ZERO(angle);

    if (y < 0)
    {
        return (-angle);
    }
    else
    {
        return (angle);
    }
}

float utils_fast_sin(float angle)
{
    while (angle < -M_PI)
    {
        angle += 2.0 * M_PI;
    }

    while (angle > M_PI)
    {
        angle -= 2.0 * M_PI;
    }

    float res = 0.0;

    if (angle < 0.0)
    {
        res = 1.27323954 * angle + 0.405284735 * angle * angle;
    }
    else
    {
        res = 1.27323954 * angle - 0.405284735 * angle * angle;
    }

    return res;
}

float utils_fast_cos(float angle)
{
    angle += 0.5 * M_PI;

    while (angle < -M_PI)
    {
        angle += 2.0 * M_PI;
    }

    while (angle > M_PI)
    {
        angle -= 2.0 * M_PI;
    }

    float res = 0.0;

    if (angle < 0.0)
    {
        res = 1.27323954 * angle + 0.405284735 * angle * angle;
    }
    else
    {
        res = 1.27323954 * angle - 0.405284735 * angle * angle;
    }

    return res;
}

/**
 * Fast sine and cosine implementation.
 *
 * See http://lab.polygonal.de/?p=205
 *
 * @param angle
 * The angle in radians
 * WARNING: Don't use too large angles.
 *
 * @param sin
 * A pointer to store the sine value.
 *
 * @param cos
 * A pointer to store the cosine value.
 */
void utils_fast_sincos(float angle, float *sin, float *cos)
{
    // always wrap input angle to -PI..PI
    while (angle < -M_PI)
    {
        angle += 2.0 * M_PI;
    }

    while (angle > M_PI)
    {
        angle -= 2.0 * M_PI;
    }

    // compute sine
    if (angle < 0.0)
    {
        *sin = 1.27323954 * angle + 0.405284735 * angle * angle;
    }
    else
    {
        *sin = 1.27323954 * angle - 0.405284735 * angle * angle;
    }

    // compute cosine: sin(x + PI/2) = cos(x)
    angle += 0.5 * M_PI;

    if (angle > M_PI)
    {
        angle -= 2.0 * M_PI;
    }

    if (angle < 0.0)
    {
        *cos = 1.27323954 * angle + 0.405284735 * angle * angle;
    }
    else
    {
        *cos = 1.27323954 * angle - 0.405284735 * angle * angle;
    }
}

/**
 * Fast sine and cosine implementation.
 *
 * See http://lab.polygonal.de/?p=205
 *
 * @param angle
 * The angle in radians
 * WARNING: Don't use too large angles.
 *
 * @param sin
 * A pointer to store the sine value.
 *
 * @param cos
 * A pointer to store the cosine value.
 */
void utils_fast_sincos_better(float angle, float *sin, float *cos)
{
    // always wrap input angle to -PI..PI
    while (angle < -M_PI)
    {
        angle += 2.0 * M_PI;
    }

    while (angle > M_PI)
    {
        angle -= 2.0 * M_PI;
    }

    // compute sine
    if (angle < 0.0)
    {
        *sin = 1.27323954 * angle + 0.405284735 * angle * angle;

        if (*sin < 0.0)
        {
            *sin = 0.225 * (*sin * -*sin - *sin) + *sin;
        }
        else
        {
            *sin = 0.225 * (*sin * *sin - *sin) + *sin;
        }
    }
    else
    {
        *sin = 1.27323954 * angle - 0.405284735 * angle * angle;

        if (*sin < 0.0)
        {
            *sin = 0.225 * (*sin * -*sin - *sin) + *sin;
        }
        else
        {
            *sin = 0.225 * (*sin * *sin - *sin) + *sin;
        }
    }

    // compute cosine: sin(x + PI/2) = cos(x)
    angle += 0.5 * M_PI;
    if (angle > M_PI)
    {
        angle -= 2.0 * M_PI;
    }

    if (angle < 0.0)
    {
        *cos = 1.27323954 * angle + 0.405284735 * angle * angle;

        if (*cos < 0.0)
        {
            *cos = 0.225 * (*cos * -*cos - *cos) + *cos;
        }
        else
        {
            *cos = 0.225 * (*cos * *cos - *cos) + *cos;
        }
    }
    else
    {
        *cos = 1.27323954 * angle - 0.405284735 * angle * angle;

        if (*cos < 0.0)
        {
            *cos = 0.225 * (*cos * -*cos - *cos) + *cos;
        }
        else
        {
            *cos = 0.225 * (*cos * *cos - *cos) + *cos;
        }
    }
}

/**
 * Calculate the values with the lowest magnitude.
 *
 * @param va
 * The first value.
 *
 * @param vb
 * The second value.
 *
 * @return
 * The value with the lowest magnitude.
 */
float utils_min_abs(float va, float vb)
{
    float res;
    if (fabsf(va) < fabsf(vb))
    {
        res = va;
    }
    else
    {
        res = vb;
    }

    return res;
}

/**
 * Calculate the values with the highest magnitude.
 *
 * @param va
 * The first value.
 *
 * @param vb
 * The second value.
 *
 * @return
 * The value with the highest magnitude.
 */
float utils_max_abs(float va, float vb)
{
    float res;
    if (fabsf(va) > fabsf(vb))
    {
        res = va;
    }
    else
    {
        res = vb;
    }

    return res;
}

/**
 * Create string representation of the binary content of a byte
 *
 * @param x
 * The byte.
 *
 * @param b
 * Array to store the string representation in.
 */
void utils_byte_to_binary(int x, char *b)
{
    b[0] = '\0';

    int z;
    for (z = 128; z > 0; z >>= 1)
    {
        strcat(b, ((x & z) == z) ? "1" : "0");
    }
}

float utils_throttle_curve(float val, float curve_acc, float curve_brake, int mode)
{
    float ret = 0.0;

    if (val < -1.0)
    {
        val = -1.0;
    }

    if (val > 1.0)
    {
        val = 1.0;
    }

    float val_a = fabsf(val);

    float curve;
    if (val >= 0.0)
    {
        curve = curve_acc;
    }
    else
    {
        curve = curve_brake;
    }

    // See
    // http://math.stackexchange.com/questions/297768/how-would-i-create-a-exponential-ramp-function-from-0-0-to-1-1-with-a-single-val
    if (mode == 0)
    { // Exponential
        if (curve >= 0.0)
        {
            ret = 1.0 - powf(1.0 - val_a, 1.0 + curve);
        }
        else
        {
            ret = powf(val_a, 1.0 - curve);
        }
    }
    else if (mode == 1)
    { // Natural
        if (fabsf(curve) < 1e-10)
        {
            ret = val_a;
        }
        else
        {
            if (curve >= 0.0)
            {
                ret = 1.0 - ((expf(curve * (1.0 - val_a)) - 1.0) / (expf(curve) - 1.0));
            }
            else
            {
                ret = (expf(-curve * val_a) - 1.0) / (expf(-curve) - 1.0);
            }
        }
    }
    else if (mode == 2)
    { // Polynomial
        if (curve >= 0.0)
        {
            ret = 1.0 - ((1.0 - val_a) / (1.0 + curve * val_a));
        }
        else
        {
            ret = val_a / (1.0 - curve * (1.0 - val_a));
        }
    }
    else
    { // Linear
        ret = val_a;
    }

    if (val < 0.0)
    {
        ret = -ret;
    }

    return ret;
}

uint32_t utils_crc32c(uint8_t *data, uint32_t len)
{
    uint32_t crc = 0xFFFFFFFF;

    for (uint32_t i = 0; i < len; i++)
    {
        uint32_t byte = data[i];
        crc = crc ^ byte;

        for (int j = 7; j >= 0; j--)
        {
            uint32_t mask = -(crc & 1);
            crc = (crc >> 1) ^ (0x82F63B78 & mask);
        }
    }

    return ~crc;
}

// Yes, this is only the average...
void utils_fft32_bin0(float *real_in, float *real, float *imag)
{
    *real = 0.0;
    *imag = 0.0;

    for (int i = 0; i < 32; i++)
    {
        *real += real_in[i];
    }

    *real /= 32.0;
}

void utils_fft32_bin1(float *real_in, float *real, float *imag)
{
    *real = 0.0;
    *imag = 0.0;
    for (int i = 0; i < 32; i++)
    {
        *real += real_in[i] * utils_tab_cos_32_1[i];
        *imag -= real_in[i] * utils_tab_sin_32_1[i];
    }
    *real /= 32.0;
    *imag /= 32.0;
}

void utils_fft32_bin2(float *real_in, float *real, float *imag)
{
    *real = 0.0;
    *imag = 0.0;
    for (int i = 0; i < 32; i++)
    {
        *real += real_in[i] * utils_tab_cos_32_2[i];
        *imag -= real_in[i] * utils_tab_sin_32_2[i];
    }
    *real /= 32.0;
    *imag /= 32.0;
}

void utils_fft16_bin0(float *real_in, float *real, float *imag)
{
    *real = 0.0;
    *imag = 0.0;

    for (int i = 0; i < 16; i++)
    {
        *real += real_in[i];
    }

    *real /= 16.0;
}

void utils_fft16_bin1(float *real_in, float *real, float *imag)
{
    *real = 0.0;
    *imag = 0.0;
    for (int i = 0; i < 16; i++)
    {
        *real += real_in[i] * utils_tab_cos_32_1[2 * i];
        *imag -= real_in[i] * utils_tab_sin_32_1[2 * i];
    }
    *real /= 16.0;
    *imag /= 16.0;
}

void utils_fft16_bin2(float *real_in, float *real, float *imag)
{
    *real = 0.0;
    *imag = 0.0;
    for (int i = 0; i < 16; i++)
    {
        *real += real_in[i] * utils_tab_cos_32_2[2 * i];
        *imag -= real_in[i] * utils_tab_sin_32_2[2 * i];
    }
    *real /= 16.0;
    *imag /= 16.0;
}

void utils_fft8_bin0(float *real_in, float *real, float *imag)
{
    *real = 0.0;
    *imag = 0.0;

    for (int i = 0; i < 8; i++)
    {
        *real += real_in[i];
    }

    *real /= 8.0;
}

void utils_fft8_bin1(float *real_in, float *real, float *imag)
{
    *real = 0.0;
    *imag = 0.0;
    for (int i = 0; i < 8; i++)
    {
        *real += real_in[i] * utils_tab_cos_32_1[4 * i];
        *imag -= real_in[i] * utils_tab_sin_32_1[4 * i];
    }
    *real /= 8.0;
    *imag /= 8.0;
}

void utils_fft8_bin2(float *real_in, float *real, float *imag)
{
    *real = 0.0;
    *imag = 0.0;
    for (int i = 0; i < 8; i++)
    {
        *real += real_in[i] * utils_tab_cos_32_2[4 * i];
        *imag -= real_in[i] * utils_tab_sin_32_2[4 * i];
    }
    *real /= 8.0;
    *imag /= 8.0;
}

// A mapping of a samsung 30q cell for % remaining capacity vs. voltage from
// 4.2 to 3.2, note that the you lose 15% of the 3Ah rated capacity in this range
float utils_batt_liion_norm_v_to_capacity(float norm_v)
{
    // constants for polynomial fit of lithium ion battery
    const float li_p[] = {
        -2.979767, 5.487810, -3.501286, 1.675683, 0.317147};
    utils_truncate_number(&norm_v, 0.0, 1.0);
    float v2 = norm_v * norm_v;
    float v3 = v2 * norm_v;
    float v4 = v3 * norm_v;
    float v5 = v4 * norm_v;
    float capacity = li_p[0] * v5 + li_p[1] * v4 + li_p[2] * v3 +
                     li_p[3] * v2 + li_p[4] * norm_v;
    return capacity;
}

static int uint16_cmp_func(const void *a, const void *b)
{
    return (*(uint16_t *)a - *(uint16_t *)b);
}

uint16_t utils_median_filter_uint16_run(uint16_t *buffer,
                                        unsigned int *buffer_index, unsigned int filter_len, uint16_t sample)
{
    buffer[(*buffer_index)++] = sample;
    *buffer_index %= filter_len;
    uint16_t buffer_sorted[filter_len]; // Assume we have enough stack space
    memcpy(buffer_sorted, buffer, sizeof(uint16_t) * filter_len);
    qsort(buffer_sorted, filter_len, sizeof(uint16_t), uint16_cmp_func);
    return buffer_sorted[filter_len / 2];
}

void utils_rotate_vector3(float *input, float *rotation, float *output, bool reverse)
{
    float s1, c1, s2, c2, s3, c3;

    if (rotation[2] != 0.0)
    {
        s1 = sinf(rotation[2]);
        c1 = cosf(rotation[2]);
    }
    else
    {
        s1 = 0.0;
        c1 = 1.0;
    }

    if (rotation[1] != 0.0)
    {
        s2 = sinf(rotation[1]);
        c2 = cosf(rotation[1]);
    }
    else
    {
        s2 = 0.0;
        c2 = 1.0;
    }

    if (rotation[0] != 0.0)
    {
        s3 = sinf(rotation[0]);
        c3 = cosf(rotation[0]);
    }
    else
    {
        s3 = 0.0;
        c3 = 1.0;
    }

    float m11 = c1 * c2;
    float m12 = c1 * s2 * s3 - c3 * s1;
    float m13 = s1 * s3 + c1 * c3 * s2;
    float m21 = c2 * s1;
    float m22 = c1 * c3 + s1 * s2 * s3;
    float m23 = c3 * s1 * s2 - c1 * s3;
    float m31 = -s2;
    float m32 = c2 * s3;
    float m33 = c2 * c3;

    if (reverse)
    {
        output[0] = input[0] * m11 + input[1] * m21 + input[2] * m31;
        output[1] = input[0] * m12 + input[1] * m22 + input[2] * m32;
        output[2] = input[0] * m13 + input[1] * m23 + input[2] * m33;
    }
    else
    {
        output[0] = input[0] * m11 + input[1] * m12 + input[2] * m13;
        output[1] = input[0] * m21 + input[1] * m22 + input[2] * m23;
        output[2] = input[0] * m31 + input[1] * m32 + input[2] * m33;
    }
}

const float utils_tab_sin_32_1[] = {
    0.000000, 0.195090, 0.382683, 0.555570, 0.707107, 0.831470, 0.923880, 0.980785,
    1.000000, 0.980785, 0.923880, 0.831470, 0.707107, 0.555570, 0.382683, 0.195090,
    0.000000, -0.195090, -0.382683, -0.555570, -0.707107, -0.831470, -0.923880, -0.980785,
    -1.000000, -0.980785, -0.923880, -0.831470, -0.707107, -0.555570, -0.382683, -0.195090};

const float utils_tab_sin_32_2[] = {
    0.000000, 0.382683, 0.707107, 0.923880, 1.000000, 0.923880, 0.707107, 0.382683,
    0.000000, -0.382683, -0.707107, -0.923880, -1.000000, -0.923880, -0.707107, -0.382683,
    -0.000000, 0.382683, 0.707107, 0.923880, 1.000000, 0.923880, 0.707107, 0.382683,
    0.000000, -0.382683, -0.707107, -0.923880, -1.000000, -0.923880, -0.707107, -0.382683};

const float utils_tab_cos_32_1[] = {
    1.000000, 0.980785, 0.923880, 0.831470, 0.707107, 0.555570, 0.382683, 0.195090,
    0.000000, -0.195090, -0.382683, -0.555570, -0.707107, -0.831470, -0.923880, -0.980785,
    -1.000000, -0.980785, -0.923880, -0.831470, -0.707107, -0.555570, -0.382683, -0.195090,
    -0.000000, 0.195090, 0.382683, 0.555570, 0.707107, 0.831470, 0.923880, 0.980785};

const float utils_tab_cos_32_2[] = {
    1.000000, 0.923880, 0.707107, 0.382683, 0.000000, -0.382683, -0.707107, -0.923880,
    -1.000000, -0.923880, -0.707107, -0.382683, -0.000000, 0.382683, 0.707107, 0.923880,
    1.000000, 0.923880, 0.707107, 0.382683, 0.000000, -0.382683, -0.707107, -0.923880,
    -1.000000, -0.923880, -0.707107, -0.382683, -0.000000, 0.382683, 0.707107, 0.923880};

static unsigned int NEXT_MOTOR_ID = 0;

motor_state_t *NewMotor(motor_config_t *config)
{
    motor_state_t *state = {};
    state->config = config;
    state->internal_encoder = new Encoder(config->internal_encoder_pin);
    state->external_encoder = new Encoder(config->external_encoder_pin);
    // state->driver = NewMotorDriver();
    state->foc_state = init_foc(state->foc_config);
    MOTOR_STATE_MAP.insert_or_assign(NEXT_MOTOR_ID, state);
    NEXT_MOTOR_ID++;
    return state;
}

void InitMotor(motor_state_t *motor, motor_driver_t *driver)
{
    motor->driver = driver;
    motor_driver_ns::InitMotorDriver(driver);
}

// disabling the motor
void disableMotor(motor_state_t *state)
{
    state->status = MOTOR_DISABLED;
    if (state->MOTOR_THREAD_STARTED)
    {
        threads.suspend(state->MOTOR_THREAD_ID);
    }
}

void disableMotor(motor_state_t *state, motor_error error)
{
    disableMotor(state);
    state->error = error;
}

void disableMotor(unsigned int id)
{
    disableMotor(MOTOR_STATE_MAP.at(id));
}

void disableMotor(unsigned int id, motor_error error)
{
    disableMotor(MOTOR_STATE_MAP.at(id), error);
}

void disableAllMotors()
{
    Serial.println("Disabling all motors on controller");
    for (auto pair : MOTOR_STATE_MAP)
    {
        disableMotor(pair.second);
    }
}

void disableAllMotors(motor_error error)
{
    for (auto pair : MOTOR_STATE_MAP)
    {
        disableMotor(pair.second);
    }
}

// enabling the motor
void enableMotor(motor_state_t *state)
{
    if (!state->error)
    {
        state->status = MOTOR_ENABLED;
        if (!state->MOTOR_THREAD_STARTED)
        {
            threads.restart(state->MOTOR_THREAD_ID);
            state->MOTOR_THREAD_STARTED = true;
        }
    }
}

void enableMotor(unsigned int id)
{
    enableMotor(MOTOR_STATE_MAP.at(id));
}

void enableAllMotors()
{
    for (auto pair : MOTOR_STATE_MAP)
    {
        enableMotor(pair.second);
    }
}

// Main motor loop
void MotorLoop(int MOTOR_ID)
{
    motor_state_t *state = MOTOR_STATE_MAP[MOTOR_ID];
    foc_drive(state->foc_state, state->internal_encoder->ReadRad());
    DriveMotorByPercent(state->driver, state->foc_state->dA, state->foc_state->dB, state->foc_state->dC);
}

void MotorLoop(motor_state_t *state)
{
    foc_drive(state->foc_state, state->internal_encoder->ReadRad());
    DriveMotorByPercent(state->driver, state->foc_state->dA, state->foc_state->dB, state->foc_state->dC);
}

// static unsigned int DRIVERS_CREATED = 0;
const float DeadTimeNs = 50.0;

namespace motor_driver_ns
{
    motor_driver_t motor_driver_1 = {
        {Sm13},
        {Sm20},
        {Sm22},
        {Tm1, Tm2},
    };

    motor_driver_t motor_driver_2 = {
        {Sm31},
        {Sm40, Sm41},
        {Sm42},
        {Tm3, Tm4},
    };
}

void motor_driver_ns::InitMotorDriver(motor_driver_t *driver)
{
    eFlex::Config myConfig;
    myConfig.setReloadLogic(kPWM_ReloadPwmFullCycle);
    myConfig.setClockSource(kPWM_Submodule0Clock);
    myConfig.setPrescale(kPWM_Prescale_Divide_1);
    myConfig.setInitializationControl(kPWM_Initialize_MasterSync);
    myConfig.setPwmFreqHz(PWM_FREQ);

    // SubModule config
    if (driver->phaseA_submodules.size() == 1)
    {
        myConfig.setPairOperation(kPWM_ComplementaryPwmA);
        if (driver->phaseA_submodules[0].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
    }
    else
    {
        myConfig.setPairOperation(kPWM_Independent);
        driver->phaseA_submodules[0].setupLevel(kPWM_HighTrue);
        if (driver->phaseA_submodules[0].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
        driver->phaseA_submodules[1].setupLevel(kPWM_LowTrue);
        if (driver->phaseA_submodules[1].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
    }

    if (driver->phaseB_submodules.size() == 1)
    {
        myConfig.setPairOperation(kPWM_ComplementaryPwmA);
        if (driver->phaseB_submodules[0].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
    }
    else
    {
        myConfig.setPairOperation(kPWM_Independent);
        driver->phaseB_submodules[0].setupLevel(kPWM_HighTrue);
        if (driver->phaseB_submodules[0].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
        driver->phaseB_submodules[1].setupLevel(kPWM_LowTrue);
        if (driver->phaseB_submodules[1].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
    }

    if (driver->phaseC_submodules.size() == 1)
    {
        myConfig.setPairOperation(kPWM_ComplementaryPwmA);
        if (driver->phaseC_submodules[0].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
    }
    else
    {
        myConfig.setPairOperation(kPWM_Independent);
        driver->phaseC_submodules[0].setupLevel(kPWM_HighTrue);
        if (driver->phaseC_submodules[0].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
        driver->phaseC_submodules[1].setupLevel(kPWM_LowTrue);
        if (driver->phaseC_submodules[1].configure(myConfig) != true)
        {
            Serial.println("eFlexPwm Submodule initialization failed");
            exit(EXIT_FAILURE);
        }
    }

    // Timer config
    uint16_t deadTimeVal = ((uint64_t)driver->timers[0].srcClockHz() * DeadTimeNs) / 1000000000;
    driver->timers[0].setupDeadtime(deadTimeVal);

    deadTimeVal = ((uint64_t)driver->timers[1].srcClockHz() * DeadTimeNs) / 1000000000;
    driver->timers[1].setupDeadtime(deadTimeVal);

    // Start timers
    bool timer1Success = driver->timers[0].begin();
    bool timer2Success = driver->timers[1].begin();

    if (timer1Success || timer2Success)
    {
        Serial.println("failed to start eFlexPwm submodules");
        exit(EXIT_FAILURE);
    }
    else
    {
        Serial.println("eFlexPwm submodules successfully started");
    }
}

void DriveMotorByPercent(motor_driver_t *driver, float dA, float dB, float dC)
{
    // Set timing for all phase A submodules
    for (auto submodule : driver->phaseA_submodules)
    {
        submodule.updateDutyCyclePercent(dA, eFlex::ChanA);
    }

    // Set timing for all phase B submodules
    for (auto submodule : driver->phaseB_submodules)
    {
        submodule.updateDutyCyclePercent(dB, eFlex::ChanA);
    }

    // Set timing for all phase C submodules
    for (auto submodule : driver->phaseC_submodules)
    {
        submodule.updateDutyCyclePercent(dC, eFlex::ChanA);
    }
}

// Not implemented
void DriveMotorByTime(motor_driver_t *driver, float dA, float dB, float dC)
{
}

motor_state_t *Motor1;
static motor_config_t motor1_config = {
    1.f,
    21U,
    20U,
};

motor_state_t *Motor2;
static motor_config_t motor2_config = {
    1.f,
    19,
    18,
};

void setup()
{
    digitalToggle(LED_BUILTIN);
    Serial.begin(115200);
    // Serial.setTimeout(0);
    Serial.println("Starting setup");
    initCan();
    Serial.println("finished motor controller setup");
    Motor1 = NewMotor(&motor1_config);
    Motor2 = NewMotor(&motor2_config);
    InitMotor(Motor1, &motor_driver_ns::motor_driver_1);
    InitMotor(Motor2, &motor_driver_ns::motor_driver_2);
    addCanRegister(0, "echo", CanEchoRead, CanEchoWrite);
    addCanRegister(1, "motor1 mode", Motor1_mode_read, Motor1_mode_write);
    addCanRegister(2, "motor2 mode", Motor2_mode_read, Motor2_mode_write);
    addCanRegister(3, "motor1 target iq", Motor1_Iq_read, Motor1_Iq_write);
    addCanRegister(4, "motor2 target iq", Motor2_Iq_read, Motor2_Iq_write);
    addCanRegister(5, "motor1 target pos", Motor1_target_pos_read, Motor1_target_pos_write);
    addCanRegister(6, "motor1 target vel", Motor1_target_vel_read, Motor1_target_vel_write);
    addCanRegister(7, "motor2 target pos", Motor2_target_pos_read, Motor2_target_pos_write);
    addCanRegister(8, "motor2 target vel", Motor2_target_vel_read, Motor2_target_vel_write);
    digitalToggle(LED_BUILTIN);
}

std::vector<uint8_t> Motor1_Iq_read()
{
    return {(uint8_t)(Motor1->foc_state->iq_target)};
}

bool Motor1_Iq_write(uint8_t buf[7])
{
    Motor1->foc_state->iq_target = (float)buf[0];
    return true;
}

std::vector<uint8_t> Motor2_Iq_read()
{
    return {(uint8_t)(Motor2->foc_state->iq_target)};
}

bool Motor2_Iq_write(uint8_t buf[7])
{
    Motor2->foc_state->iq_target = (float)buf[0];
    return true;
}

std::vector<uint8_t> Motor1_mode_read()
{
    return {(uint8_t)(Motor1->foc_state->mode)};
}

bool Motor1_mode_write(uint8_t buf[7])
{
    Motor1->foc_state->mode = (foc_mode)buf[0];
    return true;
}

std::vector<uint8_t> Motor2_mode_read()
{
    return {(u_int8_t)(Motor2->foc_state->mode)};
}

bool Motor2_mode_write(uint8_t buf[7])
{
    Motor2->foc_state->mode = (foc_mode)buf[0];
    return true;
}

std::vector<uint8_t> Motor1_target_pos_read()
{
    return utils_float_to_bytes(Motor1->foc_state->targetAngle);
}

bool Motor1_target_pos_write(uint8_t buf[7])
{
    uint8_t newBuf[4];
    for (int i = 0; i < (int)sizeof(newBuf); i++)
    {
        newBuf[i] = buf[i];
    }
    Motor1->foc_state->targetAngle = utils_bytes_to_float(newBuf);
    return true;
}

std::vector<uint8_t> Motor2_target_pos_read()
{
    return utils_float_to_bytes(Motor2->foc_state->targetAngle);
}

bool Motor2_target_pos_write(uint8_t buf[7])
{
    uint8_t newBuf[4];
    for (int i = 0; i < (int)sizeof(newBuf); i++)
    {
        newBuf[i] = buf[i];
    }
    Motor2->foc_state->targetAngle = utils_bytes_to_float(newBuf);
    return true;
}

std::vector<uint8_t> Motor1_target_vel_read()
{
    return utils_float_to_bytes(Motor1->foc_state->targetVel);
}

bool Motor1_target_vel_write(uint8_t buf[7])
{
    uint8_t newBuf[4];
    for (int i = 0; i < (int)sizeof(newBuf); i++)
    {
        newBuf[i] = buf[i];
    }
    Motor1->foc_state->targetVel = utils_bytes_to_float(newBuf);
    return true;
}

std::vector<uint8_t> Motor2_target_vel_read()
{
    return utils_float_to_bytes(Motor2->foc_state->targetVel);
}

bool Motor2_target_vel_write(uint8_t buf[7])
{
    uint8_t newBuf[4];
    for (int i = 0; i < (int)sizeof(newBuf); i++)
    {
        newBuf[i] = buf[i];
    }
    Motor2->foc_state->targetVel = utils_bytes_to_float(newBuf);
    return true;
}




void loop()
{
    MotorLoop(Motor1);
    MotorLoop(Motor2);
    delayMicroseconds(50);
}