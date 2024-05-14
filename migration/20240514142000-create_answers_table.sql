-- +migrate Down
DROP TABLE IF EXISTS `answers`;

-- +migrate Up
CREATE TABLE IF NOT EXISTS `answers` (
                                         `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
                                         `user_id` VARCHAR(36) NOT NULL,
    `question_id` INT NOT NULL,
    `user_answer` INT DEFAULT 0,
    `prefer_answer` INT DEFAULT 0,
    `importance` INT DEFAULT 0,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`,`user_id`, `question_id`),
    CONSTRAINT `fk_answers_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
