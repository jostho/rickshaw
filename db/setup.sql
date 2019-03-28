
-- connect as root

CREATE DATABASE app;

CREATE USER 'appuser'@'%' IDENTIFIED BY 'appsecret';

GRANT ALL PRIVILEGES ON app.* TO 'appuser'@'%';
