#ifndef PLANNER_MODELS_H
#define PLANNER_MODELS_H

#include <string>
#include <vector>

struct Config {
  int id = 1;
  double sleep_per_day_work_week = 7;
  double sleep_per_day_weekend = 8;
  bool prefer_morning_nap = false;
  int preferred_rests = 2;
  bool allow_short_rests = false;
  double default_min_block = 1;
  bool schedule_extra = false;
  bool prefer_spacing_required_work = false;
};

struct Interval {
  int id = 0;
  int day = 1;
  double start = 0;
  double end = 0;
};

struct Event {
  int id = 0;
  std::string name;
  double time = 0;
  double min_block = 0;
  double max_block = 0;
  std::vector<int> days;
  double hour_earliest = 0;
  double hour_latest = 0;
  double hour_start = 0;
  double hour_end = 0;
  int fixed_day = 0;
  double fixed_start = 0;
  double fixed_end = 0;
};

struct ExtraEvent {
  int id = 0;
  std::string name;
  double time = 0;
  double max_time = 0;
  double start_hour = 0;
  double end_hour = 0;
};

inline const std::vector<std::string> kDayNames = {"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"};

inline std::string interval_json(const Interval& i) {
  return "{\"day\":" + std::to_string(i.day) +
         ",\"start\":" + std::to_string(i.start) +
         ",\"end\":" + std::to_string(i.end) + "}";
}

inline std::string config_json(const Config& c) {
  return "{\"sleep_per_day_work_week\":" + std::to_string(c.sleep_per_day_work_week) +
         ",\"sleep_per_day_weekend\":" + std::to_string(c.sleep_per_day_weekend) +
         ",\"prefer_morning_nap\":" + (c.prefer_morning_nap ? "true" : "false") +
         ",\"preferred_rests\":" + std::to_string(c.preferred_rests) +
         ",\"allow_short_rests\":" + (c.allow_short_rests ? "true" : "false") +
         ",\"default_min_block\":" + std::to_string(c.default_min_block) +
         ",\"schedule_extra\":" + (c.schedule_extra ? "true" : "false") +
         ",\"prefer_spacing_required_work\":" + (c.prefer_spacing_required_work ? "true" : "false") + "}";
}

inline std::string event_json(const Event& e) {
  std::string days = "[";
  for (size_t i = 0; i < e.days.size(); ++i) {
    if (i) days += ",";
    days += std::to_string(e.days[i]);
  }
  days += "]";
  return "{\"name\":\"" + e.name +
         "\",\"time\":" + std::to_string(e.time) +
         ",\"min_block\":" + std::to_string(e.min_block) +
         ",\"max_block\":" + std::to_string(e.max_block) +
         ",\"days\":" + days +
         ",\"hour_earliest\":" + std::to_string(e.hour_earliest) +
         ",\"hour_latest\":" + std::to_string(e.hour_latest) +
         ",\"hour_start\":" + std::to_string(e.hour_start) +
         ",\"hour_end\":" + std::to_string(e.hour_end) +
         ",\"fixed_day\":" + std::to_string(e.fixed_day) +
         ",\"fixed_start\":" + std::to_string(e.fixed_start) +
         ",\"fixed_end\":" + std::to_string(e.fixed_end) + "}";
}

inline std::string extra_json(const ExtraEvent& x) {
  return "{\"name\":\"" + x.name +
         "\",\"time\":" + std::to_string(x.time) +
         ",\"max_time\":" + std::to_string(x.max_time) +
         ",\"start_hour\":" + std::to_string(x.start_hour) +
         ",\"end_hour\":" + std::to_string(x.end_hour) + "}";
}

#endif
