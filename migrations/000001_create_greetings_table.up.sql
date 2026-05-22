CREATE TABLE greetings (
    id          UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    name        TEXT        NOT NULL,
    greeted_at  TIMESTAMPTZ NOT NULL
);
