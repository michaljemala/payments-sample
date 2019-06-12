CREATE TABLE enum_scheme
(
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL
);
CREATE INDEX idx_enum_scheme_code ON enum_scheme (code);

CREATE TABLE enum_country
(
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL
);
CREATE INDEX idx_enum_country_code ON enum_country (code);

CREATE TABLE enum_currency
(
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL
);
CREATE INDEX idx_enum_currency_code ON enum_currency (code);