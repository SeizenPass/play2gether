CREATE TABLE IF NOT EXISTS games (
    id                  INTEGER         NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title               VARCHAR(255)    NOT NULL,
    description         VARCHAR(2000),
    image_link          VARCHAR(500)
);