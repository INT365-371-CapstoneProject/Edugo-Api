-- MySQL Script generated by MySQL Workbench
-- Sat Nov  9 07:12:42 2024
-- Model: New Model    Version: 1.0
-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema edugo
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema edugo
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `edugo` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci ;
USE `edugo` ;

-- -----------------------------------------------------
-- Table `edugo`.`countries`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `edugo`.`countries` (
  `country_id` INT NOT NULL,
  `name` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`country_id`),
  UNIQUE INDEX `name_UNIQUE` (`name` ASC) VISIBLE)
ENGINE = InnoDB;

-- -----------------------------------------------------
-- Table `edugo`.`posts`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `edugo`.`posts` (
  `posts_id` INT NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(100) NOT NULL,
  `description` VARCHAR(500) NOT NULL,
  `image` LONGTEXT NULL,
  `publish_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  `posts_type` ENUM('Announce', 'Subject') NOT NULL,
  `country_id` INT NOT NULL,
  PRIMARY KEY (`posts_id`),
  UNIQUE KEY `unique_image` (`image`(255)), -- กำหนด key length 255 ตัวอักษร
  INDEX `fk_posts_countries1_idx` (`country_id` ASC) VISIBLE,
  CONSTRAINT `fk_posts_countries1`
    FOREIGN KEY (`country_id`)
    REFERENCES `edugo`.`countries` (`country_id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;



-- -----------------------------------------------------
-- Table `edugo`.`comments`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `edugo`.`comments` (
  `comments_id` INT NOT NULL AUTO_INCREMENT,
  `comments_text` VARCHAR(500) NOT NULL,
  `comments_image` LONGTEXT NULL,
  `comments_type` ENUM('Announce', 'Subject') NOT NULL,
  `publish_date` DATETIME NULL DEFAULT NULL,
  `posts_id` INT NOT NULL,
  PRIMARY KEY (`comments_id`),
  UNIQUE KEY `unique_comments_image` (`comments_image`(255)), -- กำหนด key length 255 ตัวอักษร
  INDEX `fk_comments_posts1_idx` (`posts_id` ASC) VISIBLE,
  CONSTRAINT `fk_comments_posts1`
    FOREIGN KEY (`posts_id`)
    REFERENCES `edugo`.`posts` (`posts_id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;

-- -----------------------------------------------------
-- Table `edugo`.`categories`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `edugo`.`categories` (
  `category_id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`category_id`),
  UNIQUE INDEX `name_UNIQUE` (`name` ASC) VISIBLE)
ENGINE = InnoDB;

-- -----------------------------------------------------
-- Table `edugo`.`announce_posts`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `edugo`.`announce_posts` (
  `announce_id` INT NOT NULL AUTO_INCREMENT,
  `url` VARCHAR(50) NULL,
  `attach_file` LONGTEXT NULL,
  `close_date` DATETIME NOT NULL,
  `posts_id` INT NOT NULL,
  `category_id` INT NOT NULL,
  PRIMARY KEY (`announce_id`),
  UNIQUE KEY `unique_attach_file` (`attach_file`(255)), -- กำหนด key length 255 ตัวอักษร
  INDEX `fk_announce_posts_posts1_idx` (`posts_id` ASC) VISIBLE,
  INDEX `fk_announce_posts_categories1_idx` (`category_id` ASC) VISIBLE,
  CONSTRAINT `fk_announce_posts_posts1`
    FOREIGN KEY (`posts_id`)
    REFERENCES `edugo`.`posts` (`posts_id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_announce_posts_categories1`
    FOREIGN KEY (`category_id`)
    REFERENCES `edugo`.`categories` (`category_id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
SET FOREIGN_KEY_CHECKS=1;

SET GLOBAL time_zone = '+00:00';
SET time_zone = '+00:00';