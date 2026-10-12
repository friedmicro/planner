#include "widgets.h"

#include <algorithm>

#include "models.h"

NumRow::NumRow(const std::string& label, double value, double lo, double hi, double step)
    : Gtk::Box(Gtk::Orientation::HORIZONTAL, 6), adj_(Gtk::Adjustment::create(value, lo, hi, step, 0, 0)) {
  auto lbl = Gtk::make_managed<Gtk::Label>(label);
  lbl->set_halign(Gtk::Align::START);
  lbl->set_width_chars(12);
  spin_ = Gtk::make_managed<Gtk::SpinButton>(adj_);
  spin_->set_digits(2);
  append(*lbl);
  append(*spin_);
}

double NumRow::value() const { return spin_->get_value(); }
void NumRow::set_value(double v) { spin_->set_value(v); }

DaySelector::DaySelector() : Gtk::Box(Gtk::Orientation::HORIZONTAL, 2) {
  for (const auto& d : kDayNames) {
    auto cb = Gtk::make_managed<Gtk::ToggleButton>(d);
    cb->set_has_frame(false);
    toggles_.push_back(cb);
    append(*cb);
  }
}

std::vector<int> DaySelector::selected() const {
  std::vector<int> out;
  for (size_t i = 0; i < toggles_.size(); ++i)
    if (toggles_[i]->get_active()) out.push_back((int)i + 1);
  return out;
}

void DaySelector::set_days(const std::vector<int>& days) {
  for (size_t i = 0; i < toggles_.size(); ++i)
    toggles_[i]->set_active(std::find(days.begin(), days.end(), (int)i + 1) != days.end());
}