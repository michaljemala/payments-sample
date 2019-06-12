CREATE TABLE IF NOT EXISTS payment
(
    id                             UUID PRIMARY KEY,

    amount_value                   NUMERIC NOT NULL,
    amount_currency                TEXT    NOT NULL REFERENCES enum_currency (code),

    scheme_type                    TEXT    NOT NULL REFERENCES enum_scheme (code),

    creditor_name                  TEXT    NOT NULL,
    creditor_address_line1         TEXT    NOT NULL,
    creditor_address_line2         TEXT,
    creditor_address_city          TEXT    NOT NULL,
    creditor_address_region        TEXT,
    creditor_address_postal_code   TEXT    NOT NULL,
    creditor_address_country_code  TEXT    NOT NULL REFERENCES enum_country (code),
    creditor_account_name          TEXT    NOT NULL,
    creditor_account_number        TEXT    NOT NULL,
    creditor_account_provider_code TEXT    NOT NULL,
    creditor_account_provider_name TEXT,

    debtor_name                    TEXT    NOT NULL,
    debtor_address_line1           TEXT    NOT NULL,
    debtor_address_line2           TEXT,
    debtor_address_city            TEXT    NOT NULL,
    debtor_address_region          TEXT,
    debtor_address_postal_code     TEXT    NOT NULL,
    debtor_address_country_code    TEXT    NOT NULL REFERENCES enum_country (code),
    debtor_account_name            TEXT    NOT NULL,
    debtor_account_number          TEXT    NOT NULL,
    debtor_account_provider_code   TEXT    NOT NULL,
    debtor_account_provider_name   TEXT
);
CREATE INDEX idx_payment_amount_value ON payment (amount_value);
CREATE INDEX idx_payment_amount_currency ON payment (amount_currency);
CREATE INDEX idx_payment_creditor_name ON payment (creditor_name);
CREATE INDEX idx_payment_debtor_name ON payment (debtor_name);