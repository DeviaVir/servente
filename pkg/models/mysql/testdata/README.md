```
CREATE DATABASE test_servente CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'test_web'@'%';
GRANT CREATE, DROP, ALTER, INDEX, SELECT, INSERT, UPDATE, DELETE ON test_servente.* TO 'test_web'@'%';
ALTER USER 'test_web'@'%' IDENTIFIED BY 'pass';
```
