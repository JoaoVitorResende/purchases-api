-- migrations/001_init.sql
create table purchases (
    id          bigserial primary key,
    total_cents bigint not null check (total_cents > 0),
    created_at  timestamptz not null
);

create table purchase_items (
    id               bigserial primary key,
    purchase_id      bigint not null references purchases(id) on delete cascade,
    item_id          text not null,
    title            text not null,
    unit_price_cents bigint not null check (unit_price_cents > 0),
    quantity         int not null check (quantity > 0)
);

create index idx_purchase_items_purchase_id on purchase_items(purchase_id);

-- Cópia local do catálogo (alimentada pelos eventos do StoreApi)
create table catalog_items (
    item_id     text primary key,
    title       text not null,
    price_cents bigint not null check (price_cents > 0)
);