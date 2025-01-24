-- Insert into accounts
INSERT INTO `edugo`.`accounts` (`account_id`, `username`, `password`, `email`, `create_on`, `last_login`, `update_on`, `role`) 
VALUES 
(1, 'admin_user', '$2a$14$.GKpMx.V.JlLsDdYYXmay.ZJKODGZK06MoDW7ELp07rIjYRWf1/xC', 'admin@example.com', NOW(), NOW(), NOW(), 'admin'),
(2, 'provider_user', '$2a$14$.GKpMx.V.JlLsDdYYXmay.ZJKODGZK06MoDW7ELp07rIjYRWf1/xC', 'provider@example.com', NOW(), NOW(), NOW(), 'provider'),
(3, 'normal_user', '$2a$14$.GKpMx.V.JlLsDdYYXmay.ZJKODGZK06MoDW7ELp07rIjYRWf1/xC', 'user@example.com', NOW(), NOW(), NOW(), 'user');

-- Insert into admins
INSERT INTO `edugo`.`admins` (`admin_id`, `firstname`, `lastname`, `status`, `account_id`) 
VALUES 
(1, 'John', 'Doe', 'Active', 1);

-- Insert into providers
INSERT INTO `edugo`.`providers` (`provider_id`, `company_name`, `url`, `address`, `status`, `phone`, `verify`, `account_id`) 
VALUES 
(1, 'Tech Solutions', 'www.techsolutions.com', '123 Main St', 'Active', '0123456789', 'Y', 2);


-- Insert into posts
INSERT INTO `edugo`.`posts` (`posts_id`, `description`, `image`, `publish_date`, `posts_type`, `account_id`) 
VALUES 
(1, 'Welcome to our new platform!', NULL, NOW(), 'Announce', 2),
(2, 'New subject added to the curriculum.', NULL, NOW(), 'Announce', 2),
(3, 'New feature added to the platform.', NULL, NOW(), 'Subject', 3);

-- Insert into comments
INSERT INTO `edugo`.`comments` (`comments_id`, `comments_text`, `comments_image`, `publish_date`, `posts_id`, `account_id`) 
VALUES 
(1, 'Great announcement!', NULL, NOW(), 1, 3),
(2, 'Excited for this update.', NULL, NOW(), 2, 3);

-- Insert into categories
INSERT INTO `edugo`.`categories` (`category_id`, `name`) VALUES
(1, 'Full Scholarships'),
(2, 'Partial Tuition Scholarships'),
(3, 'Merit-Based Scholarships'),
(4, 'Need-Based Scholarships'),
(5, 'Research and Special Project Scholarships'),
(6, 'Government and Corporate Scholarships');

-- Insert into countries
INSERT INTO `edugo`.`countries` (`country_id`, `name`) VALUES
(1, 'Afghanistan'),
(2, 'Albania'),
(3, 'Algeria'),
(4, 'Andorra'),
(5, 'Angola'),
(6, 'Antigua and Barbuda'),
(7, 'Argentina'),
(8, 'Armenia'),
(9, 'Australia'),
(10, 'Austria'),
(11, 'Azerbaijan'),
(12, 'Bahamas'),
(13, 'Bahrain'),
(14, 'Bangladesh'),
(15, 'Barbados'),
(16, 'Belarus'),
(17, 'Belgium'),
(18, 'Belize'),
(19, 'Benin'),
(20, 'Bhutan'),
(21, 'Bolivia'),
(22, 'Bosnia and Herzegovina'),
(23, 'Botswana'),
(24, 'Brazil'),
(25, 'Brunei'),
(26, 'Bulgaria'),
(27, 'Burkina Faso'),
(28, 'Burundi'),
(29, 'Cambodia'),
(30, 'Cameroon'),
(31, 'Canada'),
(32, 'Cape Verde'),
(33, 'Central African Republic'),
(34, 'Chad'),
(35, 'Chile'),
(36, 'China'),
(37, 'Colombia'),
(38, 'Comoros'),
(39, 'Congo (Democratic Republic)'),
(40, 'Congo (Republic)'),
(41, 'Costa Rica'),
(42, 'Croatia'),
(43, 'Cuba'),
(44, 'Cyprus'),
(45, 'Czech Republic'),
(46, 'Denmark'),
(47, 'Djibouti'),
(48, 'Dominica'),
(49, 'Dominican Republic'),
(50, 'Ecuador'),
(51, 'Egypt'),
(52, 'El Salvador'),
(53, 'Equatorial Guinea'),
(54, 'Eritrea'),
(55, 'Estonia'),
(56, 'Eswatini'),
(57, 'Ethiopia'),
(58, 'Fiji'),
(59, 'Finland'),
(60, 'France'),
(61, 'Gabon'),
(62, 'Gambia'),
(63, 'Georgia'),
(64, 'Germany'),
(65, 'Ghana'),
(66, 'Greece'),
(67, 'Grenada'),
(68, 'Guatemala'),
(69, 'Guinea'),
(70, 'Guinea-Bissau'),
(71, 'Guyana'),
(72, 'Haiti'),
(73, 'Honduras'),
(74, 'Hungary'),
(75, 'Iceland'),
(76, 'India'),
(77, 'Indonesia'),
(78, 'Iran'),
(79, 'Iraq'),
(80, 'Ireland'),
(81, 'Israel'),
(82, 'Italy'),
(83, 'Jamaica'),
(84, 'Japan'),
(85, 'Jordan'),
(86, 'Kazakhstan'),
(87, 'Kenya'),
(88, 'Kiribati'),
(89, 'Korea (North)'),
(90, 'Korea (South)'),
(91, 'Kosovo'),
(92, 'Kuwait'),
(93, 'Kyrgyzstan'),
(94, 'Laos'),
(95, 'Latvia'),
(96, 'Lebanon'),
(97, 'Lesotho'),
(98, 'Liberia'),
(99, 'Libya'),
(100, 'Liechtenstein'),
(101, 'Lithuania'),
(102, 'Luxembourg'),
(103, 'Madagascar'),
(104, 'Malawi'),
(105, 'Malaysia'),
(106, 'Maldives'),
(107, 'Mali'),
(108, 'Malta'),
(109, 'Marshall Islands'),
(110, 'Mauritania'),
(111, 'Mauritius'),
(112, 'Mexico'),
(113, 'Micronesia'),
(114, 'Moldova'),
(115, 'Monaco'),
(116, 'Mongolia'),
(117, 'Montenegro'),
(118, 'Morocco'),
(119, 'Mozambique'),
(120, 'Myanmar (Burma)'),
(121, 'Namibia'),
(122, 'Nauru'),
(123, 'Nepal'),
(124, 'Netherlands'),
(125, 'New Zealand'),
(126, 'Nicaragua'),
(127, 'Niger'),
(128, 'Nigeria'),
(129, 'North Macedonia'),
(130, 'Norway'),
(131, 'Oman'),
(132, 'Pakistan'),
(133, 'Palau'),
(134, 'Palestine'),
(135, 'Panama'),
(136, 'Papua New Guinea'),
(137, 'Paraguay'),
(138, 'Peru'),
(139, 'Philippines'),
(140, 'Poland'),
(141, 'Portugal'),
(142, 'Qatar'),
(143, 'Romania'),
(144, 'Russia'),
(145, 'Rwanda'),
(146, 'Saint Kitts and Nevis'),
(147, 'Saint Lucia'),
(148, 'Saint Vincent and the Grenadines'),
(149, 'Samoa'),
(150, 'San Marino'),
(151, 'Sao Tome and Principe'),
(152, 'Saudi Arabia'),
(153, 'Senegal'),
(154, 'Serbia'),
(155, 'Seychelles'),
(156, 'Sierra Leone'),
(157, 'Singapore'),
(158, 'Slovakia'),
(159, 'Slovenia'),
(160, 'Solomon Islands'),
(161, 'Somalia'),
(162, 'South Africa'),
(163, 'South Sudan'),
(164, 'Spain'),
(165, 'Sri Lanka'),
(166, 'Sudan'),
(167, 'Suriname'),
(168, 'Sweden'),
(169, 'Switzerland'),
(170, 'Syria'),
(171, 'Tajikistan'),
(172, 'Tanzania'),
(173, 'Thailand'),
(174, 'Timor-Leste'),
(175, 'Togo'),
(176, 'Tonga'),
(177, 'Trinidad and Tobago'),
(178, 'Tunisia'),
(179, 'Turkey'),
(180, 'Turkmenistan'),
(181, 'Tuvalu'),
(182, 'Uganda'),
(183, 'Ukraine'),
(184, 'United Arab Emirates'),
(185, 'United Kingdom'),
(186, 'United States'),
(187, 'Uruguay'),
(188, 'Uzbekistan'),
(189, 'Vanuatu'),
(190, 'Vatican City'),
(191, 'Venezuela'),
(192, 'Vietnam'),
(193, 'Yemen'),
(194, 'Zambia'),
(195, 'Zimbabwe');

-- Insert sample data for otps table
INSERT INTO `edugo`.`otps` (`code`, `is_used`, `expired_at`, `account_id`) 
VALUES 
('123456', false, DATE_ADD(NOW(), INTERVAL 15 MINUTE), 1),
('654321', true, DATE_ADD(NOW(), INTERVAL 15 MINUTE), 2);

-- Insert into announce_posts
INSERT INTO `edugo`.`announce_posts` (`announce_id`, `title`, `url`, `attach_file`, `close_date`, `posts_id`, `category_id`, `country_id`) 
VALUES 
(1, 'Scholarship Announcement', NULL, NULL, '2025-12-31 23:59:59', 1, 1, 1),
(2, 'New Online Course', 'www.onlinecourse.com', NULL, '2025-06-30 23:59:59', 2, 1, 2);
