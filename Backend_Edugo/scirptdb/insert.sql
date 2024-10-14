use edugo_v3;
-- Insert into accounts
INSERT INTO `edugo_v3`.`accounts` 
(`account_id`, `phone_number`, `create_on`, `update_on`, `last_login`, `username`, `password`, `email`, `role`)
VALUES 
(1, '0812345678', NOW(), NOW(), NOW(), 'admin_user', 'password123', 'admin@example.com', 'admin'),
(2, '0823456789', NOW(), NOW(), NOW(), 'user_one', 'password456', 'user1@example.com', 'user'),
(3, '0834567890', NOW(), NOW(), NOW(), 'provider_one', 'password789', 'provider1@example.com', 'provider');

-- Insert into admins
INSERT INTO `edugo_v3`.`admins` 
(`admin_id`, `firstname`, `lastname`, `status`, `account_id`)
VALUES 
(1, 'John', 'Doe', 'Active', 1);

-- Insert into categories
INSERT INTO `edugo_v3`.`categories` 
(`category_id`, `name`, `description`)
VALUES 
(1, 'Technology', 'All things related to technology and innovation.'),
(2, 'Education', 'Content focused on educational topics.');

-- Insert into tags
INSERT INTO `edugo_v3`.`tags` 
(`tag_id`, `tag_name`)
VALUES 
(1, 'Python'), 
(2, 'JavaScript');

-- Insert into users
INSERT INTO `edugo_v3`.`users` 
(`user_id`, `firstname`, `lastname`, `account_id`)
VALUES 
(1, 'Jane', 'Smith', 2),
(2, 'Bob', 'Johnson', 3);

-- Insert into subjects
INSERT INTO `edugo_v3`.`subjects` 
(`subject_id`, `title`, `description`, `attach_file`, `tag_id`, `user_id`)
VALUES 
(1, 'Introduction to Python', 'A beginner guide to Python programming.', NULL, 1, 1),
(2, 'JavaScript Best Practices', 'Best practices for writing clean JavaScript.', NULL, 2, 2);

-- Insert into providers
INSERT INTO `edugo_v3`.`providers` 
(`provider_id`, `provider_name`, `url`, `address`, `status`, `verify`, `account_id`)
VALUES 
(1, 'Provider X', 'www.providerx.com', '123 Main St.', 'Active', 'Y', 3);

-- Insert into posts
INSERT INTO `edugo_v3`.`posts` 
(`post_id`, `title`, `description`, `url`, `attach_file`, `category_id`, `provider_id`)
VALUES 
(1, 'Latest Tech Trends', 'Discussion about the latest trends in technology.', NULL, NULL, 1, 1);

-- Insert into postments
INSERT INTO `edugo_v3`.`postments` 
(`postment_id`, `comment_text`, `comment_file`, `post_id`, `user_id`, `provider_id`)
VALUES 
(1, 'Great post about tech!', NULL, 1, 1, 1);

-- Insert into subjectments
INSERT INTO `edugo_v3`.`subjectments` 
(`subjectment_id`, `comment_text`, `comment_file`, `commu_id`, `user_id`, `provider_id`)
VALUES 
(1, 'Interesting subject about Python!', NULL, 1, 1, 1),
(2, 'Very useful JavaScript tips.', NULL, 2, 2, 1);