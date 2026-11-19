DROP TABLE IF EXISTS category_book;
DROP TABLE IF EXISTS reading;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS book;

-- monastery and users reference each other.
SET @previous_foreign_key_checks = @@SESSION.foreign_key_checks;
SET SESSION foreign_key_checks = 0;
DROP TABLE IF EXISTS users, monastery;
SET SESSION foreign_key_checks = @previous_foreign_key_checks;
