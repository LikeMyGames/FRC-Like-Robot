#include <FlexCAN_T4.h>
#include <string>
#include <unordered_map>
#include <EEPROM.h>
#include <Arduino.h>
#include <eFlexPwm.h>

#define CAN_ID_ADDRESS 0
#define MOTOR_CURRENT_LIMIT_ADDRESS 1
#define PHASE_A_PIN 2
#define PHASE_B_PIN 3
#define PHASE_C_PIN 4
#define NUM_TX_MAILBOXES 2
#define NUM_RX_MAILBOXES 6
#define PWM_FREQ 20000
#define PWM_CLK 150000000.0f
// #define PWM_MOD ((uint16_t)(PWM_CLK / (2.0f * PWM_FREQ)))
#define PWM_PERIOD ((uint16_t)(1000000U / PWM_FREQ)) 
#define DEADTIME_NS 50

struct Register {
  String name;
  std::function<void()> action;
};

struct abc_t {
  float a, b, c;
};

struct ab_t {
  float alpha, beta;
};

struct clarke_input_t {
  float ia;
  float ib;
  float ic;
};

struct clarke_output_t {
  float ialpha;
  float ibeta;
};

struct park_input_t {
  float ialpha;
  float ibeta;
  float theta;
};

struct park_output_t {
  float id;
  float iq;
  float theta;
};

struct svpwm_input_t {
  float valpha;
  float vbeta;
};

struct PIController {
  float kp;
  float ki;
  float integral;
  float out_min;
  float out_max;
};

typedef struct {

} motor_t;

FlexCAN_T4<CAN2, RX_SIZE_256, TX_SIZE_16> Can0;

std::unordered_map<uint8_t, Register> REG_MAP;

uint32_t CAN_ID = 0;
uint8_t MOTOR_CURRENT_LIMIT = 0;
double TARGET_ANGLE = 0;
double MOTOR_ENCODER_ANGLE = 0;
double EXTERNAL_ENCODER_ANGLE = 0;


const uint VDC = 12;

const float Id_ref = 0;
float Iq_ref = 0.1;
const float Ts = 1.0f / PWM_FREQ;
const float MOTOR_1_RESISTANCE = 0.4;
const float MOTOR_1_INDUCTANCE = 1e-3;
PIController pi_d = { kp: MOTOR_1_RESISTANCE, ki: MOTOR_1_RESISTANCE *MOTOR_1_RESISTANCE / MOTOR_1_INDUCTANCE, out_min: -(float)VDC, out_max: (float)VDC };
PIController pi_q = { kp: MOTOR_1_RESISTANCE, ki: MOTOR_1_RESISTANCE *MOTOR_1_RESISTANCE / MOTOR_1_INDUCTANCE, out_min: -(float)VDC, out_max: (float)VDC };
// PIController pi_d = { kp: 0.1f, ki: 500.f, out_min: -(float)VDC, out_max: (float)VDC };
// PIController pi_q = { kp: 0.1f, ki: 500.f, out_min: -(float)VDC, out_max: (float)VDC };


// Motor 1 eFlexPWM submodules (Hardware > PWM1: SM[3]; Hardware > PWM2: SM[0], SM[2])
eFlex::SubModule Sm13(8, 7);   // Motor 1 Phase A (pins 8 and 7)
eFlex::SubModule Sm20(4, 33);  // Motor 1 Phase B (pins 4 and 33)
eFlex::SubModule Sm22(6, 9);   // Motor 1 Phase C (pins 6 and 9)

// Motor 2 eFlexPWM submodules (Hardware > PWM3: SM[1]; Hardware > PWM4: SM[0], SM[1], SM[2])
eFlex::SubModule Sm31(29, 28);  // Motor 2 Phase A (pins 29 and 28)
eFlex::SubModule Sm40(22);      // Motor 2 Phase B Chan A (pins 22)
eFlex::SubModule Sm41(23);      // Motor 2 Phase B Chan B (pins 23)
eFlex::SubModule Sm42(2, 3);    // Motor 2 Phase C (pins 2 and 3)

eFlex::Timer &Tm1 = Sm13.timer();
eFlex::Timer &Tm2 = Sm20.timer();
eFlex::Timer &Tm3 = Sm31.timer();
eFlex::Timer &Tm4 = Sm40.timer();

void setup(void) {
  CAN_ID = EEPROM.read(CAN_ID_ADDRESS);
  MOTOR_CURRENT_LIMIT = EEPROM.read(MOTOR_CURRENT_LIMIT_ADDRESS);
  // pi_d
  // pi_q
  initRegMap();
  Serial.begin(115200);
  delay(400);
  Can0.begin();
  Can0.setBaudRate(1 * 1000 * 1000);
  Can0.setMaxMB(NUM_TX_MAILBOXES + NUM_RX_MAILBOXES);
  for (int i = 0; i < NUM_RX_MAILBOXES; i++) {
    Can0.setMB((FLEXCAN_MAILBOX)i, RX, STD);
  }
  for (int i = NUM_RX_MAILBOXES; i < (NUM_TX_MAILBOXES + NUM_RX_MAILBOXES); i++) {
    Can0.setMB((FLEXCAN_MAILBOX)i, TX, STD);
  }
  Can0.setMBFilter(REJECT_ALL);
  Can0.enableMBInterrupts();
  Can0.onReceive(MB0, canSniff);
  Can0.onReceive(MB1, canSniff);
  Can0.onReceive(MB2, canSniff);
  Can0.onReceive(MB3, canSniff);
  Can0.setMBUserFilter(MB0, 0x01, 0xFF);
  Can0.setMBUserFilter(MB1, 0x03, 0xFF);
  Can0.setMBUserFilter(MB2, 0x0B, 0xFF);
  Can0.setMBUserFilter(MB3, CAN_ID, 0xFF);
  Can0.enhanceFilter(MB3);
  Can0.mailboxStatus();
  Serial.setTimeout(0);
  pwm_init_3phase();
}

void canSniff(const CAN_message_t &msg) {
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
  for (uint8_t i = 0; i < msg.len; i++) {
    Serial.print(msg.buf[i], HEX);
    Serial.print(" ");
  }
  Serial.println();
}

void ReadSerial() {
  String str = Serial.readString(8);
  if (str.equals("help")) {
    regHelp();
    return;
  } else if (str.length() > 0) {
    uint8_t reg = str.toInt();
    if (REG_MAP.count(reg)) {
      REG_MAP[reg].action();
    } else {
      Serial.printf("%d: Not assigned as active register\n", reg);
    }
  }
}

void initRegMap() {
  Register reg;
  reg.name = "Read CAN_ID";
  reg.action = [] {
    Serial.printf("CAN_ID: %d\n", CAN_ID);
  };
  REG_MAP[1] = reg;

  reg.name = "Write CAN_ID";
  reg.action = [] {
    changeCanID();
  };
  REG_MAP[2] = reg;

  reg.name = "Read MOTOR_CONTROLLER_LIMIT";
  reg.action = [] {
    Serial.printf("MOTOR_CURRENT_LIMIT: %d\n", MOTOR_CURRENT_LIMIT);
  };
  REG_MAP[3] = reg;

  reg.name = "Write MOTOR_CURRENT_LIMIT";
  reg.action = [] {
    changeCurrentLimit();
  };
  REG_MAP[4] = reg;

  reg.name = "CAN Mailbox Status";
  reg.action = [] {
    Can0.mailboxStatus();
  };
  REG_MAP[5] = reg;
}

void regHelp() {
  Serial.println("Valid Registers:");
  for (auto pair : REG_MAP) {
    Register reg = pair.second;
    Serial.printf("%d: %s\n", pair.first, reg.name.c_str());
  }
}

void changeCanID() {
  Serial.print("New CAN_ID: ");
  int id;
  while (true) {
    String str = Serial.readString(4);
    if (str.length() > 0) {
      id = str.toInt();
      if (id <= 255 && id >= 0) {
        break;
      }
    }
  }
  CAN_ID = id;
  EEPROM.write(CAN_ID_ADDRESS, CAN_ID);
  Serial.println(CAN_ID);
}

void changeCurrentLimit() {
  Serial.print("New MOTOR_CURRENT_LIMIT (uint8): ");
  int lim;
  while (true) {
    String str = Serial.readString(4);
    if (str.length() > 0) {
      lim = str.toInt();
      if (lim <= 255 && lim >= 0) {
        break;
      }
    }
  }
  MOTOR_CURRENT_LIMIT = lim;
  EEPROM.write(MOTOR_CURRENT_LIMIT_ADDRESS, MOTOR_CURRENT_LIMIT);
  Serial.println(MOTOR_CURRENT_LIMIT);
}

void pwm_init_3phase() {
  pinMode(LED_BUILTIN, OUTPUT);
  digitalWrite(LED_BUILTIN, HIGH);

  Serial.println("Configuring PWM timers");

  eFlex::Config myConfig;

  /* Use full cycle reload */
  myConfig.setReloadLogic(kPWM_ReloadPwmFullCycle);
  /* PWM A & PWM B form a complementary PWM pair */
  myConfig.setPairOperation(kPWM_ComplementaryPwmA);
  myConfig.setPwmFreqHz(PWM_FREQ);

  /* Initialize submodule 0 */
  if (Sm13.configure(myConfig) != true) {
    Serial.println("13 initialization failed");
    exit(EXIT_FAILURE);
  }

  /* Initialize submodule 2, make it use same counter clock as submodule 0. */

  // myConfig.setPrescale (kPWM_Prescale_Divide_1);


  if (Sm20.configure(myConfig) != true) {
    Serial.println("20 initialization failed");
    exit(EXIT_FAILURE);
  }

  myConfig.setClockSource(kPWM_Submodule0Clock);
  myConfig.setInitializationControl(kPWM_Initialize_MasterSync);

  /* Initialize the rest of the submodules the same way as submodule 2 */
  if (Sm22.configure(myConfig) != true) {
    Serial.println("22 initialization failed");
    exit(EXIT_FAILURE);
  }

  myConfig.setClockSource(kPWM_BusClock);

  if (Sm31.configure(myConfig) != true) {
    Serial.println("31 initialization failed");
    exit(EXIT_FAILURE);
  }

  if (Sm40.configure(myConfig) != true) {
    Serial.println("40 initialization failed");
    exit(EXIT_FAILURE);
  }

  myConfig.setClockSource(kPWM_Submodule0Clock);
  // myConfig.setLevel(kPWM_LowTrue);

  if (Sm41.configure(myConfig) != true) {
    Serial.println("41 initialization failed");
    exit(EXIT_FAILURE);
  }

  Sm41.setupLevel(kPWM_LowTrue);

  // myConfig.setLevel(kPWM_HighTrue);

  if (Sm42.configure(myConfig) != true) {
    Serial.println("42 initialization failed");
    exit(EXIT_FAILURE);
  }

  uint16_t deadTimeVal = ((uint64_t)Tm1.srcClockHz() * DEADTIME_NS) / 1000000000;
  Tm1.setupDeadtime(deadTimeVal);

  deadTimeVal = ((uint64_t)Tm2.srcClockHz() * DEADTIME_NS) / 1000000000;
  Tm2.setupDeadtime(deadTimeVal);

  deadTimeVal = ((uint64_t)Tm3.srcClockHz() * DEADTIME_NS) / 1000000000;
  Tm3.setupDeadtime(deadTimeVal);

  deadTimeVal = ((uint64_t)Tm4.srcClockHz() * DEADTIME_NS) / 1000000000;
  Tm4.setupDeadtime(deadTimeVal);

  // synchronize registers and start all submodules
  if (Tm1.begin() != true) {
    Serial.println("Failed to start module 1");
    exit(EXIT_FAILURE);
  } else {

    Serial.println("Submodules successfuly started");
  }

  // synchronize registers and start all submodules
  if (Tm2.begin() != true) {
    Serial.println("Failed to start module 2");
    exit(EXIT_FAILURE);
  } else {

    Serial.println("Submodules successfuly started");
  }

  // synchronize registers and start all submodules
  if (Tm3.begin() != true) {
    Serial.println("Failed to start module 3");
    exit(EXIT_FAILURE);
  } else {

    Serial.println("Submodules successfuly started");
  }

  // synchronize registers and start all submodules
  if (Tm4.begin() != true) {
    Serial.println("Failed to start module 4");
    exit(EXIT_FAILURE);
  } else {

    Serial.println("Submodules successfuly started");
  }

  analogReadResolution(12);
  analogWriteResolution(12);

  // Serial.print("PWM_MOD: "); Serial.println(PWM_MOD);

  // end of PWM setup
  digitalWrite(LED_BUILTIN, LOW);
  pinMode(14, INPUT_PULLDOWN);
}

// For next 4 functions refer to:
// https://ww1.microchip.com/downloads/aemdocuments/documents/fpga/ProductDocuments/UserGuides/sf2_mc_park_invpark_clarke_invclarke_transforms_ug.pdf
// for more information

// Converts three-phase quantities to a two-axis orthogonal stationary reference frame
clarke_output_t clarke(clarke_input_t in) {
  float sqrt3 = 1.73205080757f;
  clarke_output_t out;
  out.ialpha = (2.f * in.ia / 3.f) - ((in.ib - in.ic) / 3.f);
  out.ialpha = (2.f * (in.ib - in.ic) / sqrt3);
  return out;
}

// Converts a two-axis orthogonal stationary reference from to a three-phase stationary reference frame
clarke_input_t inv_clarke(clarke_output_t in) {
  float sqrt3 = 1.73205080757f;
  clarke_input_t out;
  out.ia = in.ialpha;
  out.ib = (-in.ialpha + (sqrt3 * in.ibeta)) / 2.f;
  out.ic = (-in.ialpha - (sqrt3 * in.ibeta)) / 2.f;
  return out;
}

// Converts two-axis orthogonal stationary reference frame quantites to rotated reference frame quantities
park_output_t park(park_input_t in) {
  park_output_t out;
  out.id = (in.ialpha * cosf(in.theta)) + (in.ibeta * sinf(in.theta));
  out.iq = (in.ibeta * cosf(in.theta)) - (in.ialpha * sinf(in.theta));
  return out;
}

// Converts two-axis orthogonal stationary reference frame quantites to rotated reference frame quantities
park_input_t inv_park(park_output_t in) {
  park_input_t out;
  out.ialpha = (in.id * cosf(in.theta)) - (in.iq * sinf(in.theta));
  out.ibeta = (in.iq * cosf(in.theta)) + (in.id * sinf(in.theta));
  return out;
}

clarke_input_t phaseCurrentRead() {
  return { ia: .2, ib: .2, ic: .6 };
}

abc_t svpwm(float Valpha, float Vbeta) {
  // αβ → phase voltages
  float Va = Valpha;
  float Vb = -0.5f * Valpha + 0.8660254f * Vbeta;
  float Vc = -0.5f * Valpha - 0.8660254f * Vbeta;

  // Serial.print("Va: ");
  // Serial.print(Va);
  // Serial.print("\tVb: ");
  // Serial.print(Vb);
  // Serial.print("\tVc: ");
  // Serial.println(Vc);

  // Zero-vector (T0) computation
  float Vmax = fmaxf(Va, fmaxf(Vb, Vc));
  float Vmin = fminf(Va, fminf(Vb, Vc));
  float Voffset = 0.5f * (Vmax + Vmin);

  // Duty cycles (centered)
  float Da = (Va - Voffset) / VDC + 0.5f;
  float Db = (Vb - Voffset) / VDC + 0.5f;
  float Dc = (Vc - Voffset) / VDC + 0.5f;

  // Clamp
  Da = fminf(fmaxf(Da, 0.0f), 1.0f);
  Db = fminf(fmaxf(Db, 0.0f), 1.0f);
  Dc = fminf(fmaxf(Dc, 0.0f), 1.0f);

  // Convert to timer counts
  // Da * PWM_MOD;
  // Db * PWM_MOD;
  // Dc * PWM_MOD;
  Da *= 100;
  Db *= 100;
  Dc *= 100;

// Serial.print("Da: "); Serial.print(Da); Serial.print("\tDb: "); Serial.print(Db); Serial.print("\tDc: "); Serial.println(Dc);

  return abc_t{a: Da, b: Db, c: Dc};
}


abc_t foc() {
  // 1. Read currents
  // adc_read(&Ia, &Ib);
  auto Iabc = phaseCurrentRead();

  // 2. Clarke
  float Ialpha = Iabc.ia;
  float Ibeta = (Iabc.ia + 2.0f * Iabc.ib) * 0.5773503f;
  // Serial.print("Ialpha: ");
  // Serial.print(Ialpha);
  // Serial.print("\tIbeta: ");
  // Serial.println(Ibeta);

  // 3. Park
  float sin_t = sinf(MOTOR_ENCODER_ANGLE);
  float cos_t = cosf(MOTOR_ENCODER_ANGLE);

  float Id = Ialpha * cos_t + Ibeta * sin_t;
  float Iq = -Ialpha * sin_t + Ibeta * cos_t;

  // Serial.print("Id_ref: ");
  // Serial.print(Id_ref);
  // Serial.print("\tIq_ref: ");
  // Serial.println(Iq_ref);

  // 4. Current PI
  float Vd = pi_update(&pi_d, Id_ref - Id, Ts);
  float Vq = pi_update(&pi_q, Iq_ref - Iq, Ts);

  // Serial.print("Vd: ");
  // Serial.print(Vd);
  // Serial.print("\tVq: ");
  // Serial.println(Vq);

  // 5. Inverse Park
  float Valpha = Vd * cos_t - Vq * sin_t;
  float Vbeta = Vd * sin_t + Vq * cos_t;

  // 6. SVPWM
  return svpwm(Valpha, Vbeta);
}

void updateEncoders() {
  MOTOR_ENCODER_ANGLE++;
  uint16_t encoderVal = analogRead(14);
  analogWrite(LED_BUILTIN, encoderVal);
  Iq_ref = (((float)encoderVal / 2048) * MOTOR_CURRENT_LIMIT) - MOTOR_CURRENT_LIMIT;
  Serial.print("Iq_ref: ");
  Serial.println(Iq_ref);
}

static inline float pi_update(PIController *pi, float error, float Ts) {
  float oldIntegral = pi->integral;
  pi->integral += pi->ki * error * Ts;

  // Serial.print("Old Integral: ");
  // Serial.print(oldIntegral);
  // Serial.print("\tNew Integral: ");
  // Serial.println(pi->integral);

  // Clamp integrator
  if (pi->integral > pi->out_max) pi->integral = pi->out_max;
  if (pi->integral < pi->out_min) pi->integral = pi->out_min;

  float output = pi->kp * error + pi->integral;

  // Clamp output
  if (output > pi->out_max) output = pi->out_max;
  if (output < pi->out_min) output = pi->out_min;

  return output;
}

void loop() {
  Can0.events();
  ReadSerial();
  if (Can0.getRXQueueCount() != 0) {
    Serial.println(Can0.getRXQueueCount());
  }
  updateEncoders();
  auto foc1 = foc();

  Sm13.updateDutyCyclePercent(foc1.a, eFlex::ChanA);
  Sm20.updateDutyCyclePercent(foc1.b, eFlex::ChanA);
  Sm22.updateDutyCyclePercent(foc1.c, eFlex::ChanA);

  // delay(500);
  delayMicroseconds ( (1000000U / PWM_FREQ) * 50);
}

// abc_t svpwm(ab_t v_ab, float vbus, uint16_t pwm_max)
// {
//     abc_t duty;

//     // Normalize by DC bus
//     float v_alpha = v_ab.alpha / vbus;
//     float v_beta = v_ab.beta / vbus;

//     // Inverse Clarke: get three-phase voltages
//     float v_a = v_alpha;
//     float v_b = -0.5f * v_alpha + (sqrtf(3) / 2.0f) * v_beta;
//     float v_c = -0.5f * v_alpha - (sqrtf(3) / 2.0f) * v_beta;

//     // Find min/max for re-centering
//     float v_min = fminf(fminf(v_a, v_b), v_c);
//     float v_max = fmaxf(fmaxf(v_a, v_b), v_c);

//     // Normalize to 0..1 range
//     duty.a = (v_a - v_min) / (v_max - v_min);
//     duty.b = (v_b - v_min) / (v_max - v_min);
//     duty.c = (v_c - v_min) / (v_max - v_min);

//     // Scale to PWM wrap
//     duty.a *= pwm_max;
//     duty.b *= pwm_max;
//     duty.c *= pwm_max;

//     return duty;
// }

// void pwm_update_3phase(abc_t duty)
// {
//     // Clamp values
//     uint16_t da = duty.a < 0 ? 0 : (duty.a > pwm_wrap ? pwm_wrap : (uint16_t)duty.a);
//     uint16_t db = duty.b < 0 ? 0 : (duty.b > pwm_wrap ? pwm_wrap : (uint16_t)duty.b);
//     uint16_t dc = duty.c < 0 ? 0 : (duty.c > pwm_wrap ? pwm_wrap : (uint16_t)duty.c);

//     // Apply new levels
//     analogWrite(PHASE_A_PIN, da);
//     analogWrite(PHASE_B_PIN, db);
//     analogWrite(PHASE_C_PIN, dc);
//     // pwm_set_chan_level(slice_a, 0, da);
//     // pwm_set_chan_level(slice_b, 0, db);
//     // pwm_set_chan_level(slice_c, 0, dc);
// }