-- +goose Up
ALTER TABLE compras DROP CONSTRAINT fk_compras_cliente;
ALTER TABLE compras ADD CONSTRAINT fk_compras_cliente
    FOREIGN KEY (cliente_id) REFERENCES clientes(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE compras DROP CONSTRAINT fk_compras_cliente;
ALTER TABLE compras ADD CONSTRAINT fk_compras_cliente
    FOREIGN KEY (cliente_id) REFERENCES clientes(id);

