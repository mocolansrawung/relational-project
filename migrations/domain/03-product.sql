-- DROP TABLE IF EXISTS `product`

-- CREATE TABLE IF NOT EXISTS `product` (
-- 	  `id` CHAR(36) NOT NULL,
-- 	  `name` VARCHAR(255) NOT NULL,
-- 	  `shop_name` VArCHAR(255) NOT NULL,
-- 	  `price` INT NOT NULL,
-- 	  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
-- 	  `updated_at` TIMESTAMP NULL DEFAULT NULL,
--       PRIMARY KEY (`id`)
-- ) ENGINE=InnoDB
-- DEFAULT CHARSET=utf8

-- INSERT INTO `product`
-- (`id`, `name`, `shop_name`, `price`, `created_at`)
-- VALUES
-- ('123e4567-e89b-12d3-a456-426614174000', 'iPhone 12', 'Apple Store', 69900, NOW()),
-- ('123e4567-e89b-12d3-a456-426614174001', 'Galaxy S21', 'Samsung Store', 79900, NOW()),
-- ('123e4567-e89b-12d3-a456-426614174002', 'Pixel 5', 'Google Store', 69900, NOW()),
-- ('123e4567-e89b-12d3-a456-426614174003', 'OnePlus 9', 'OnePlus Store', 72900, NOW()),
-- ('123e4567-e89b-12d3-a456-426614174004', 'Huawei P50', 'Huawei Store', 59900, NOW());

DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `brands`;
DROP TABLE IF EXISTS `products`;
DROP TABLE IF EXISTS `variants`;
DROP TABLE IF EXISTS `images`;
DROP TABLE IF EXISTS `warehouse`;
DROP TABLE IF EXISTS `variant_warehouse`;

CREATE TABLE users (
    user_id CHAR(36) PRIMARY KEY NOT NULL,
    user_type ENUM('ADMIN', 'REGULAR') NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by CHAR(36) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by CHAR(36) NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by CHAR(36) NULL DEFAULT NULL
);

CREATE TABLE brands (
    brand_id CHAR(36) PRIMARY KEY NOT NULL,
    brand_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by CHAR(36) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by CHAR(36) NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by CHAR(36) NULL DEFAULT NULL
);

CREATE TABLE products (
    product_id CHAR(36) PRIMARY KEY NOT NULL,
    brand_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    product_name VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by CHAR(36) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by CHAR(36) NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by CHAR(36) NULL DEFAULT NULL,
    FOREIGN KEY (brand_id) REFERENCES brands(brand_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;

CREATE INDEX idx_products_brand ON products(brand_id);
CREATE INDEX idx_products_user ON products(user_id);

CREATE TABLE variants (
    variant_id CHAR(36) PRIMARY KEY NOT NULL,
    product_id CHAR(36) NOT NULL,
    variant_name VARCHAR(255),
    price INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by CHAR(36) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by CHAR(36) NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by CHAR(36) NULL DEFAULT NULL,
    FOREIGN KEY (product_id) REFERENCES products(product_id)
);

CREATE INDEX idx_variants_product ON variants(product_id);

CREATE TABLE images (
    image_id CHAR(36) PRIMARY KEY NOT NULL,
    variant_id CHAR(36) NOT NULL,
    image_url VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by CHAR(36) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by CHAR(36) NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by CHAR(36) NULL DEFAULT NULL,
    FOREIGN KEY (variant_id) REFERENCES variants(variant_id)
);

CREATE INDEX idx_images_variant ON images(variant_id);

CREATE TABLE warehouse (
    warehouse_id CHAR(36) PRIMARY KEY,
	warehouse_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by CHAR(36) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by CHAR(36) NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by CHAR(36) NULL DEFAULT NULL
);

CREATE TABLE variant_warehouse (
    variant_id CHAR(36) NOT NULL,
    warehouse_id CHAR(36) NOT NULL,
    stock INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by CHAR(36) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by CHAR(36) NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by CHAR(36) NULL DEFAULT NULL,
    FOREIGN KEY (variant_id) REFERENCES variants(variant_id),
    FOREIGN KEY (warehouse_id) REFERENCES warehouse(warehouse_id)
);

CREATE INDEX idx_variant_warehouse_variant ON variant_warehouse(variant_id);
CREATE INDEX idx_variant_warehouse_warehouse ON variant_warehouse(warehouse_id);
CREATE INDEX idx_variant_warehouse_composite ON variant_warehouse(variant_id, warehouse_id);

DELIMITER //
CREATE TRIGGER update_product_updated_time
AFTER UPDATE ON variants
FOR EACH ROW
BEGIN
   UPDATE products 
   SET updated_at = NOW() 
   WHERE product_id = OLD.product_id;
END; //
DELIMITER ;
