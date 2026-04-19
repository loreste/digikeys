ALTER TABLE citizens DROP COLUMN IF EXISTS origin_country;
ALTER TABLE cards DROP COLUMN IF EXISTS origin_country;
ALTER TABLE embassies DROP COLUMN IF EXISTS origin_country;
DELETE FROM countries WHERE code = 'CD';
