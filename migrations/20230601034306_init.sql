-- +goose Up
create table trades
(
    -- Скорее всего не понадобится, но пусть будет
    id               int generated always as identity primary key,
    tid              text     not null,
    pair_id          smallint not null,
    price            text     not null,
    amount           text     not null,
    side_id          smallint not null,
    ts               bigint   not null,
    is_1m_processed  bool     not null        default false, -- Для выбора тех которые нужно превратить в свечки
    is_15m_processed bool     not null        default false,
    is_1h_processed  bool     not null        default false, -- Для выбора тех которые нужно превратить в свечки
    is_15d_processed bool     not null        default false, -- опечатка is_15d_processed на самом деле 1 день
    created_at       timestamp with time zone default now() not null
);


create unique index trades_tid_uniq_idx on trades (tid);
create index trades_pair_ts_side_idx on trades (pair_id, ts, side_id);

comment
    on column trades.pair_id is
    ' 0 - Unknown
     1 - BTC_USDT
     2 - TRX_USDT
     3 - ETH_USDT
     4 - DOGE_USDT
     5 - BCH_USDT';

comment
    on column trades.side_id is
    ' 0 - Unknown
     1 - BUY
     2 - SELL';

-- Вообще изначально думал что буду по интервалам делать,
-- но тк пар и интервалов немного и их ограниченное ограниченно, одна таблица выглядит норм вариантом
create table candles
(
    pair_id    smallint not null,
    begin_ts   bigint   not null,
    end_ts     bigint   not null,
    time_frame smallint not null,
    data       json     not null        default '{}'::json, -- упущение, что мы не собираемся делать какую-либо выборку
    created_at timestamp with time zone default now() not null,
    PRIMARY KEY (pair_id, time_frame, begin_ts)             -- Чтоб избежать повторений
);
-- create index candles_pair_time_frame_idx on candles (pair_id, time_frame);
-- create index candles_begin_end_idx on candles (begin_ts, end_ts);

comment
    on column candles.time_frame is
    ' 0 - Unknown
     1 - MINUTE_1
     2 - MINUTE_15
     3 - HOUR_1
     4 - DAY_1';

-- +goose Down
