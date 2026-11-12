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
        REFERENCES category(id)
);
