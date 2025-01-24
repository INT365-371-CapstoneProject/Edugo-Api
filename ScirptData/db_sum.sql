-- MySQL Script generated by MySQL Workbench
-- Tue Jan 14 15:37:01 2025
-- Model: New Model    Version: 1.0
-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
-- -----------------------------------------------------
-- Schema edugo
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema edugo
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `edugo` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci ;
USE `edugo` ;

-- -----------------------------------------------------
-- Table `edugo`.`accounts`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `edugo`.`accounts` (
  `account_id` INT NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(50) NOT NULL,
  `password` VARCHAR(500) NOT NULL,
  `email` VARCHAR(50) NOT NULL,
  `create_on` DATETIME NOT NULL,
  `last_login` DATETIME NULL DEFAULT NULL,
  `update_on` DATETIME NOT NULL,
  `role` ENUM('admin', 'user', 'provider', 'superadmin') NOT NULL,
  PRIMARY KEY (`account_id`),
  UNIQUE INDEX `username` (`username` ASC) VISIBLE,
  UNIQUE INDEX `email` (`email` ASC) VISIBLE)
ENGINE = InnoDB
AUTO_INCREMENT = 4
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;

-- -----------------------------------------------------
-- Table `edugo`.`otps`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `edugo`.`otps` (
  `otp_id` INT NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(6) NOT NULL,
  `is_used` BOOLEAN DEFAULT FALSE,
  `attempt_count` INT DEFAULT 0,
  `expired_at` DATETIME NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `account_id` INT NOT NULL,
  PRIMARY KEY (`otp_id`),
  INDEX `idx_account_code` (`account_id`, `code`),
  INDEX `idx_expired_at` (`expired_at`),
  CONSTRAINT `fk_otps_accounts`
    FOREIGN KEY (`account_id`)
    REFERENCES `edugo`.`accounts` (`account_id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;

-- -----------------------------------------------------
-- Table `edugo`.`admins`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `edugo`.`admins` (
  `admin_id` INT NOT NULL AUTO_INCREMENT,
  `firstname` VARCHAR(50) NOT NULL,
  `lastname` VARCHAR(50) NOT NULL,
  `status` ENUM('Active', 'Inactive') NOT NULL,
  `account_id` INT NOT NULL,
  PRIMARY KEY (`admin_id`),
  INDEX `fk_admins_accounts1_idx` (`account_id` ASC) VISIBLE,
  CONSTRAINT `fk_admins_accounts1`
    FOREIGN KEY (`account_id`)
    REFERENCES `edugo`.`accounts` (`account_id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB
AUTO_INCREMENT = 2
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `edugo`.`providers`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `edugo`.`providers` (
  `provider_id` INT NOT NULL AUTO_INCREMENT,
  `company_name` VARCHAR(50) NOT NULL,
  `url` VARCHAR(500) NOT NULL,
  `address` VARCHAR(100) NOT NULL,
  `status` ENUM('Active', 'Inactive') NOT NULL,
  `phone` VARCHAR(10) NOT NULL,
  `verify` ENUM('Y', 'N') NOT NULL,
  `account_id` INT NOT NULL,
  PRIMARY KEY (`provider_id`),
  INDEX `fk_providers_accounts1_idx` (`account_id` ASC) VISIBLE,
  CONSTRAINT `fk_providers_accounts1`
    FOREIGN KEY (`account_id`)
    REFERENCES `edugo`.`accounts` (`account_id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB
AUTO_INCREMENT = 3
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


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
  `description` VARCHAR(3000) NOT NULL,
  `image` LONGBLOB NULL DEFAULT NULL,
  `publish_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  `posts_type` ENUM('Announce', 'Subject') NOT NULL,
  `account_id` INT NOT NULL,
  PRIMARY KEY (`posts_id`),
  UNIQUE INDEX `unique_image` (`image`(255) ASC) VISIBLE,
  INDEX `fk_posts_accounts1_idx` (`account_id` ASC) VISIBLE,
  CONSTRAINT `fk_posts_accounts1`
    FOREIGN KEY (`account_id`)
    REFERENCES `edugo`.`accounts` (`account_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `edugo`.`comments`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `edugo`.`comments` (
  `comments_id` INT NOT NULL AUTO_INCREMENT,
  `comments_text` VARCHAR(3000) NOT NULL,
  `comments_image` LONGBLOB NULL DEFAULT NULL,
  `publish_date` DATETIME NULL DEFAULT NULL,
  `posts_id` INT NOT NULL,
  `account_id` INT NOT NULL,
  PRIMARY KEY (`comments_id`),
  UNIQUE INDEX `unique_comments_image` (`comments_image`(255) ASC) VISIBLE,
  INDEX `fk_comments_posts1_idx` (`posts_id` ASC) VISIBLE,
  INDEX `fk_comments_accounts1_idx` (`account_id` ASC) VISIBLE,
  CONSTRAINT `fk_comments_posts1`
    FOREIGN KEY (`posts_id`)
    REFERENCES `edugo`.`posts` (`posts_id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_comments_accounts1`
    FOREIGN KEY (`account_id`)
    REFERENCES `edugo`.`accounts` (`account_id`)
    ON DELETE NO ACTION
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
  `title` VARCHAR(100) NOT NULL,
  `url` VARCHAR(255) NULL DEFAULT NULL,
  `attach_name` VARCHAR(255) NULL DEFAULT NULL,
  `attach_file` LONGBLOB NULL DEFAULT NULL,
  `close_date` DATETIME NOT NULL,
  `posts_id` INT NOT NULL,
  `category_id` INT NOT NULL,
  `country_id` INT NOT NULL,
  PRIMARY KEY (`announce_id`),
  UNIQUE INDEX `unique_attach_file` (`attach_file`(255) ASC) VISIBLE,
  INDEX `fk_announce_posts_posts1_idx` (`posts_id` ASC) VISIBLE,
  INDEX `fk_announce_posts_categories1_idx` (`category_id` ASC) VISIBLE,
  INDEX `fk_announce_posts_countries1_idx` (`country_id` ASC) VISIBLE,
  CONSTRAINT `fk_announce_posts_posts1`
    FOREIGN KEY (`posts_id`)
    REFERENCES `edugo`.`posts` (`posts_id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_announce_posts_categories1`
    FOREIGN KEY (`category_id`)
    REFERENCES `edugo`.`categories` (`category_id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_announce_posts_countries1`
    FOREIGN KEY (`country_id`)
    REFERENCES `edugo`.`countries` (`country_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
SET FOREIGN_KEY_CHECKS=1;

SET GLOBAL time_zone = '+00:00';
SET time_zone = '+00:00';
