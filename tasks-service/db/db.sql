\c tasks;

CREATE TABLE tbl_book
(
    book_id     SERIAL PRIMARY KEY,
    book_name   VARCHAR(45) NOT NULL UNIQUE,
    book_author VARCHAR(45) NOT NULL
);

CREATE TABLE tbl_student
(
    student_id       SERIAL PRIMARY KEY,
    student_name     VARCHAR(45) NOT NULL UNIQUE,
    student_age      INT NOT NULL,
    street_address   VARCHAR(45),
    street_address_2 VARCHAR(45),
    city             VARCHAR(45),
    state            VARCHAR(45),
    zip_code         VARCHAR(45),
    country          VARCHAR(45)
);

CREATE TABLE tbl_student_book
(
    student_book_id SERIAL PRIMARY KEY,
    student_id      INT NOT NULL REFERENCES tbl_student(student_id) ON UPDATE CASCADE,
    book_id         INT NOT NULL REFERENCES tbl_book(book_id) ON UPDATE CASCADE
);