cmake_minimum_required(VERSION 3.16)
project(ai_planner_gtk VERSION 0.1.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# pkg-config is required to locate gtkmm4 / libsoup3 / json-glib.
find_package(PkgConfig REQUIRED)
pkg_check_modules(GTKMM REQUIRED gtkmm-4.0)
pkg_check_modules(JSON_GLIB REQUIRED json-glib-1.0)
pkg_check_modules(LIBSOUP REQUIRED libsoup-3.0)

add_executable(ai-planner-gtk
  src/main.cpp
  src/http.cpp
  src/widgets.cpp
  src/planner_window.cpp
  src/planner_window_build.cpp
  src/planner_window_render.cpp
)

target_include_directories(ai-planner-gtk PRIVATE
  ${GTKMM_INCLUDE_DIRS}
  ${JSON_GLIB_INCLUDE_DIRS}
  ${LIBSOUP_INCLUDE_DIRS}
)

target_compile_options(ai-planner-gtk PRIVATE ${GTKMM_CFLAGS_OTHER} ${JSON_GLIB_CFLAGS_OTHER} ${LIBSOUP_CFLAGS_OTHER})

target_link_libraries(ai-planner-gtk PRIVATE
  ${GTKMM_LIBRARIES}
  ${JSON_GLIB_LIBRARIES}
  ${LIBSOUP_LIBRARIES}
)
