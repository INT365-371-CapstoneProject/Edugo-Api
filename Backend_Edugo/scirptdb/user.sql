create user edugo@localhost identified by 'mysql@edugo';

GRANT ALL PRIVILEGES ON edugo_v3.* TO 'edugo'@'localhost';

FLUSH PRIVILEGES;
