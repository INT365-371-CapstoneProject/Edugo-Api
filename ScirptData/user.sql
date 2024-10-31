create user edugo@localhost identified by 'mysql@edugo';

GRANT ALL PRIVILEGES ON edugo.* TO 'edugo'@'localhost';

FLUSH PRIVILEGES;
