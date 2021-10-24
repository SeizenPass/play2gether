CREATE TABLE IF NOT EXISTS games_ownership (
     id                  INTEGER         NOT NULL PRIMARY KEY AUTO_INCREMENT,
     game_id             INTEGER         NOT NULL,
     user_id             INTEGER         NOT NULL,
     CONSTRAINT game_id_fk FOREIGN KEY (game_id) REFERENCES games(id)
                                           ON DELETE CASCADE
                                           ON UPDATE CASCADE,
    CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users(id)
                                           ON DELETE CASCADE
                                           ON UPDATE CASCADE
);