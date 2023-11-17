CREATE DATABASE IF NOT EXISTS `CA_Tech_Dojo` DEFAULT CHARACTER SET utf8mb4 ;
USE `CA_Tech_Dojo` ;

SET CHARSET utf8mb4;

CREATE TABLE IF NOT EXISTS `user` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'ユーザID',
  `name` VARCHAR(128) NOT NULL COMMENT '名前',
  `high_score` INT NOT NULL DEFAULT 0 COMMENT 'ハイスコア',
  `coin` INT NOT NULL DEFAULT 0 COMMENT '所持コイン数',
  `auth_token` VARCHAR(128) NOT NULL COMMENT 'UUIDを用いた認証用トークン',
  PRIMARY KEY (`id`))
ENGINE = InnoDB
COMMENT = 'ユーザ';

CREATE TABLE IF NOT EXISTS `game_settings` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'ゲーム設定のID',
  `gacha_coin_consumption` INT NOT NULL COMMENT 'ガチャ1回あたりのコイン消費量',
  `ranking_list_limit` INT NOT NULL COMMENT 'ランキングリスト取得時のユーザ数上限',
  `n_weight` INT NOT NULL COMMENT 'Nアイテムの重み',
  `r_weight` INT NOT NULL COMMENT 'Rアイテムの重み',
  `sr_weight` INT NOT NULL COMMENT 'SRアイテムの重み',
  `max_gacha_times` INT NOT NULL COMMENT 'ガチャの最大回数',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '作成日時',
  `is_active` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '現在適用されている設定かどうか',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ゲーム設定情報';

CREATE TABLE IF NOT EXISTS `item` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'アイテムID',
  `name` VARCHAR(128) NOT NULL COMMENT '名称',
  `rarity` INT NOT NULL COMMENT 'レアリティ(1=N, 2=R, 3=SR)',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='アイテムマスター';

CREATE TABLE IF NOT EXISTS `user_items` (
  `user_id` INT NOT NULL COMMENT 'user.id',
  `item_id` INT NOT NULL COMMENT 'item.id',
  PRIMARY KEY (`user_id`, `item_id`),
  FOREIGN KEY (`user_id`) REFERENCES `user`(`id`),
  FOREIGN KEY (`item_id`) REFERENCES `item`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ユーザが保持するアイテム';

CREATE TABLE IF NOT EXISTS `user_scores` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` INT NOT NULL COMMENT 'user.id',
  `score` INT NOT NULL COMMENT 'スコア',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  PRIMARY KEY (`id`),
  FOREIGN KEY (`user_id`) REFERENCES `user`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ユーザのスコア（ここからランキングを算出する）';