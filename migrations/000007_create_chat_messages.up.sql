CREATE TABLE IF NOT EXISTS chat_messages
(
    id                      INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    sender_id               INTEGER NOT NULL,
    receiver_id             INTEGER NOT NULL,
    content                 VARCHAR(250) NOT NULL,
    is_read                 BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at              TIMESTAMP    NOT NULL DEFAULT current_timestamp(),
    CONSTRAINT sender_id_fk FOREIGN KEY (sender_id) REFERENCES users(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT receiver_id_fk FOREIGN KEY (receiver_id) REFERENCES users(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);