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
    email VARCHAR(254) NULL,
    password_hash VARCHAR(255) NULL,
    phone_number VARCHAR(12) NULL,

    UNIQUE KEY uq_users_email (email),
    CONSTRAINT fk_users_monastery 
        FOREIGN KEY (monastery_id)
        REFERENCES monastery(id)
);

CREATE TABLE IF NOT EXISTS reading (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    book_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    monk_id BIGINT,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    created_at DATETIME(6) NULL,
    assigned_at DATETIME(6) NULL,

    CONSTRAINT fk_reading_monk FOREIGN KEY (monk_id)
        REFERENCES users(id),

    CONSTRAINT fk_reading_book FOREIGN KEY (book_id)
        REFERENCES book(id),

    CONSTRAINT fk_reading_customer FOREIGN KEY (customer_id)
        REFERENCES users(id)
);


CREATE TABLE IF NOT EXISTS category (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS category_book (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    book_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,

    CONSTRAINT fk_category_book_book FOREIGN KEY (book_id)
        REFERENCES book(id),
    CONSTRAINT fk_category_book_category FOREIGN KEY (category_id)
        REFERENCES category(id),

    INDEX idx_category_book_category_book (category_id, book_id)
);
