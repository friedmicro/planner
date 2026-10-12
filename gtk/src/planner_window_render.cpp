#include "planner_window.h"

#include "http.h"

void PlannerWindow::render_intervals(const std::vector<Interval>& list, Gtk::Box* rows,
                                     const std::string& path) {
  while (auto child = rows->get_first_child()) rows->remove(*child);
  for (const auto& i : list) {
    auto box = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::HORIZONTAL, 8);
    auto lbl = Gtk::make_managed<Gtk::Label>(
        kDayNames[i.day - 1] + "  " + fmt_h(i.start) + "–" + fmt_h(i.end));
    lbl->set_halign(Gtk::Align::START);
    lbl->set_hexpand(true);
    auto del = Gtk::make_managed<Gtk::Button>("✕");
    int id = i.id;
    del->signal_clicked().connect([this, path, id] { del_interval(path, id); });
    box->append(*lbl);
    box->append(*del);
    rows->append(*box);
  }
}

void PlannerWindow::reload_interval(std::vector<Interval>& list, Gtk::Box* rows) {
  const bool is_sleep = (&list == &sleep_list_);
  const std::string path = is_sleep ? "/api/sleep" : "/api/blacklist";
  try {
    auto node = parse_json(http_json("GET", path));
    auto arr = node ? json_node_get_array(node) : nullptr;
    list.clear();
    if (arr) {
      guint len = json_array_get_length(arr);
      for (guint i = 0; i < len; ++i) {
        auto obj = json_array_get_object_element(arr, i);
        Interval it;
        it.id = (int)json_object_get_int_member(obj, "id");
        it.day = (int)json_object_get_int_member(obj, "day");
        it.start = json_object_get_double_member(obj, "start");
        it.end = json_object_get_double_member(obj, "end");
        list.push_back(it);
      }
    }
    render_intervals(list, rows, path);
  } catch (const std::exception& e) {
    set_status(e.what(), true);
  }
}

void PlannerWindow::reload_events() {
  event_list_.clear();
  try {
    auto node = parse_json(http_json("GET", "/api/events"));
    auto arr = node ? json_node_get_array(node) : nullptr;
    if (arr) {
      guint len = json_array_get_length(arr);
      for (guint i = 0; i < len; ++i) {
        auto obj = json_array_get_object_element(arr, i);
        Event e;
        e.id = (int)json_object_get_int_member(obj, "id");
        e.name = jstr(json_object_get_member(obj, "name"));
        e.time = json_object_get_double_member(obj, "time");
        e.fixed_day = (int)json_object_get_int_member(obj, "fixed_day");
        e.days = jday_array(json_object_get_member(obj, "days"));
        event_list_.push_back(e);
      }
    }
    while (auto child = event_rows_->get_first_child()) event_rows_->remove(*child);
    for (const auto& e : event_list_) {
      auto box = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::HORIZONTAL, 8);
      auto lbl = Gtk::make_managed<Gtk::Label>(e.name + "  (" + std::to_string(e.time) + "h)");
      lbl->set_halign(Gtk::Align::START);
      lbl->set_hexpand(true);
      auto del = Gtk::make_managed<Gtk::Button>("✕");
      int id = e.id;
      del->signal_clicked().connect([this, id] { del_event(id); });
      box->append(*lbl);
      box->append(*del);
      event_rows_->append(*box);
    }
  } catch (const std::exception& e) {
    set_status(e.what(), true);
  }
}

void PlannerWindow::reload_extras() {
  extra_list_.clear();
  try {
    auto node = parse_json(http_json("GET", "/api/extras"));
    auto arr = node ? json_node_get_array(node) : nullptr;
    if (arr) {
      guint len = json_array_get_length(arr);
      for (guint i = 0; i < len; ++i) {
        auto obj = json_array_get_object_element(arr, i);
        ExtraEvent x;
        x.id = (int)json_object_get_int_member(obj, "id");
        x.name = jstr(json_object_get_member(obj, "name"));
        x.time = json_object_get_double_member(obj, "time");
        extra_list_.push_back(x);
      }
    }
    while (auto child = extra_rows_->get_first_child()) extra_rows_->remove(*child);
    for (const auto& x : extra_list_) {
      auto box = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::HORIZONTAL, 8);
      auto lbl = Gtk::make_managed<Gtk::Label>(x.name + "  (" + std::to_string(x.time) + "h)");
      lbl->set_halign(Gtk::Align::START);
      lbl->set_hexpand(true);
      auto del = Gtk::make_managed<Gtk::Button>("✕");
      int id = x.id;
      del->signal_clicked().connect([this, id] { del_extra(id); });
      box->append(*lbl);
      box->append(*del);
      extra_rows_->append(*box);
    }
  } catch (const std::exception& e) {
    set_status(e.what(), true);
  }
}

void PlannerWindow::reload_all() {
  reload_interval(sleep_list_, sleep_rows_);
  reload_interval(black_list_, black_rows_);
  reload_events();
  reload_extras();
  try {
    auto node = parse_json(http_json("GET", "/api/config"));
    auto obj = as_object(node);
    if (obj) {
      cfg_work_->set_value(json_object_get_double_member(obj, "sleep_per_day_work_week"));
      cfg_weekend_->set_value(json_object_get_double_member(obj, "sleep_per_day_weekend"));
      cfg_rests_->set_value(json_object_get_int_member(obj, "preferred_rests"));
      cfg_min_block_->set_value(json_object_get_double_member(obj, "default_min_block"));
      cfg_morning_->set_active(json_object_get_boolean_member(obj, "prefer_morning_nap"));
      cfg_short_->set_active(json_object_get_boolean_member(obj, "allow_short_rests"));
      cfg_extra_->set_active(json_object_get_boolean_member(obj, "schedule_extra"));
      cfg_space_->set_active(json_object_get_boolean_member(obj, "prefer_spacing_required_work"));
    }
  } catch (const std::exception& e) {
    set_status(e.what(), true);
  }
}
