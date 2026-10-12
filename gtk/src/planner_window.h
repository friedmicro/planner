#ifndef PLANNER_WINDOW_H
#define PLANNER_WINDOW_H

#include <gtkmm.h>

#include <string>
#include <vector>

#include "models.h"
#include "widgets.h"

class PlannerWindow : public Gtk::Window {
 public:
  PlannerWindow();

 private:
  static std::string api_path_interval(const char* name);

  Gtk::Box* root_ = nullptr;

  NumRow* cfg_work_ = nullptr;
  NumRow* cfg_weekend_ = nullptr;
  NumRow* cfg_rests_ = nullptr;
  NumRow* cfg_min_block_ = nullptr;
  Gtk::CheckButton* cfg_morning_ = nullptr;
  Gtk::CheckButton* cfg_short_ = nullptr;
  Gtk::CheckButton* cfg_extra_ = nullptr;
  Gtk::CheckButton* cfg_space_ = nullptr;

  Gtk::Box* sleep_form_ = nullptr;
  Gtk::Box* black_form_ = nullptr;
  Gtk::Box* sleep_rows_ = nullptr;
  Gtk::Box* black_rows_ = nullptr;
  std::vector<Interval> sleep_list_;
  std::vector<Interval> black_list_;

  Gtk::Box* event_form_ = nullptr;
  Gtk::Box* event_rows_ = nullptr;
  Gtk::Entry* ev_name_ = nullptr;
  Gtk::SpinButton* ev_time_ = nullptr;
  Gtk::SpinButton* ev_min_ = nullptr;
  Gtk::SpinButton* ev_max_ = nullptr;
  DaySelector* ev_days_ = nullptr;
  NumRow* ev_earliest_ = nullptr;
  NumRow* ev_latest_ = nullptr;
  Gtk::SpinButton* ev_fixed_day_ = nullptr;
  NumRow* ev_fixed_start_ = nullptr;
  NumRow* ev_fixed_end_ = nullptr;
  std::vector<Event> event_list_;

  Gtk::Box* extra_form_ = nullptr;
  Gtk::Box* extra_rows_ = nullptr;
  Gtk::Entry* ex_name_ = nullptr;
  Gtk::SpinButton* ex_time_ = nullptr;
  Gtk::SpinButton* ex_max_ = nullptr;
  NumRow* ex_start_ = nullptr;
  NumRow* ex_end_ = nullptr;
  std::vector<ExtraEvent> extra_list_;

  Gtk::Label* status_ = nullptr;

  static Gtk::SpinButton* make_spin(double val, double lo, double hi, double step);
  Gtk::Box* labeled_row(const std::string& label, Gtk::Widget& w);

  void build_config();
  void build_interval_section(const std::string& title, const std::string& path,
                              std::vector<Interval>&, Gtk::Box*& form, Gtk::Box*& rows);
  void build_event_section();
  void build_extra_section();
  void build_actions();

  void set_status(const std::string& msg, bool error = false);

  void save_config();
  void add_interval(const std::string& path, Interval i);
  void del_interval(const std::string& path, int id);
  void add_event();
  void del_event(int id);
  void add_extra();
  void del_extra(int id);
  void solve();

  static std::string fmt_h(double h);

  void render_intervals(const std::vector<Interval>& list, Gtk::Box* rows,
                        const std::string& path);
  void reload_interval(std::vector<Interval>& list, Gtk::Box* rows);
  void reload_events();
  void reload_extras();
  void reload_all();
};

#endif
