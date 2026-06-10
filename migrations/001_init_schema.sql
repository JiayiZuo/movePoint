CREATE DATABASE IF NOT EXISTS `movepoint`
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE `movepoint`;

CREATE TABLE IF NOT EXISTS `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(191) NOT NULL,
  `email` varchar(191) NOT NULL,
  `password` longtext NOT NULL,
  `weight` double DEFAULT NULL,
  `height` double DEFAULT NULL,
  `birth_date` datetime(3) DEFAULT NULL,
  `avatar_url` longtext,
  `bio` text,
  `achievements` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `climbing_records` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` int unsigned NOT NULL,
  `type` varchar(20) NOT NULL,
  `start_time` datetime(3) DEFAULT NULL,
  `end_time` datetime(3) DEFAULT NULL,
  `duration` bigint DEFAULT NULL,
  `grade` varchar(10) DEFAULT NULL,
  `color` varchar(20) DEFAULT NULL,
  `attempts` varchar(10) DEFAULT NULL,
  `success` boolean DEFAULT NULL,
  `rating` bigint DEFAULT NULL,
  `location` varchar(255) DEFAULT NULL,
  `notes` text,
  `media_urls` text,
  `calories` double DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_climbing_records_deleted_at` (`deleted_at`),
  KEY `idx_climbing_records_user_id` (`user_id`),
  CONSTRAINT `chk_climbing_records_rating` CHECK (`rating` >= 1 AND `rating` <= 5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `climbing_analyses` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `date` datetime(3) DEFAULT NULL,
  `data` text,
  PRIMARY KEY (`id`),
  KEY `idx_climbing_analyses_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
