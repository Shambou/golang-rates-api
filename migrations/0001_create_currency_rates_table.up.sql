CREATE TABLE IF NOT EXISTS currency_rates
(
    id             serial constraint currency_rates_pk primary key,
    base_currency  char(3)        not null,
    quote_currency char(3)        not null,
    rate           decimal(12, 6) not null,
    date           DATE
);
CREATE INDEX IF NOT EXISTS "base_currency_quote_currency_index" ON "public"."currency_rates" USING BTREE ("base_currency","quote_currency");
