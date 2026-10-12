-- ai-planner schema (MariaDB/MySQL).
-- Hours are DECIMAL(5,2) in [0,24). Day of week is 1..7 (Mon..Sun) matching the
-- POC. All tables use InnoDB/utf8mb4.

CREATE TABLE IF NOT EXISTS planner_config (
    id                          INT PRIMARY KEY,
    sleep_per_day_work_week     DECIMAL(5,2) NOT NULL DEFAULT 7.00,
    sleep_per_day_weekend       DECIMAL(5,2) NOT NULL DEFAULT 8.00,
    prefer_morning_nap          TINYINT(1)   NOT NULL DEFAULT 0,
    preferred_rests             INT          NOT NULL DEFAULT 2,
    allow_short_rests           TINYINT(1)   NOT NULL DEFAULT 0,
    default_min_block           DECIMAL(5,2) NOT NULL DEFAULT 1.00,
    schedule_extra              TINYINT(1)   NOT NULL DEFAULT 0,
    prefer_spacing_required_work TINYINT(1)  NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS sleep_times (
    id     INT AUTO_INCREMENT PRIMARY KEY,
    day    TINYINT     NOT NULL,              -- 1..7
    start_h DECIMAL(5,2) NOT NULL,
    end_h   DECIMAL(5,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS blacklisted_times (
    id     INT AUTO_INCREMENT PRIMARY KEY,
    day    TINYINT     NOT NULL,
    start_h DECIMAL(5,2) NOT NULL,
    end_h   DECIMAL(5,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
    id                       INT AUTO_INCREMENT PRIMARY KEY,
    name                     VARCHAR(255) NOT NULL,
    time_h                   DECIMAL(7,2) NOT NULL,
    min_block                DECIMAL(5,2) NOT NULL DEFAULT 0,
    max_block                DECIMAL(5,2) NOT NULL DEFAULT 0,
    days                     JSON         NULL,   -- e.g. [1,3,5]
    hour_earliest            DECIMAL(5,2) NOT NULL DEFAULT 0,
    hour_latest              DECIMAL(5,2) NOT NULL DEFAULT 0,
    hour_start               DECIMAL(5,2) NOT NULL DEFAULT 0,
    hour_end                 DECIMAL(5,2) NOT NULL DEFAULT 0,
    fixed_day                TINYINT      NOT NULL DEFAULT 0,
    fixed_start              DECIMAL(5,2) NOT NULL DEFAULT 0,
    fixed_end                DECIMAL(5,2) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS extra_events (
    id         INT AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    time_h     DECIMAL(7,2) NOT NULL,
    max_time   DECIMAL(7,2) NOT NULL DEFAULT 0,
    start_hour DECIMAL(5,2) NOT NULL DEFAULT 0,
    end_hour   DECIMAL(5,2) NOT NULL DEFAULT 0
);

-- The most recent solved week, keyed by its Monday.
CREATE TABLE IF NOT EXISTS solved_schedules (
    week_start DATE    PRIMARY KEY,           -- Monday of the week
    schedule   JSON    NOT NULL,              -- serialized solver.Solution
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Seed the singleton config row with POC defaults.
INSERT IGNORE INTO planner_config (id) VALUES (1);
