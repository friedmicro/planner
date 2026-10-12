#include <gtkmm.h>

#include "planner_window.h"

int main(int argc, char** argv) {
  auto app = Gtk::Application::create("lab.friedmicro.aiplanner");
  return app->make_window_and_run<PlannerWindow>(argc, argv);
}