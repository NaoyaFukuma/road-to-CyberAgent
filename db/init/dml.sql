use CA_Tech_Dojo;

SET NAMES utf8mb4;

INSERT INTO `game_settings` (`gacha_coin_consumption`, `ranking_list_limit`, `n_weight`, `r_weight`, `sr_weight`, `max_gacha_times`) VALUES (100, 10, 5, 3, 1, 30);
INSERT INTO `game_settings` (`gacha_coin_consumption`, `ranking_list_limit`, `n_weight`, `r_weight`, `sr_weight`, `max_gacha_times`, `is_active`) VALUES (10, 10, 5, 3, 1, 50, true);

INSERT INTO `item` (`name`, `rarity`) VALUES ('ノーマル1', 1);
INSERT INTO `item` (`name`, `rarity`) VALUES ('ノーマル2', 1);
INSERT INTO `item` (`name`, `rarity`) VALUES ('ノーマル3', 1);
INSERT INTO `item` (`name`, `rarity`) VALUES ('ノーマル4', 1);
INSERT INTO `item` (`name`, `rarity`) VALUES ('ノーマル5', 1);
INSERT INTO `item` (`name`, `rarity`) VALUES ('ノーマル6', 1);
INSERT INTO `item` (`name`, `rarity`) VALUES ('ノーマル7', 1);
INSERT INTO `item` (`name`, `rarity`) VALUES ('レア1', 2);
INSERT INTO `item` (`name`, `rarity`) VALUES ('レア2', 2);
INSERT INTO `item` (`name`, `rarity`) VALUES ('レア3', 2);
INSERT INTO `item` (`name`, `rarity`) VALUES ('レア4', 2);
INSERT INTO `item` (`name`, `rarity`) VALUES ('レア5', 2);
INSERT INTO `item` (`name`, `rarity`) VALUES ('スーパーレア1', 3);
INSERT INTO `item` (`name`, `rarity`) VALUES ('スーパーレア2', 3);
