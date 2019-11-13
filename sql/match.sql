USE lottery;


CREATE TABLE `sport` (
  `id` bigint(20) unsigned NOT NULL,
  `name` varchar(20)  NOT NULL,
   PRIMARY KEY (`id`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



CREATE TABLE `match_01` (
  `id` bigint(20) unsigned NOT NULL,
  `sport_id` bigint(20) unsigned NOT NULL,
  `home_name` varchar(20)  NOT NULL,
  `away_name` varchar(20)  NOT NULL,
  `time` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_sport` (`sport_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE `odds_01` (
  `id` bigint(20) unsigned NOT NULL,
  `match_id` bigint(20) unsigned NOT NULL,
  `bet_type_name` varchar(10) NOT NULL,
  `selection1` decimal  NOT NULL,
  `selection2` decimal  NOT NULL,
  `time` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_match_id` (`match_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

