-- +goose Up
ALTER TABLE compra_items
    DROP CONSTRAINT fk_compras_itens;

ALTER TABLE compra_items
    ADD CONSTRAINT fk_compras_itens
        FOREIGN KEY (compra_id) REFERENCES compras(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE compra_items
    DROP CONSTRAINT fk_compras_itens;

ALTER TABLE compra_items
    ADD CONSTRAINT fk_compras_itens
        FOREIGN KEY (compra_id) REFERENCES compras(id);
