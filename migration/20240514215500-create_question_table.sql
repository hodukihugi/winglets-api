-- +migrate Down
DROP TABLE IF EXISTS `questions`;

-- +migrate Up
CREATE TABLE IF NOT EXISTS `questions` (
   `question_id` int AUTO_INCREMENT NOT NULL,
   `content` TEXT NOT NULL,
   `answers` TEXT NOT NULL,
   `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
   `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   `deleted_at` DATETIME DEFAULT NULL,
   PRIMARY KEY (`question_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;