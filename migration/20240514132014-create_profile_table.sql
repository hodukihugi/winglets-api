-- +migrate Down
DROP TABLE IF EXISTS `profiles`;
-- +migrate Up
CREATE TABLE IF NOT EXISTS `profiles` (
    `id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `gender` VARCHAR(10),
    `birthday` DATETIME,
    `height` INT,
    `horoscope` VARCHAR(50),
    `hobby` TEXT,
    `language` TEXT,
    `education` VARCHAR(255),
    `home_town` VARCHAR(255),
    `coordinates` VARCHAR(255),
    `image_id_1` VARCHAR(50) DEFAULT NULL,
    `image_id_2` VARCHAR(50) DEFAULT NULL,
    `image_id_3` VARCHAR(50) DEFAULT NULL,
    `image_id_4` VARCHAR(50) DEFAULT NULL,
    `image_id_5` VARCHAR(50) DEFAULT NULL,
    `image_url_1` TEXT DEFAULT NULL,
    `image_url_2` TEXT DEFAULT NULL,
    `image_url_3` TEXT DEFAULT NULL,
    `image_url_4` TEXT DEFAULT NULL,
    `image_url_5` TEXT DEFAULT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_profile_user_id` FOREIGN KEY (`id`) REFERENCES `users` (`id`) ON DELETE CASCADE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
