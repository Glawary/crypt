CREATE TABLE public.cryptoexchange
(
    cryptoexchange_id               SERIAL PRIMARY KEY,
    cryptoexchange_name             TEXT UNIQUE,
    cryptocurrency_create_timestamp TIMESTAMP DEFAULT NOW()
);

CREATE TABLE public.cryptocurrency
(
    cryptocurrency_id               SERIAL PRIMARY KEY,
    cryptocurrency_ticker           TEXT UNIQUE,
    cryptocurrency_create_timestamp TIMESTAMP DEFAULT NOW()
);

CREATE TABLE public.cryptocurrency_data
(
    cryptocurrency_id     INT REFERENCES cryptocurrency (cryptocurrency_id),
    cryptoexchange_id     INT REFERENCES cryptoexchange (cryptoexchange_id),
    data_olhcv            json,
    data_order_book       json,
    data_trades           json,
    data_create_timestamp TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (cryptocurrency_id, cryptoexchange_id)
);