CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    username      TEXT NOT NULL,
    email         TEXT NOT NULL,
    hash_password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS contollers 
(
    id_contorller SERIAL PRIMARY KEY,
    type_controller INTEGER,
    number_controller INTEGER,
);

CREATE TABLE IF NOT EXISTS messanges 
(
    id_messange SERIAL PRIMARY KEY,
    status_controller INTEGER,
    charge_controller INTEGER,
    temperature_MK_controller INTEGER
);

CREATE TABLE IF NOT EXISTS messanges_contollers_leack
(
    id_messange references messanges(id_messange),
    leack INTEGER
);

CREATE TABLE IF NOT EXISTS messanges_contollers_module
(
    id_messange references messanges(id_messange),
    temperature INTEGER,
    humidity INTEGER,
    pressure INTEGER,
    gas INTEGER
);

CREATE TABLE IF NOT EXISTS messanges_contollers_enviroment
(
    id_messange references messanges(id_messange),
    temperature INTEGER,
    humidity INTEGER,
    pressure INTEGER,
    VOC INTEGER,
    gas1 INTEGER,
    gas2 INTEGER,
    gas3 INTEGER,
    pm1 INTEGER,
    pm25 INTEGER,
    pm10 INTEGER,
    fire INTEGER,
    smoke INTEGER
);
