-- +migrate Down
DROP TABLE IF EXISTS `users`;

-- +migrate Up
CREATE TABLE IF NOT EXISTS `users` (
    `id` VARCHAR(36) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `verification_code` VARCHAR(30) DEFAULT NULL,
    `verification_status` INT DEFAULT 0,
    `verification_time` DATETIME DEFAULT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT email_unique UNIQUE(email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;