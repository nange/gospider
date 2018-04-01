# noinspection SqlNoDataSourceInspectionForFile

CREATE TABLE `t_task` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `task_name` VARCHAR(64) NOT NULL,
  `task_rule_name` VARCHAR(64) NOT NULL,
  `task_desc` VARCHAR(512) NOT NULL DEFAULT '',
  `status` TINYINT NOT NULL DEFAULT 0,
  `counts` INT NOT NULL DEFAULT 0,
  `opt_user_agent` VARCHAR(128) NOT NULL DEFAULT '',
  `opt_max_depth` INT NOT NULL DEFAULT 0,
  `opt_allowed_domains` VARCHAR(512) NOT NULL DEFAULT '',
  `opt_url_filters` VARCHAR(512) NOT NULL DEFAULT '',
  `opt_allow_url_revisit` BOOL NOT NULL DEFAULT 0,
  `opt_max_body_size` INT NOT NULL DEFAULT 0,
  `opt_ignore_robots_txt` BOOL NOT NULL DEFAULT 1,
  `opt_parse_http_error_response` BOOL NOT NULL DEFAULT 0,
  `limit_enable` BOOL NOT NULL DEFAULT 0,
  `limit_domain_regexp` VARCHAR(128) NOT NULL DEFAULT '',
  `limit_domain_glob` VARCHAR(128) NOT NULL DEFAULT '',
  `limit_delay` INT NOT NULL DEFAULT 0,
  `limit_random_delay` INT NOT NULL DEFAULT 0,
  `limit_parallelism` INT NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT current_timestamp,
  `updated_at` DATETIME NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp,
  PRIMARY KEY (`id`),
  KEY `idx_task_name` (`task_name`),
  KEY `idx_task_rule_name` (`task_rule_name`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

