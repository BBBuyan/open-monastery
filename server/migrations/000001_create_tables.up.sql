CREATE TABLE IF NOT EXISTS monastery (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS book (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price INT UNSIGNED NOT NULL
) CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    user_type VARCHAR(50) NOT NULL,
    monastery_id BIGINT,

    CONSTRAINT fk_users_monastery 
        FOREIGN KEY (monastery_id)
        REFERENCES monastery(id)
);

CREATE TABLE IF NOT EXISTS reading (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    book_id BIGINT NOT NULL,
    monk_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,

    CONSTRAINT fk_reading_monk FOREIGN KEY (monk_id)
        REFERENCES users(id),

    CONSTRAINT fk_reading_book FOREIGN KEY (book_id)
        REFERENCES book(id),

    CONSTRAINT fk_reading_customer FOREIGN KEY (customer_id)
        REFERENCES users(id)
);

