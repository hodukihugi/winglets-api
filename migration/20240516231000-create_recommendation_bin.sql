-- +migrate Down
DROP TABLE IF EXISTS `recommendation_bins`;

-- +migrate Up
CREATE TABLE IF NOT EXISTS `recommendation_bins` (
    `user_id` VARCHAR(36) NOT NULL,
    `recommended_user_id` VARCHAR(36) NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`user_id`),
    CONSTRAINT `fk_recommended_bin_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_recommended_bin_recommended_user_id` FOREIGN KEY (`recommended_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
