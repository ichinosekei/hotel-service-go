-- Create table for hotels
CREATE TABLE Hotels (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(100) NOT NULL,
                        location VARCHAR(100) NOT NULL
);

-- Create table for rooms
CREATE TABLE Rooms (
                       id SERIAL PRIMARY KEY,
                       hotel_id INTEGER REFERENCES Hotels(id) ON DELETE CASCADE,
                       room_number VARCHAR(10) NOT NULL,
                       price DECIMAL(10, 2) NOT NULL
);
