-- Insert into countries
INSERT INTO `edugo`.`countries` (`country_id`, `name`) VALUES
(1, 'Thailand'),
(2, 'Japan'),
(3, 'United States'),
(4, 'United Kingdom'),
(5, 'Australia'),
(6, 'Canada');

-- Insert into posts
INSERT INTO `edugo`.`posts` (`posts_id`, `title`, `description`, `image`, `publish_date`, `posts_type`, `country_id`) VALUES
(1, 'New Announcement', 'This is a new announcement for all users', NULL, NOW(), 'Announce', 1),
(2, 'Subject Update', 'An update on the subject matter', NULL, "2024-11-10 12:29:29", 'Announce', 2);
(3, 'Scholarship Opportunities', 'Various scholarships available for students', NULL, NOW(), 'Subject', 3),
(4, 'Study Abroad Programs', 'Information on study abroad programs', NULL, NOW(), 'Subject', 4);


-- Insert into comments
INSERT INTO `edugo`.`comments` (`comments_id`, `comments_text`, `comments_image`, `comments_type`, `publish_date`, `posts_id`) VALUES
(1, 'Great post!', NULL, 'Announce', NOW(), 1);

-- Insert into categories
INSERT INTO `edugo`.`categories` (`category_id`, `name`) VALUES
(1, 'Full Scholarships'),
(2, 'Partial Tuition Scholarships'),
(3, 'Merit-Based Scholarships'),
(4, 'Need-Based Scholarships'),
(5, 'Research and Special Project Scholarships'),
(6, 'Government and Corporate Scholarships');

-- Insert into announce_posts
INSERT INTO `edugo`.`announce_posts` (`announce_id`, `url`, `attach_file`, `close_date`, `posts_id`, `category_id`) VALUES
(1, 'https://example.com/announcement1', NULL, "2024-12-20 12:29:29", 1, 1),
(2, 'https://example.com/announcement2', NULL, "2024-11-15 12:29:29", 2, 2);
