CREATE TABLE `videos` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
    `user_id` bigint NOT NULL DEFAULT 0 COMMENT 'User id',
    `play_url` varchar(128) DEFAULT './video/test.mp4' COMMENT 'video src',
    `cover_url` varchar(128) NOT NULL DEFAULT './image/test.jpg' COMMENT "cover path",
    `favorite_count` bigint NOT NULL DEFAULT 0 COMMENT "video favorite count",
    `comment_count` bigint NOT NULL DEFAULT 0 COMMENT "video comment count",
    `title` varchar(128) NOT NULL DEFAULT 'untitle' COMMENT "video title",
    PRIMARY KEY (`id`),
    KEY `idx_id` (`user_id`) COMMENT 'User id index'
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'Video information table';