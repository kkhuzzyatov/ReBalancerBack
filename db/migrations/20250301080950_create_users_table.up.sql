CREATE TABLE users (
    id SERIAL PRIMARY KEY,
	email VARCHAR(255),
	password VARCHAR(255),
	curAllocation TEXT,
	targetAllocation TEXT,
	taxRate INT
);