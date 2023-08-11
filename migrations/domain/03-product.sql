DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `brands`;
DROP TABLE IF EXISTS `products`;
DROP TABLE IF EXISTS `variants`;
DROP TABLE IF EXISTS `images`;
DROP TABLE IF EXISTS `warehouse`;
DROP TABLE IF EXISTS `variant_warehouse`;

CREATE TABLE `users` (
    `id` CHAR(36) PRIMARY KEY NOT NULL,
    `type` ENUM('ADMIN', 'REGULAR') NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` CHAR(36) NOT NULL,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` CHAR(36) NOT NULL,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    `deleted_by` CHAR(36) NULL DEFAULT NULL
);

CREATE TABLE `brands` (
    `id` CHAR(36) PRIMARY KEY NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` CHAR(36) NOT NULL,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` CHAR(36) NOT NULL,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    `deleted_by` CHAR(36) NULL DEFAULT NULL
);

CREATE TABLE `products` (
    `id` CHAR(36) PRIMARY KEY NOT NULL,
    `brand_id` CHAR(36) NOT NULL,
    `user_id` CHAR(36) NOT NULL,
    `name` VARCHAR(255),
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` CHAR(36) NOT NULL,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` CHAR(36) NOT NULL,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    `deleted_by` CHAR(36) NULL DEFAULT NULL,
    FOREIGN KEY (`brand_id`) REFERENCES `brands`(`brand_id`),
    FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;

CREATE INDEX `idx_products_brand` ON `products`(`brand_id`);
CREATE INDEX `idx_products_user` ON `products`(`user_id`);

CREATE TABLE `variants` (
    `id` CHAR(36) PRIMARY KEY NOT NULL,
    `product_id` CHAR(36) NOT NULL,
    `name` VARCHAR(255),
    `stock` INT NULL,
    `price` INT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` CHAR(36) NOT NULL,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` CHAR(36) NOT NULL,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    `deleted_by` CHAR(36) NULL DEFAULT NULL,
    FOREIGN KEY (`product_id`) REFERENCES `products`(`product_id`)
);

CREATE INDEX `idx_variants_product` ON `variants`(`product_id`);

CREATE TABLE `images` (
    `id` CHAR(36) PRIMARY KEY NOT NULL,
    `product_id` CHAR(36) NOT NULL,
    `image_url` VARCHAR(255),
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` CHAR(36) NOT NULL,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` CHAR(36) NOT NULL,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    `deleted_by` CHAR(36) NULL DEFAULT NULL,
    FOREIGN KEY (`variant_id`) REFERENCES `variants`(`variant_id`)
);

CREATE INDEX `idx_images_variant` ON `images`(`variant_id`);

CREATE TABLE `warehouse` (
    `id` CHAR(36) PRIMARY KEY,
	`name` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` CHAR(36) NOT NULL,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` CHAR(36) NOT NULL,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    `deleted_by` CHAR(36) NULL DEFAULT NULL
);

CREATE TABLE `variant_warehouse` (
    `variant_id` CHAR(36) NOT NULL,
    `warehouse_id` CHAR(36) NOT NULL,
    `stock` INT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` CHAR(36) NOT NULL,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` CHAR(36) NOT NULL,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    `deleted_by` CHAR(36) NULL DEFAULT NULL,
    FOREIGN KEY (`variant_id`) REFERENCES `variants`(`variant_id`),
    FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse`(`warehouse_id`)
);

CREATE INDEX `idx_variant_warehouse_variant` ON `variant_warehouse`(`variant_id`);
CREATE INDEX `idx_variant_warehouse_warehouse` ON `variant_warehouse`(`warehouse_id`);
CREATE INDEX `idx_variant_warehouse_composite` ON `variant_warehouse`(`variant_id`, warehouse_id);

DELIMITER //
CREATE TRIGGER `update_product_updated_time`
AFTER UPDATE ON `variants`
FOR EACH ROW
BEGIN
   UPDATE `products` 
   SET `updated_at` = NOW() 
   WHERE `product_id` = OLD.`product_id`;
END; //
DELIMITER ;
