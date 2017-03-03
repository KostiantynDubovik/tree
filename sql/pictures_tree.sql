DROP DATABASE IF EXISTS `pictures_tree`;

CREATE DATABASE `pictures_tree` /*!40100 DEFAULT CHARACTER SET utf8 */;

CREATE TABLE `pictures_tree`.`pictures_tree_nodes` (
    `node_id` INT(10) NOT NULL AUTO_INCREMENT,
    `left_key` INT(10) NOT NULL DEFAULT 0,
    `right_key` INT(10) NOT NULL DEFAULT 0,
    `parent_id` INT(10) NOT NULL DEFAULT 0,
    `image_name` VARCHAR(45) NOT NULL DEFAULT '',
    PRIMARY KEY `node_id` (`node_id` ASC),
    UNIQUE INDEX `node_id` (`node_id` ASC)
);

INSERT INTO `pictures_tree`.`pictures_tree_nodes` (`node_id`, `left_key`, `right_key`, `parent_id`, `image_name`) VALUES ('1', '1', '2', '0', 'nil');


DELIMITER // 
   
CREATE PROCEDURE `pictures_tree`.`insert_node` (IN `recieved_left_key` INT,`recieved_right_key` INT, `recieved_level` INT, `recieved_image_name` VARCHAR(45))	 
LANGUAGE SQL 
DETERMINISTIC 
SQL SECURITY DEFINER
COMMENT 'insert_node procedure' 
BEGIN 
    UPDATE `pictures_tree`.`pictures_tree_nodes` SET `right_key` = `right_key` + 2, `left_key` = IF(`left_key` > `recieved_left_key`, `left_key` + 2, `left_key`) WHERE `right_key` >= `recieved_right_key`;
    INSERT INTO `pictures_tree`.`pictures_tree_nodes` SET `left_key` = `recieved_right_key`, `right_key` = `recieved_right_key` + 1, `level` = `recieved_level` + 1, `image_name`=`recieved_image_name`;
END// 

CREATE PROCEDURE `pictures_tree`.`delete_node` (IN `recieved_left_key` INT,`recieved_right_key` INT, `recieved_level` INT, `image_name` VARCHAR(45))	 
LANGUAGE SQL 
DETERMINISTIC 
SQL SECURITY DEFINER 
COMMENT 'delete_node procedure' 
BEGIN 
	DELETE FROM `pictures_tree`.`pictures_tree_nodes` WHERE `left_key` >= `recieved_left_key` AND `right_key` <= `recieved_right_key`;
	UPDATE `pictures_tree`.`pictures_tree_nodes` SET `left_key` = IF(`left_key` > `recieved_left_key`, `left_key` - (`recieved_right_key` - `recieved_left_key` + 1), `left_key`), `right_key` = `right_key` - (`recieved_right_key` - `recieved_left_key` + 1) WHERE `right_key` > `recieved_right_key`;
END//


