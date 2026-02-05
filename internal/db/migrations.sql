-- =========================
-- USERS
-- =========================
create table if not exists users (
                                     id bigserial primary key,
                                     first_name varchar(80) not null,
    last_name varchar(80) not null,
    email varchar(200) not null unique,
    phone_number varchar(40),
    password_hash text not null,
    role varchar(20) not null default 'user',
    created_at timestamptz not null default now()
    );

-- =========================
-- MOVIES
-- =========================
create table if not exists movies (
                                      id bigserial primary key,
                                      name varchar(200) not null,
    duration int not null,
    description text,
    poster_url text,
    authors text,
    rating double precision not null default 0,
    created_at timestamptz not null default now()
    );

-- =========================
-- GENRES
-- =========================
create table if not exists genres (
                                      id bigserial primary key,
                                      name varchar(80) not null unique
    );

create table if not exists movies_genres (
                                             movie_id bigint not null references movies(id) on delete cascade,
    genre_id bigint not null references genres(id) on delete cascade,
    primary key(movie_id, genre_id)
    );

-- =========================
-- HALLS
-- =========================
create table if not exists halls (
                                     id bigserial primary key,
                                     name varchar(120) not null,
    location varchar(200),
    total_rows int not null,
    seats_per_row int not null
    );

-- =========================
-- SESSIONS
-- =========================
create table if not exists sessions (
                                        id bigserial primary key,
                                        movie_id bigint not null references movies(id) on delete cascade,
    hall_id bigint not null references halls(id) on delete cascade,
    start_time timestamptz not null,
    end_time timestamptz not null,
    price int not null
    );

-- =========================
-- TICKETS
-- =========================
create table if not exists tickets (
                                       id bigserial primary key,
                                       user_id bigint not null references users(id) on delete cascade,
    session_id bigint not null references sessions(id) on delete cascade,
    row_number int not null,
    seat_number int not null,
    status varchar(30) not null default 'booked',
    created_at timestamptz not null default now(),
    unique(session_id, row_number, seat_number)
    );

-- =========================
-- PAYMENTS
-- =========================
create table if not exists payments (
                                        id bigserial primary key,

                                        ticket_id bigint not null
                                        references tickets(id)
    on delete cascade,

    payment_source varchar(20) not null
    check (payment_source in ('card', 'cash')),

    payment_order varchar(100) not null unique,
    -- example: ORD-2026-02-05-0001

    payment_status varchar(30) not null default 'paid',
    -- paid | failed | refunded

    amount int not null
    check (amount > 0),

    created_at timestamptz not null default now()
    );

-- =========================
-- REVIEWS
-- =========================
create table if not exists reviews (
                                       id bigserial primary key,
                                       user_id bigint not null references users(id) on delete cascade,
    movie_id bigint not null references movies(id) on delete cascade,
    rating int not null check (rating between 1 and 10),
    comment text,
    created_at timestamptz not null default now(),
    unique(user_id, movie_id)
    );

-- =========================
-- INDEXES
-- =========================
create index if not exists idx_sessions_movie on sessions(movie_id);
create index if not exists idx_tickets_user on tickets(user_id);
create index if not exists idx_reviews_movie on reviews(movie_id);
create index if not exists idx_payments_ticket on payments(ticket_id);
create index if not exists idx_payments_created on payments(created_at);