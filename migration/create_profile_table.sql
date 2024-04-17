-- +migrate Down
DROP TABLE IF EXISTS `profiles`;

-- +migrate Up
CREATE TABLE IF NOT EXISTS `profiles` (
    `id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `gender` VARCHAR(10),
    `birthday` DATETIME,
    `height` DECIMAL(5,2),
    `horoscope` VARCHAR(50),
    `hobby` VARCHAR(255),
    `language` VARCHAR(50),
    `education` VARCHAR(100),
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_profile_user_id` FOREIGN KEY (`id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
