# ai-planner

A schedule creation tool for converting a list of responsibilities into a simple schedule. While this is not an incredibly complex problem to solve I could not find a utility that would allow me to do so with an API and a client that I like with the exact handling of with respect to organizing my time.

## Why

I have historically had a pretty busy life over the years between hobbies, career, partners, and continuing to keep my engineering\programming skills sharp. This is further compounded by several positions I have worked at over the years which were hybrid\remote with somewhat flexible hours and even my own company at one point. When you combine in having lived in 4 states across multiple timezones throughout my adult life things get very complicated keeping in touch with people. It can also be a really dynamic schedule ranging from your traditional scheduled hours, dates, music events, developer gatherings, and time chatting with friends. While I could keep track of all of this manually...that is really not optimal.

Additionally what I expect from a schedule solver is fairly specific: I have a polyphasic sleep schedule and have for around a decade (2 3 hour blocks), I regularly fluctuate on my sleep schedule depending on need\what's going on, and I have both hard responsibilties (work, partners, etc) as well as soft responsibilities which I am willing to fold on in the name of getting sleep\taking care of my physical health.

The other thing is I have a pretty strong need for an API to automate adding events to my calendar. For example let's say I want an AI agent of mine to add an event to keep in mind or I want to do something like automatically add an event from my email (ex an ics file\calendar) into the planner. This API also makes it insanely trivial for me to do something like run automations on what I am doing, while yes this can be done with just the calendar itself this is entirely local.

The last detail here is I wanted a native Linux app to make entry of my responsibilities faster into the solver. And if we're completely honest there are not that many schedule solver apps with a native client and event fewer with a Linux one.

## Solver rules

1. **Fixed** events (exact day + start/end slot)
2. **Day-constrained required** events (allowed days / time windows)
3. **Floating required** events
4. **Extras** (best-effort, only if `schedule_extra`)
5. **Naps**

- **Hard** constraint = blackout (never violated)
- **Soft** constraint = sleep
- If required work overflows the available non-blackout time, sleep is sacrificed

## Configure and run

Easiest way to run this project is to build the docker image and supply the following variables: `AI_PLANNER_DB_HOST`, `AI_PLANNER_DB_PORT`, `AI_PLANNER_DB_USER`, `AI_PLANNER_DB_PASSWORD`, `AI_PLANNER_DB_NAME`, and `AI_PLANNER_ADDR`

To build the GTK GUI app you need: **gtkmm4**, **libsoup3**, and **json-glib**

```sh
cd gtk
cmake -B build
cmake --build build
./build/ai-planner-gtk
```

Update the code or set the `AI_PLANNER_API` environment variable to the correct domain\IP.
