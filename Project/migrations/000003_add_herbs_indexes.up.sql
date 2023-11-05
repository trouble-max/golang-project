CREATE INDEX IF NOT EXISTS herbs_name_idx ON herbs USING gin (to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS herbs_culinary_uses_idx ON herbs USING gin (culinary_uses);