USE lottery;

CREATE TABLE `tran` (
  `id` bigint(20) unsigned NOT NULL,
  `user_id` bigint(20) unsigned NOT NULL,
  `amount` int(11) NOT NULL,
  `type` tinyint(3) unsigned NOT NULL,
  `time` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_tran_user_id_type_time` (`user_id`,`type`,`time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
