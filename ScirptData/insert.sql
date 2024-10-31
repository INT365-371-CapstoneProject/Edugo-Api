-- Insert into accounts table
INSERT INTO `edugo`.`accounts` (`account_id`, `phone_number`, `create_on`, `update_on`, `last_login`, `username`, `password`, `email`, `role`)
VALUES 
(1, '0912345678', NOW(), NOW(), NOW(), 'admin_user', 'password123', 'admin@example.com', 'admin'),
(2, '0923456789', NOW(), NOW(), NOW(), 'normal_user', 'password123', 'user@example.com', 'user'),
(3, '0934567890', NOW(), NOW(), NOW(), 'provider_user', 'password123', 'provider@example.com', 'provider');

-- Insert into admins table
INSERT INTO `edugo`.`admins` (`admin_id`, `firstname`, `lastname`, `status`, `account_id`)
VALUES 
(1, 'Admin', 'One', 'Active', 1);

-- Insert into users table
INSERT INTO `edugo`.`users` (`user_id`, `firstname`, `lastname`, `account_id`)
VALUES 
(1, 'User', 'One', 2);

-- Insert into providers table
INSERT INTO `edugo`.`providers` (`provider_id`, `provider_name`, `url`, `address`, `status`, `verify`, `account_id`)
VALUES 
(1, 'Provider One', 'https://providerone.com', '123 Provider St.', 'Active', 'Y', 3);

-- Insert into posts table
INSERT INTO `edugo`.`posts` (`posts_id`, `title`, `description`, `url`, `posts_type`, `publish_date`, `close_date`, `provider_id`, `user_id`)
VALUES 
(1, 'Announcement Title', 'This is an announcement description', NULL, 'Announce', NOW(), NULL, NULL, 1),
(2, 'Subject Title', 'This is a subject description', 'https://example.com', 'Subject', NOW(), NOW(), 1, NULL);

-- Insert into comments table
INSERT INTO `edugo`.`comments` (`comments_id`, `comments_text`, `comments_type`, `publishDate`, `posts_id`, `user_id`, `provider_id`)
VALUES 
(1, 'This is a comment on announcement', 'Announce', NOW(), 1, 1, NULL),
(2, 'This is a comment on subject', 'Subject', NOW(), 2, NULL, 1);

-- Insert into categories table
INSERT INTO `edugo`.`categories` (`category_id`, `category_name`, `category_type`, `posts_id`)
VALUES 
(1, 'General Announcement', 'Announce', 1),
(2, 'Math Subject', 'Subject', 2);