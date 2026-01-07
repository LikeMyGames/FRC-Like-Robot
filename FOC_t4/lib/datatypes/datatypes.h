#ifndef DATATYPES_H_
#define DATATYPES_H_

#include <Arduino.h>

typedef struct
{
    // Limits
    float l_current_max;
    float l_current_min;
    float l_in_current_max;
    float l_in_current_min;
    float l_in_current_map_start;
    float l_in_current_map_filter;
    float l_abs_current_max;
    float l_min_erpm;
    float l_max_erpm;
    float l_erpm_start;
    float l_max_erpm_fbrake;
    float l_max_erpm_fbrake_cc;
    float l_min_vin;
    float l_max_vin;
    float l_battery_cut_start;
    float l_battery_cut_end;
    float l_battery_regen_cut_start;
    float l_battery_regen_cut_end;
    bool l_slow_abs_current;
    float l_temp_fet_start;
    float l_temp_fet_end;
    float l_temp_motor_start;
    float l_temp_motor_end;
    float l_temp_accel_dec;
    float l_min_duty;
    float l_max_duty;
    float l_watt_max;
    float l_watt_min;
    float l_current_max_scale;
    float l_current_min_scale;
    float l_duty_start;
    // Overridden limits (Computed during runtime)
    float lo_current_max;
    float lo_current_min;
    float lo_in_current_max;
    float lo_in_current_min;

    // Sensorless (bldc)
    float sl_min_erpm;
    float sl_min_erpm_cycle_int_limit;
    float sl_max_fullbreak_current_dir_change;
    float sl_cycle_int_limit;
    float sl_phase_advance_at_br;
    float sl_cycle_int_rpm_br;
    float sl_bemf_coupling_k;
    // Hall sensor
    int8_t hall_table[8];
    float hall_sl_erpm;

    // FOC
    float foc_current_kp;
    float foc_current_ki;
    float foc_f_zv;
    float foc_dt_us;
    float foc_encoder_offset;
    bool foc_encoder_inverted;
    float foc_encoder_ratio;
    float foc_motor_l;
    float foc_motor_ld_lq_diff;
    float foc_motor_r;
    float foc_motor_flux_linkage;
    float foc_observer_gain;
    float foc_observer_gain_slow;
    float foc_observer_offset;
    float foc_pll_kp;
    float foc_pll_ki;
    float foc_duty_dowmramp_kp;
    float foc_duty_dowmramp_ki;
    float foc_start_curr_dec;
    float foc_start_curr_dec_rpm;
    float foc_openloop_rpm;
    float foc_openloop_rpm_low;
    float foc_sl_openloop_hyst;
    float foc_sl_openloop_time;
    float foc_sl_openloop_time_lock;
    float foc_sl_openloop_time_ramp;
    float foc_sl_openloop_boost_q;
    float foc_sl_openloop_max_q;
    uint8_t foc_hall_table[8];
    float foc_hall_interp_erpm;
    float foc_sl_erpm_start;
    float foc_sl_erpm;
    float foc_sat_comp;
    bool foc_temp_comp;
    float foc_temp_comp_base_temp;
    float foc_current_filter_const;
    float foc_hfi_amb_current;
    uint8_t foc_hfi_amb_tres;
    float foc_hfi_voltage_start;
    float foc_hfi_voltage_run;
    float foc_hfi_voltage_max;
    float foc_hfi_gain;
    float foc_hfi_max_err;
    float foc_hfi_hyst;
    float foc_sl_erpm_hfi;
    float foc_hfi_reset_erpm;
    uint16_t foc_hfi_start_samples;
    float foc_hfi_obs_ovr_sec;
    uint8_t foc_offsets_cal_mode;
    float foc_offsets_current[3];
    float foc_offsets_voltage[3];
    float foc_offsets_voltage_undriven[3];
    bool foc_phase_filter_enable;
    bool foc_phase_filter_disable_fault;
    float foc_phase_filter_max_erpm;
    // Field Weakening
    float foc_fw_current_max;
    float foc_fw_duty_start;
    float foc_fw_ramp_time;
    float foc_fw_q_current_factor;
    // FOC_SPEED_SRC foc_speed_soure;
    bool foc_short_ls_on_zero_duty;
    float foc_overmod_factor;

    // PID_RATE sp_pid_loop_rate;

    // Speed PID
    float s_pid_kp;
    float s_pid_ki;
    float s_pid_kd;
    float s_pid_kd_filter;
    float s_pid_min_erpm;
    bool s_pid_allow_braking;
    float s_pid_ramp_erpms_s;
    // S_PID_SPEED_SRC s_pid_speed_source;

    // Pos PID
    float p_pid_kp;
    float p_pid_ki;
    float p_pid_kd;
    float p_pid_kd_proc;
    float p_pid_kd_filter;
    float p_pid_ang_div;
    float p_pid_gain_dec_angle;
    float p_pid_offset;

    // Current controller
    float cc_startup_boost_duty;
    float cc_min_current;
    float cc_gain;
    float cc_ramp_step_max;

    // Misc
    int32_t m_fault_stop_time_ms;
    float m_duty_ramp_step;
    float m_current_backoff_gain;
    uint32_t m_encoder_counts;
    float m_encoder_sin_offset;
    float m_encoder_sin_amp;
    float m_encoder_cos_offset;
    float m_encoder_cos_amp;
    float m_encoder_sincos_filter_constant;
    float m_encoder_sincos_phase_correction;
    bool m_invert_direction;
    int m_drv8301_oc_adj;
    float m_bldc_f_sw_min;
    float m_bldc_f_sw_max;
    float m_dc_f_sw;
    float m_ntc_motor_beta;
    float m_ptc_motor_coeff;
    int m_hall_extra_samples;
    int m_batt_filter_const;
    float m_ntcx_ptcx_temp_base;
    float m_ntcx_ptcx_res;
    // Setup info
    uint8_t si_motor_poles;
    float si_gear_ratio;
    float si_wheel_diameter;
    int si_battery_cells;
    float si_battery_ah;
    float si_motor_nl_current;
} mc_configuration;

// typedef enum
// {
//     MC_STATE_OFF = 0,
//     MC_STATE_DETECTING,
//     MC_STATE_RUNNING,
//     MC_STATE_FULL_BRAKE,
// } mc_state;

typedef enum
{
    FOC_SPEED_SRC_CORRECTED = 0,
    FOC_SPEED_SRC_OBSERVER,
} FOC_SPEED_SRC;

typedef enum
{
    CONTROL_MODE_DUTY = 0,
    CONTROL_MODE_SPEED,
    CONTROL_MODE_CURRENT,
    CONTROL_MODE_CURRENT_BRAKE,
    CONTROL_MODE_POS,
    CONTROL_MODE_HANDBRAKE,
    CONTROL_MODE_OPENLOOP,
    CONTROL_MODE_OPENLOOP_PHASE,
    CONTROL_MODE_OPENLOOP_DUTY,
    CONTROL_MODE_OPENLOOP_DUTY_PHASE,
    CONTROL_MODE_NONE
} mc_control_mode;

#endif