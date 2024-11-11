-- Insert into countries
INSERT INTO `edugo`.`countries` (`country_id`, `name`) VALUES
(1, 'Thailand'),
(2, 'Japan'),
(3, 'United States');

-- Insert into posts
INSERT INTO `edugo`.`posts` (`posts_id`, `title`, `description`, `image`, `publish_date`, `posts_type`, `country_id`) VALUES
(1, 'New Announcement', 'This is a new announcement for all users', NULL, NOW(), 'Announce', 1),
(2, 'Subject Update', 'An update on the subject matter', NULL, NOW(), 'Announce', 2);

-- Insert into comments
INSERT INTO `edugo`.`comments` (`comments_id`, `comments_text`, `comments_image`, `comments_type`, `publish_date`, `posts_id`) VALUES
(1, 'Great post!', NULL, 'Announce', NOW(), 1),
(2, 'Thanks for the update.', NULL, 'Subject', NOW(), 2);

-- Insert into categories
INSERT INTO `edugo`.`categories` (`category_id`, `name`) VALUES
(1, 'General'),
(2, 'Updates'),
(3, 'Announcements');

-- Insert into announce_posts
INSERT INTO `edugo`.`announce_posts` (`announce_id`, `url`, `attach_file`, `close_date`, `posts_id`, `category_id`) VALUES
(1, 'https://example.com/announcement1', NULL, NULL, 1, 1),
(2, 'https://example.com/announcement2', NULL, NULL, 2, 2);
