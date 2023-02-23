CREATE TABLE `users` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
    `name` varchar(128) NOT NULL DEFAULT '' COMMENT 'User name',
    `passwd_hash` char(32) NOT NULL DEFAULT '' COMMENT 'User password hash',
    `signature` varchar(128) DEFAULT '' COMMENT 'User signature',
    `follow_count` bigint unsigned NOT NULL DEFAULT 0 COMMENT "User follow count",
    `follower_count` bigint unsigned NOT NULL DEFAULT 0 COMMENT 'User follwer count',
    `avatar` varchar(128) NOT NULL DEFAULT '/test.jpg' COMMENT "User avator",
    `background_image` varchar(128) NOT NULL DEFAULT '/testBG.jpg' COMMENT "User background",
    `total_favorited` bigint NOT NULL DEFAULT 0 COMMENT 'User total_favorited',
    `work_count` bigint NOT NULL DEFAULT 0 COMMENT 'User work_count',
    `favorite_count` bigint NOT NULL DEFAULT 0 COMMENT 'favorite_count',
    PRIMARY KEY (`id`),
    KEY `idx_name` (`name`) COMMENT 'User name index'
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'User information table';

CREATE TABLE `user_follows` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id` bigint DEFAULT '0' COMMENT 'user id',
    `followed_id` bigint DEFAULT '0' COMMENT 'followed user ID',
    `status` tinyint(1) DEFAULT '0' COMMENT 'follow status: 0 for cancel, 1 for follow',
    PRIMARY KEY (`id`),
    UNIQUE KEY `vip_followed_indx` (`user_id`, `followed_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'User follow table';

CREATE TABLE `favorites` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id` bigint DEFAULT '0' COMMENT 'user id',
    `video_id` bigint DEFAULT '0' COMMENT 'followed user ID',
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_video` (`user_id`, `video_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'User favorites video';

CREATE TABLE `user_comments` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id` bigint DEFAULT '0' COMMENT 'user id',
    `video_id` bigint DEFAULT '0' COMMENT 'video ID',
    `content` varchar(1000) NOT NULL DEFAULT "",
    `create_date` varchar(128),
    PRIMARY KEY (`id`),
    UNIQUE KEY `comments_user_video` (`user_id`, `video_id`, `create_date`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'User comments table';