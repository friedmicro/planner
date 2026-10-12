#ifndef PLANNER_WIDGETS_H
#define PLANNER_WIDGETS_H

#include <gtkmm.h>

#include <string>
#include <vector>

class NumRow : public Gtk::Box {
 public:
  NumRow(const std::string& label, double value, double lo, double hi, double step);
  double value() const;
  void set_value(double v);

 private:
  Glib::RefPtr<Gtk::Adjustment> adj_;
  Gtk::SpinButton* spin_ = nullptr;
};

class DaySelector : public Gtk::Box {
 public:
  DaySelector();
  std::vector<int> selected() const;
  void set_days(const std::vector<int>& days);

 private:
  std::vector<Gtk::ToggleButton*> toggles_;
};

#endif
