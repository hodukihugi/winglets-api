-- +migrate Down
DROP TABLE IF EXISTS `matches`;

-- +migrate Up
CREATE TABLE IF NOT EXISTS `matches` (
    `matcher_id` VARCHAR(36) NOT NULL,
    `matchee_id` VARCHAR(36) NOT NULL,
    `match_status` INT DEFAULT 0,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`matcher_id`, `matchee_id`),
    CONSTRAINT `fk_match_id` FOREIGN KEY (`matcher_id`, `matchee_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
