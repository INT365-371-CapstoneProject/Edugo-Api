USE edugo;

-- ตาราง accounts
INSERT INTO `edugo`.`accounts` (`account_id`, `phone_number`, `create_on`, `update_on`, `last_login`, `username`, `password`, `email`, `role`)
VALUES 
(1, '0812345678', NOW(), NOW(), NOW(), 'admin01', 'pass1234', 'admin01@example.com', 'admin'),
(2, '0823456789', NOW(), NOW(), NOW(), 'user01', 'pass1234', 'user01@example.com', 'user'),
(3, '0834567890', NOW(), NOW(), NOW(), 'provider01', 'pass1234', 'provider01@example.com', 'provider');

-- ตาราง admins
INSERT INTO `edugo`.`admins` (`admin_id`, `firstname`, `lastname`, `status`, `account_id`)
VALUES 
(1, 'John', 'Doe', 'Active', 1),
(2, 'Jane', 'Smith', 'Inactive', 2);

-- ตาราง users
INSERT INTO `edugo`.`users` (`user_id`, `firstname`, `lastname`, `account_id`)
VALUES 
(1, 'Alice', 'Brown', 2),
(2, 'Bob', 'White', 3);

-- ตาราง providers
INSERT INTO `edugo`.`providers` (`provider_id`, `provider_name`, `url`, `address`, `status`, `verify`, `account_id`)
VALUES 
(1, 'Provider A', 'http://provider-a.com', '123 Main St', 'Active', 'Y', 3),
(2, 'Provider B', 'http://provider-b.com', '456 Second St', 'Inactive', 'N', 3);

-- ตาราง posts
INSERT INTO `edugo`.`posts` (`posts_id`, `title`, `description`, `url`, `image`, `attach_file`, `posts_type`, `publish_date`, `close_date`, `provider_id`, `user_id`)
VALUES 
(1, 'Post Title 1', 'Description of post 1', 'http://post1.com', 'image_data_1', 'file_data_1', 'Announce', NOW(), NULL, 1, NULL),
(2, 'Post Title 2', 'Description of post 2', NULL, NULL, NULL, 'Subject', NOW(), '2025-12-31 23:59:59', NULL, 1);

-- ตาราง comments
INSERT INTO `edugo`.`comments` (`comments_id`, `comments_text`, `comments_image`, `comments_type`, `publish_date`, `user_id`, `provider_id`, `posts_id`)
VALUES 
(1, 'This is a comment', NULL, 'Announce', NOW(), 1, NULL, 1),
(2, 'Another comment', 'comment_image_data', 'Subject', NOW(), NULL, 2, 1);

-- ตาราง tags
INSERT INTO `edugo`.`tags` (`tags_id`, `tags_name`)
VALUES 
(1, 'Education'),
(2, 'Technology'),
(3, 'Science');

-- ตาราง post_tags
INSERT INTO `edugo`.`post_tags` (`posts_id`, `tags_id`)
VALUES 
(1, 1),
(1, 2),
(2, 3);
