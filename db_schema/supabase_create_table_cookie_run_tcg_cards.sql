-- create table cards (id integer, name text, name_eng text, code text, rarity text, rarity_abb text, card_type text, color text, color_sub text, level smallint, plain_text_eng text, plain_text text, expansion text, illustrator text, link text, image_link text);--done in prod and test
 -- create policy "Enable read access for all users" on "public"."cards" as permissive
-- for
-- select to public using (true);--done in prod and test
 -- ALTER TABLE cards ADD PRIMARY KEY (code); --done in prod and test

ALTER TABLE cards RENAME COLUMN name TO name_kr; --done in test plan to do in prod


ALTER TABLE cards RENAME COLUMN level TO card_level; --done in test plan to do in prod


ALTER TABLE cards
DROP CONSTRAINT cards_pkey; --done in test


ALTER TABLE cards ADD PRIMARY KEY (id); --done in test


ALTER TABLE cards ADD CONSTRAINT code_unique UNIQUE (code); --done in test


CREATE EXTENSION IF NOT EXISTS pg_trgm; --done in test


CREATE EXTENSION IF NOT EXISTS pgroonga; --done in test


ALTER TABLE cards ADD COLUMN name_eng_lower TEXT GENERATED ALWAYS AS (LOWER(name_eng)) STORED; --done in test


CREATE INDEX trgm_idx_cards_name_eng ON cards USING gin (name_eng_lower gin_trgm_ops); --done in test

-- DROP INDEX IF EXISTS trgm_idx_cards_name_eng;

CREATE INDEX pgroonga_name_kr_index ON cards USING pgroonga (name_kr); --done in test

-- DROP INDEX IF EXISTS pgroonga_name_kr_index;
 -- CREATE INDEX pgroonga_name_eng_lower_index ON cards USING pgroonga (name_eng_lower);
 -- DROP INDEX IF EXISTS pgroonga_name_eng_lower_index;
 -- CREATE INDEX pgroonga_name_eng_index ON cards USING pgroonga (name_eng);
 -- DROP INDEX IF EXISTS pgroonga_name_eng_index;
 -- CREATE INDEX card_code_index ON cards (code);
 -- DROP INDEX IF EXISTS card_code_index;