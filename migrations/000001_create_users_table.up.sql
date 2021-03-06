CREATE TABLE IF NOT EXISTS users (
    id                  INTEGER         NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name                VARCHAR(255)    NOT NULL,
    email               VARCHAR(255)    NOT NULL,
    hashed_password     CHAR(60)        NOT NULL,
    created             DATETIME        NOT NULL,
    active              BOOLEAN         NOT NULL DEFAULT TRUE,
    image_link          VARCHAR(500)
);