DROP TABLE Feedback;
DROP INDEX IF EXISTS idx_feedback_status;
DROP INDEX IF EXISTS idx_feedback_created_at;

DROP TABLE Operate_Record;
DROP INDEX IF EXISTS idx_operate_record_user_id;
DROP INDEX IF EXISTS idx_operate_record_update_at;

DROP TABLE User_Food;
DROP INDEX IF EXISTS idx_user_food_user_id;
DROP INDEX IF EXISTS idx_user_food_food_id;

DROP TABLE Food;
DROP INDEX IF EXISTS idx_food_restaurant_id;
DROP INDEX IF EXISTS idx_food_edit_by;

DROP TABLE User_Restaurant;
DROP INDEX IF EXISTS idx_user_restaurant_user_id;
DROP INDEX IF EXISTS idx_user_restaurant_restaurant_id;

DROP TABLE Restaurant;

DROP INDEX IF EXISTS idx_user_line_id;
DROP TABLE "User";
DROP TABLE Role;
