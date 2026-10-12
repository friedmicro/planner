#include "planner_window.h"

#include <cmath>
#include <cstdio>

#include "http.h"

PlannerWindow::PlannerWindow() : Gtk::Window() {
  set_title("planner");
  set_default_size(720, 640);

  auto scroller = Gtk::make_managed<Gtk::ScrolledWindow>();
  scroller->set_policy(Gtk::PolicyType::NEVER, Gtk::PolicyType::AUTOMATIC);
  root_ = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::VERTICAL, 8);
  root_->set_margin(12);

  build_config();
  build_interval_section("Sleep", api_path_interval("sleep"), sleep_list_, sleep_form_, sleep_rows_);
  build_interval_section("Blackout", api_path_interval("blacklist"), black_list_, black_form_, black_rows_);
  build_event_section();
  build_extra_section();
  build_actions();

  scroller->set_child(*root_);
  set_child(*scroller);

  reload_all();
}

std::string PlannerWindow::api_path_interval(const char* name) {
  return std::string("/api/") + name;
}

Gtk::SpinButton* PlannerWindow::make_spin(double val, double lo, double hi, double step) {
  auto adj = Gtk::Adjustment::create(val, lo, hi, step, 0, 0);
  auto s = Gtk::make_managed<Gtk::SpinButton>(adj);
  s->set_digits(2);
  return s;
}

Gtk::Box* PlannerWindow::labeled_row(const std::string& label, Gtk::Widget& w) {
  auto box = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::HORIZONTAL, 6);
  auto lbl = Gtk::make_managed<Gtk::Label>(label);
  lbl->set_width_chars(12);
  box->append(*lbl);
  box->append(w);
  return box;
}

void PlannerWindow::set_status(const std::string& msg, bool error) {
  status_->set_markup(std::string("<span color='") + (error ? "red" : "gray") + "'>" + msg + "</span>");
}

void PlannerWindow::save_config() {
  Config c;
  c.sleep_per_day_work_week = cfg_work_->value();
  c.sleep_per_day_weekend = cfg_weekend_->value();
  c.preferred_rests = (int)cfg_rests_->value();
  c.default_min_block = cfg_min_block_->value();
  c.prefer_morning_nap = cfg_morning_->get_active();
  c.allow_short_rests = cfg_short_->get_active();
  c.schedule_extra = cfg_extra_->get_active();
  c.prefer_spacing_required_work = cfg_space_->get_active();
  try {
    http_json("PUT", "/api/config", config_json(c));
    set_status("Config saved");
  } catch (const std::exception& e) {
    set_status(e.what(), true);
  }
}

void PlannerWindow::add_interval(const std::string& path, Interval i) {
  try {
    auto resp = parse_json(http_json("POST", path, interval_json(i)));
    auto obj = as_object(resp);
    auto idnode = obj ? json_object_get_member(obj, "id") : nullptr;
    i.id = jint(idnode);
    reload_interval(path == "/api/sleep" ? sleep_list_ : black_list_,
                    path == "/api/sleep" ? sleep_rows_ : black_rows_);
    set_status("Added");
  } catch (const std::exception& e) {
    set_status(e.what(), true);
  }
}

void PlannerWindow::del_interval(const std::string& path, int id) {
  try {
    http_json("DELETE", path + "/" + std::to_string(id));
    reload_interval(path == "/api/sleep" ? sleep_list_ : black_list_,
                    path == "/api/sleep" ? sleep_rows_ : black_rows_);
  } catch (const std::exception& e) {
    set_status(e.what(), true);
  }
}

void PlannerWindow::add_event() {
  Event e;
  e.name = ev_name_->get_text();
  e.time = ev_time_->get_value();
  e.min_block = ev_min_->get_value();
  e.max_block = ev_max_->get_value();
  e.days = ev_days_->selected();
  e.hour_earliest = ev_earliest_->value();
  e.hour_latest = ev_latest_->value();
  e.fixed_day = (int)ev_fixed_day_->get_value();
  e.fixed_start = ev_fixed_start_->value();
  e.fixed_end = ev_fixed_end_->value();
  try {
    http_json("POST", "/api/events", event_json(e));
    reload_events();
    set_status("Commitment added");
  } catch (const std::exception& ex) {
    set_status(ex.what(), true);
  }
}

void PlannerWindow::del_event(int id) {
  try {
    http_json("DELETE", "/api/events/" + std::to_string(id));
    reload_events();
  } catch (const std::exception& e) {
    set_status(e.what(), true);
  }
}

void PlannerWindow::add_extra() {
  ExtraEvent x;
  x.name = ex_name_->get_text();
  x.time = ex_time_->get_value();
  x.max_time = ex_max_->get_value();
  x.start_hour = ex_start_->value();
  x.end_hour = ex_end_->value();
  try {
    http_json("POST", "/api/extras", extra_json(x));
    reload_extras();
    set_status("Activity added");
  } catch (const std::exception& e) {
    set_status(e.what(), true);
  }
}

void PlannerWindow::del_extra(int id) {
  try {
    http_json("DELETE", "/api/extras/" + std::to_string(id));
    reload_extras();
  } catch (const std::exception& e) {
    set_status(e.what(), true);
  }
}

void PlannerWindow::solve() {
  try {
    http_json("POST", "/api/solve");
    set_status("Week solved (see /api/schedule)");
  } catch (const std::exception& e) {
    set_status(e.what(), true);
  }
}

std::string PlannerWindow::fmt_h(double h) {
  int hh = (int)h;
  int mm = (int)std::round((h - hh) * 60);
  if (mm == 60) { hh += 1; mm = 0; }
  char buf[16];
  std::snprintf(buf, sizeof(buf), "%02d:%02d", hh % 24, mm);
  return buf;
}
