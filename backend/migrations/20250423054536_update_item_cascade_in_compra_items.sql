-- +goose Up
ALTER TABLE compra_items
    DROP CONSTRAINT fk_compra_items_item;

ALTER TABLE compra_items
    ADD CONSTRAINT fk_compra_items_item
        FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE compra_items
    DROP CONSTRAINT fk_compra_items_item;

ALTER TABLE compra_items
    ADD CONSTRAINT fk_compra_items_item
        FOREIGN KEY (item_id) REFERENCES items(id);
