CREATE TABLE users (
    user_id UUID NOT NULL PRIMARY KEY,
    first_name CHARACTER VARYING(45) NOT NULL,
    last_name CHARACTER VARYING(45) NOT NULL,
    login VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    phone_number VARCHAR(9) NOT NULL,
    balance NUMERIC DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE books (
    book_id UUID NOT NULL PRIMARY KEY,
    title VARCHAR NOT NULL,
    author VARCHAR(150),
    price NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE orders (
    order_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES users(user_id),
    book_id UUID NOT NULL REFERENCES books(book_id),
    payed NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- INSERT INTO users(user_id, first_name, last_name, phone_number, updated_at) VALUES
-- ('8b85ff2a-8091-11ed-a1eb-0242ac120002','Samandar', 'Foziljonov', '997191323', now()); 

-- INSERT INTO books(book_id, title, author, price, updated_at) VALUES
-- ('f23fa31a-8091-11ed-a1eb-0242ac120002', 'Tokaki', 'Murakami', 1500, now());

-- INSERT INTO orders(order_id, user_id, book_id, payed) VALUES
-- ('5af01b9c-8092-11ed-a1eb-0242ac120002', '8b85ff2a-8091-11ed-a1eb-0242ac120002', 'f23fa31a-8091-11ed-a1eb-0242ac120002', 1500);

-- SELECT 
--     users.first_name || ' ' || users.last_name as fullname,
--     books.title,
--     orders.payed,
--     orders.created_at,
--     orders.updated_at
-- FROM
--     orders
-- JOIN users ON orders.user_id = users.user_id
-- JOIN books ON orders.book_id = books.book_id;