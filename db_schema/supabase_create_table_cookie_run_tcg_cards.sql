create table cards (id integer, name text, name_eng text, code text, rarity text, rarity_abb text, card_type text, color text, color_sub text, level smallint, plain_text_eng text, plain_text text, expansion text, illustrator text, link text, image_link text);


create policy "Enable read access for all users" on "public"."cards" as permissive
for
select to public using (true);


ALTER TABLE cards ADD PRIMARY KEY (code);