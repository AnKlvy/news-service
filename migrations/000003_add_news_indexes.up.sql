CREATE INDEX IF NOT EXISTS news_title_idx ON news USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS news_categories_idx ON news USING GIN (categories);
