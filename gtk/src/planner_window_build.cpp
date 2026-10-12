#include "planner_window.h"

void PlannerWindow::build_config() {
  auto frame = Gtk::make_managed<Gtk::Frame>("Config");
  auto box = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::VERTICAL, 6);
  box->set_margin(8);

  auto grid = Gtk::make_managed<Gtk::Grid>();
  grid->set_column_spacing(12);
  grid->set_row_spacing(6);

  auto add = [&](int r, int c, Gtk::Widget& w) { grid->attach(w, c, r, 1, 1); };

  cfg_work_ = Gtk::make_managed<NumRow>("Sleep (workday)", 7, 0, 24, 0.5);
  cfg_weekend_ = Gtk::make_managed<NumRow>("Sleep (weekend)", 8, 0, 24, 0.5);
  cfg_rests_ = Gtk::make_managed<NumRow>("Rests", 2, 1, 12, 1);
  cfg_min_block_ = Gtk::make_managed<NumRow>("Min block (h)", 1, 0, 24, 0.5);
  add(0, 0, *cfg_work_); add(0, 1, *cfg_weekend_);
  add(1, 0, *cfg_rests_); add(1, 1, *cfg_min_block_);

  cfg_morning_ = Gtk::make_managed<Gtk::CheckButton>("Prefer morning nap");
  cfg_short_ = Gtk::make_managed<Gtk::CheckButton>("Allow short rests");
  cfg_extra_ = Gtk::make_managed<Gtk::CheckButton>("Schedule extras");
  cfg_space_ = Gtk::make_managed<Gtk::CheckButton>("Space required work");
  auto checks = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::HORIZONTAL, 12);
  checks->append(*cfg_morning_);
  checks->append(*cfg_short_);
  checks->append(*cfg_extra_);
  checks->append(*cfg_space_);

  auto save = Gtk::make_managed<Gtk::Button>("Save config");
  save->signal_clicked().connect([this] { save_config(); });

  box->append(*grid);
  box->append(*checks);
  box->append(*save);
  frame->set_child(*box);
  root_->append(*frame);
}

void PlannerWindow::build_interval_section(const std::string& title, const std::string& path,
                                           std::vector<Interval>&,
                                           Gtk::Box*& form, Gtk::Box*& rows) {
  auto frame = Gtk::make_managed<Gtk::Frame>(title);
  auto vbox = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::VERTICAL, 6);
  vbox->set_margin(8);

  form = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::HORIZONTAL, 6);
  auto day_spin = make_spin(1, 1, 7, 1);
  auto start_spin = make_spin(22, 0, 24, 0.5);
  auto end_spin = make_spin(8, 0, 24, 0.5);
  auto add_btn = Gtk::make_managed<Gtk::Button>("Add");
  add_btn->signal_clicked().connect(
      [this, path, day_spin, start_spin, end_spin] {
        Interval i;
        i.day = (int)day_spin->get_value();
        i.start = start_spin->get_value();
        i.end = end_spin->get_value();
        add_interval(path, i);
      });

  form->append(*Gtk::make_managed<Gtk::Label>("Day"));
  form->append(*day_spin);
  form->append(*Gtk::make_managed<Gtk::Label>("Start"));
  form->append(*start_spin);
  form->append(*Gtk::make_managed<Gtk::Label>("End"));
  form->append(*end_spin);
  form->append(*add_btn);

  rows = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::VERTICAL, 2);

  vbox->append(*form);
  vbox->append(*rows);
  frame->set_child(*vbox);
  root_->append(*frame);
}

void PlannerWindow::build_event_section() {
  auto frame = Gtk::make_managed<Gtk::Frame>("Commitments (required events)");
  auto vbox = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::VERTICAL, 6);
  vbox->set_margin(8);

  event_form_ = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::VERTICAL, 4);
  ev_name_ = Gtk::make_managed<Gtk::Entry>();
  ev_name_->set_placeholder_text("Name");
  ev_time_ = make_spin(1, 0, 168, 0.5);
  ev_min_ = make_spin(0, 0, 24, 0.5);
  ev_max_ = make_spin(0, 0, 24, 0.5);
  ev_days_ = Gtk::make_managed<DaySelector>();
  ev_earliest_ = Gtk::make_managed<NumRow>("Earliest hr", 0, 0, 24, 0.5);
  ev_latest_ = Gtk::make_managed<NumRow>("Latest hr", 24, 0, 24, 0.5);
  ev_fixed_day_ = make_spin(0, 0, 7, 1);
  ev_fixed_start_ = Gtk::make_managed<NumRow>("Fixed start", 0, 0, 24, 0.5);
  ev_fixed_end_ = Gtk::make_managed<NumRow>("Fixed end", 0, 0, 24, 0.5);

  auto row1 = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::HORIZONTAL, 6);
  row1->append(*ev_name_);
  row1->append(*Gtk::make_managed<Gtk::Label>("hrs"));
  row1->append(*ev_time_);
  row1->append(*Gtk::make_managed<Gtk::Label>("min"));
  row1->append(*ev_min_);
  row1->append(*Gtk::make_managed<Gtk::Label>("max"));
  row1->append(*ev_max_);

  auto row2 = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::HORIZONTAL, 6);
  row2->append(*Gtk::make_managed<Gtk::Label>("Days"));
  row2->append(*ev_days_);

  auto row3 = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::HORIZONTAL, 6);
  row3->append(*ev_earliest_);
  row3->append(*ev_latest_);

  auto row4 = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::HORIZONTAL, 6);
  row4->append(*Gtk::make_managed<Gtk::Label>("Fixed day (0=no)"));
  row4->append(*ev_fixed_day_);
  row4->append(*ev_fixed_start_);
  row4->append(*ev_fixed_end_);

  auto add_btn = Gtk::make_managed<Gtk::Button>("Add commitment");
  add_btn->signal_clicked().connect([this] { add_event(); });

  event_form_->append(*row1);
  event_form_->append(*row2);
  event_form_->append(*row3);
  event_form_->append(*row4);
  event_form_->append(*add_btn);

  event_rows_ = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::VERTICAL, 2);

  vbox->append(*event_form_);
  vbox->append(*event_rows_);
  frame->set_child(*vbox);
  root_->append(*frame);
}

void PlannerWindow::build_extra_section() {
  auto frame = Gtk::make_managed<Gtk::Frame>("Optional activities");
  auto vbox = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::VERTICAL, 6);
  vbox->set_margin(8);

  extra_form_ = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::HORIZONTAL, 6);
  ex_name_ = Gtk::make_managed<Gtk::Entry>();
  ex_name_->set_placeholder_text("Name");
  ex_time_ = make_spin(1, 0, 168, 0.5);
  ex_max_ = make_spin(0, 0, 168, 0.5);
  ex_start_ = Gtk::make_managed<NumRow>("Start", 0, 0, 24, 0.5);
  ex_end_ = Gtk::make_managed<NumRow>("End", 24, 0, 24, 0.5);
  auto add_btn = Gtk::make_managed<Gtk::Button>("Add");
  add_btn->signal_clicked().connect([this] { add_extra(); });

  extra_form_->append(*ex_name_);
  extra_form_->append(*Gtk::make_managed<Gtk::Label>("hrs"));
  extra_form_->append(*ex_time_);
  extra_form_->append(*Gtk::make_managed<Gtk::Label>("max"));
  extra_form_->append(*ex_max_);
  extra_form_->append(*ex_start_);
  extra_form_->append(*ex_end_);
  extra_form_->append(*add_btn);

  extra_rows_ = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::VERTICAL, 2);

  vbox->append(*extra_form_);
  vbox->append(*extra_rows_);
  frame->set_child(*vbox);
  root_->append(*frame);
}

void PlannerWindow::build_actions() {
  auto box = Gtk::make_managed<Gtk::Box>(Gtk::Orientation::HORIZONTAL, 12);
  auto solve_btn = Gtk::make_managed<Gtk::Button>("Solve week");
  solve_btn->signal_clicked().connect([this] { solve(); });
  auto refresh = Gtk::make_managed<Gtk::Button>("Reload");
  refresh->signal_clicked().connect([this] { reload_all(); });
  status_ = Gtk::make_managed<Gtk::Label>("");
  status_->set_halign(Gtk::Align::START);
  box->append(*solve_btn);
  box->append(*refresh);
  box->append(*status_);
  root_->append(*box);
}
