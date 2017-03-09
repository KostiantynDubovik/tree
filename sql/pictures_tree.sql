DROP DATABASE IF EXISTS `images_tree`;

CREATE DATABASE `images_tree` /*!40100 DEFAULT CHARACTER SET utf8 */;

CREATE TABLE `images_tree`.`images_tree_nodes` (
  `nodeId`    INT(10)     NOT NULL AUTO_INCREMENT,
  `leftKey`   INT(10)     NOT NULL DEFAULT 0,
  `rightKey`  INT(10)     NOT NULL DEFAULT 0,
  `level`     INT(10)     NOT NULL DEFAULT 0,
  `parentId`  INT(10)     NOT NULL DEFAULT 0,
  `imageName` VARCHAR(45) NOT NULL DEFAULT '',
  PRIMARY KEY (`nodeId`),
  UNIQUE INDEX `nodeId_UNIQUE` (`nodeId` ASC)
);

INSERT INTO `images_tree`.`images_tree_nodes` (`nodeId`, `leftKey`, `rightKey`, `level`, `imageName`)
VALUES (1, 1, 2, 0, 'nil');


DELIMITER //

CREATE PROCEDURE `images_tree`.`insert_node`(IN `node_id` INT, `parent_id` INT, `recieved_image_name` VARCHAR(45))
LANGUAGE SQL
DETERMINISTIC
  SQL SECURITY DEFINER
  COMMENT 'insert_node procedure'
  BEGIN
    DECLARE `rc` INT;
    DECLARE `lvl` INT;
    SET `rc` = (SELECT `rightKey`
                FROM `images_tree`.`images_tree_nodes`
                WHERE `nodeId` = `parent_id`);
    SET `lvl` = (SELECT `level`
                 FROM `images_tree`.`images_tree_nodes`
                 WHERE `nodeId` = `parent_id`);


    UPDATE `images_tree`.`images_tree_nodes`
    SET `rightKey` = `rightKey` + 2, `leftKey` = IF(`leftKey` > `rc`, `leftKey` + 2, `leftKey`)
    WHERE `rightKey` >= `rc`;
    INSERT INTO `images_tree`.`images_tree_nodes`
    SET `nodeId` = `node_id`, `leftKey` = `rc`, `rightKey` = `rc` + 1,
      `level`     = `lvl` + 1, `imageName` = `recieved_image_name`;
  END//

CREATE PROCEDURE `images_tree`.`delete_node`(IN `node_id` INT)
LANGUAGE SQL
DETERMINISTIC
  SQL SECURITY DEFINER
  COMMENT 'delete_node procedure'
  BEGIN
    DECLARE `rc` INT;
    DECLARE `lc` INT;
    SET `rc` = (SELECT `rightKey`
                FROM `images_tree`.`images_tree_nodes`
                WHERE `nodeId` = `node_id`);
    SET `lc` = (SELECT `leftKey`
                FROM `images_tree`.`images_tree_nodes`
                WHERE `nodeId` = `node_id`);

    DELETE FROM `images_tree`.`images_tree_nodes`
    WHERE `leftKey` >= `lc` AND `rightKey` <= `rc`;

    UPDATE `images_tree`.`images_tree_nodes`
    SET `leftKey` = IF(`leftKey` > `lc`, `leftKey` - (`rc` - lc + 1), `leftKey`),
      `rightKey`  = `rightKey` - (`rc` - `lc` + 1)
    WHERE `rightKey` > `rc`;
  END//

CREATE PROCEDURE `images_tree`.`delete_all_nodes`()
LANGUAGE SQL
DETERMINISTIC
  SQL SECURITY DEFINER
  COMMENT 'delete_all_nodes procedure'
  BEGIN
    DELETE FROM `images_tree`.`images_tree_nodes`
    WHERE `leftKey` > 1;
    UPDATE `images_tree`.`images_tree_nodes`
    SET `rightKey` = 2
    WHERE `leftKey` = 1;
  END//