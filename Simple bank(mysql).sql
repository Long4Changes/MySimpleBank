-- Table for accounts
CREATE TABLE `accounts` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `owner` VARCHAR(255) NOT NULL,
  `balance` BIGINT NOT NULL,
  `currency` VARCHAR(255) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

-- Table for account entries (transactions)
CREATE TABLE `entries` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `account_id` BIGINT UNSIGNED NOT NULL,
  `amount` BIGINT NOT NULL COMMENT 'can be negative or positive',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

-- Table for transfers between accounts
CREATE TABLE `transfers` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `from_account_id` BIGINT UNSIGNED NOT NULL,
  `to_account_id` BIGINT UNSIGNED NOT NULL,
  `amount` BIGINT NOT NULL COMMENT 'must be positive',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

-- Indexes for faster lookups
CREATE INDEX `accounts_owner_index` ON `accounts` (`owner`);
CREATE INDEX `entries_account_id_index` ON `entries` (`account_id`);
CREATE INDEX `transfers_from_account_id_index` ON `transfers` (`from_account_id`);
CREATE INDEX `transfers_to_account_id_index` ON `transfers` (`to_account_id`);
CREATE INDEX `transfers_from_to_account_id_index` ON `transfers` (`from_account_id`, `to_account_id`);

-- Foreign key constraints to ensure data integrity
ALTER TABLE `entries` ADD CONSTRAINT `fk_entries_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`);
ALTER TABLE `transfers` ADD CONSTRAINT `fk_transfers_from_account` FOREIGN KEY (`from_account_id`) REFERENCES `accounts` (`id`);
ALTER TABLE `transfers` ADD CONSTRAINT `fk_transfers_to_account` FOREIGN KEY (`to_account_id`) REFERENCES `accounts` (`id`);