USE edugo;

-- Insert into accounts
INSERT INTO accounts (account_id, phone_number, create_on, update_on, last_login, username, password, email, role)
VALUES
(1, '0812345678', NOW(), NOW(), NOW(), 'adminuser', 'pass123', 'admin@example.com', 'admin'),
(2, '0823456789', NOW(), NOW(), NOW(), 'normaluser', 'pass456', 'user@example.com', 'user'),
(3, '0834567890', NOW(), NOW(), NOW(), 'provideruser', 'pass789', 'provider@example.com', 'provider');

-- Insert into admins
INSERT INTO admins (admin_id, firstname, lastname, status, account_id)
VALUES
(1, 'Admin', 'One', 'Active', 1);

-- Insert into users
INSERT INTO users (user_id, firstname, lastname, account_id)
VALUES
(1, 'User', 'One', 2),
(2, 'User', 'Two', 2);

-- Insert into providers
INSERT INTO providers (provider_id, provider_name, url, address, status, verify, account_id)
VALUES
(1, 'Provider One', 'http://provider1.com', '123 Provider St', 'Active', 'Y', 3),
(2, 'Provider Two', NULL, '456 Provider Rd', 'Inactive', 'N', 3);

-- Insert into posts
INSERT INTO posts (posts_id, title, description, url, posts_type, publish_date, close_date, provider_id, user_id)
VALUES
(1, 'First Announcement', 'This is the first announcement', NULL, 'Announce', NOW(), NULL, 1, NULL),
(2, 'Subject Post', 'This is a subject-related post', NULL, 'Subject', NOW(), NULL, NULL, 1);

-- Insert into comments
INSERT INTO comments (comments_id, comments_text, comments_type, publish_date, user_id, provider_id, posts_id)
VALUES
(1, 'Great announcement!', 'Announce', NOW(), 1, NULL, 1),
(2, 'Looking forward to the event', 'Subject', NOW(), NULL, 1, 2);

-- Insert into tags
INSERT INTO tags (tags_id, tags_name)
VALUES
(1, 'Education'),
(2, 'Announcement');

-- Insert into post_tags
INSERT INTO post_tags (posts_id, tags_id)
VALUES
(1, 2),
(2, 1);