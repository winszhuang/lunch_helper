CREATE TABLE IF NOT EXISTS Role (
    id SERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL PRIMARY KEY,
    line_id VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(150) NOT NULL,
    picture TEXT NOT NULL,
    google_maps_api_call_count SMALLINT DEFAULT 5 NOT NULL,
    role_id INTEGER DEFAULT 2 NOT NULL,
    FOREIGN KEY (role_id) REFERENCES Role(id)
);

CREATE INDEX IF NOT EXISTS idx_user_line_id ON "user"(line_id);

CREATE TABLE IF NOT EXISTS Restaurant (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    rating DECIMAL(3,2),
    user_ratings_total INTEGER,
    address VARCHAR(255) NOT NULL,
    google_map_place_id VARCHAR(255) NOT NULL,
    google_map_url TEXT NOT NULL,
    phone_number VARCHAR(25) NOT NULL
);

CREATE TABLE IF NOT EXISTS User_Restaurant (
    user_id INTEGER NOT NULL,
    restaurant_id INTEGER NOT NULL,
    PRIMARY KEY(user_id, restaurant_id),
    FOREIGN KEY (user_id) REFERENCES "user"(id),
    FOREIGN KEY (restaurant_id) REFERENCES Restaurant(id)
);

CREATE INDEX IF NOT EXISTS idx_user_restaurant_user_id ON User_Restaurant(user_id);
CREATE INDEX IF NOT EXISTS idx_user_restaurant_restaurant_id ON User_Restaurant(restaurant_id);

CREATE TABLE IF NOT EXISTS Food (
    id SERIAL PRIMARY KEY,
    name VARCHAR(75) NOT NULL,
    price DECIMAL(5,2) NOT NULL,
    image TEXT,
    description TEXT,
    restaurant_id INTEGER NOT NULL,
    version SMALLINT DEFAULT 0 NOT NULL,
    edit_by INTEGER,
    FOREIGN KEY (restaurant_id) REFERENCES Restaurant(id),
    FOREIGN KEY (edit_by) REFERENCES "user"(id)
);

CREATE INDEX IF NOT EXISTS idx_food_restaurant_id ON Food(restaurant_id);
CREATE INDEX IF NOT EXISTS idx_food_edit_by ON Food(edit_by);

CREATE TABLE IF NOT EXISTS User_Food (
    user_id INTEGER NOT NULL,
    food_id INTEGER NOT NULL,
    PRIMARY KEY(user_id, food_id),
    FOREIGN KEY (user_id) REFERENCES "user"(id),
    FOREIGN KEY (food_id) REFERENCES Food(id)
);

CREATE INDEX IF NOT EXISTS idx_user_food_user_id ON User_Food(user_id);
CREATE INDEX IF NOT EXISTS idx_user_food_food_id ON User_Food(food_id);

CREATE TABLE IF NOT EXISTS Operate_Record (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    food_id INTEGER NOT NULL,
    before TEXT,
    after TEXT,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    operate_category SMALLINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user"(id),
    FOREIGN KEY (food_id) REFERENCES Food(id)
);

CREATE INDEX IF NOT EXISTS idx_operate_record_user_id ON Operate_Record(user_id);
CREATE INDEX IF NOT EXISTS idx_operate_record_update_at ON Operate_Record(update_at);

CREATE TABLE IF NOT EXISTS Feedback (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    food_id INTEGER NOT NULL,
    edit_by INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    "status" VARCHAR(20) DEFAULT 'todo' CHECK ("status" IN ('todo', 'done')) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user"(id),
    FOREIGN KEY (food_id) REFERENCES Food(id),
    FOREIGN KEY (edit_by) REFERENCES "user"(id)
);

CREATE INDEX IF NOT EXISTS idx_feedback_created_at ON Feedback (created_at);
CREATE INDEX IF NOT EXISTS idx_feedback_status ON Feedback ("status");