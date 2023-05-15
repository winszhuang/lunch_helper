CREATE TABLE Role (
    id SERIAL PRIMARY KEY,
    name VARCHAR(30),
    description TEXT
);

CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    line_id VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(150) NOT NULL,
    picture TEXT NOT NULL,
    google_maps_api_call_count SMALLINT DEFAULT 5,
    role_id INTEGER DEFAULT 2,
    FOREIGN KEY (role_id) REFERENCES Role(id)
);

CREATE INDEX idx_user_line_id ON "user"(line_id);

CREATE TABLE Restaurant (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150),
    rating DECIMAL(3,2),
    user_ratings_total INTEGER,
    address VARCHAR(255),
    google_map_place_id VARCHAR(255),
    google_map_url TEXT,
    phone_number VARCHAR(25)
);

CREATE TABLE User_Restaurant (
    user_id INTEGER,
    restaurant_id INTEGER,
    PRIMARY KEY(user_id, restaurant_id),
    FOREIGN KEY (user_id) REFERENCES "user"(id),
    FOREIGN KEY (restaurant_id) REFERENCES Restaurant(id)
);

CREATE INDEX idx_user_restaurant_user_id ON User_Restaurant(user_id);
CREATE INDEX idx_user_restaurant_restaurant_id ON User_Restaurant(restaurant_id);

CREATE TABLE Food (
    id SERIAL PRIMARY KEY,
    name VARCHAR(75),
    price DECIMAL(5,2),
    image TEXT,
    description TEXT,
    restaurant_id INTEGER,
    version SMALLINT DEFAULT 0,
    edit_by INTEGER,
    FOREIGN KEY (restaurant_id) REFERENCES Restaurant(id),
    FOREIGN KEY (edit_by) REFERENCES "user"(id)
);

CREATE INDEX idx_food_restaurant_id ON Food(restaurant_id);
CREATE INDEX idx_food_edit_by ON Food(edit_by);

CREATE TABLE User_Food (
    user_id INTEGER,
    food_id INTEGER,
    PRIMARY KEY(user_id, food_id),
    FOREIGN KEY (user_id) REFERENCES "user"(id),
    FOREIGN KEY (food_id) REFERENCES Food(id)
);

CREATE INDEX idx_user_food_user_id ON User_Food(user_id);
CREATE INDEX idx_user_food_food_id ON User_Food(food_id);

CREATE TABLE Operate_Record (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    food_id INTEGER,
    before TEXT,
    after TEXT,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    operate_category SMALLINT,
    FOREIGN KEY (user_id) REFERENCES "user"(id),
    FOREIGN KEY (food_id) REFERENCES Food(id)
);

CREATE INDEX idx_operate_record_user_id ON Operate_Record(user_id);
CREATE INDEX idx_operate_record_update_at ON Operate_Record(update_at);

CREATE TABLE Feedback (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    food_id INTEGER,
    edit_by INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    "status" VARCHAR(20) DEFAULT 'todo' CHECK ("status" IN ('todo', 'done')),
    FOREIGN KEY (user_id) REFERENCES "user"(id),
    FOREIGN KEY (food_id) REFERENCES Food(id),
    FOREIGN KEY (edit_by) REFERENCES "user"(id)
);

CREATE INDEX idx_feedback_created_at ON Feedback (created_at);
CREATE INDEX idx_feedback_status ON Feedback ("status");