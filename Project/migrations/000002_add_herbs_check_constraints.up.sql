ALTER TABLE herbs ADD CONSTRAINT herbs_price_check CHECK (price >= 0);
ALTER TABLE herbs ADD CONSTRAINT culinary_uses_length_check CHECK (array_length(culinary_uses, 1) BETWEEN 1 AND 5);