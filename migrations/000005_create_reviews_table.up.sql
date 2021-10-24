CREATE TABLE IF NOT EXISTS reviews (
   id                       INTEGER         NOT NULL PRIMARY KEY AUTO_INCREMENT,
   review_text              VARCHAR(1500)   NOT NULL,
   reviewer_id             INTEGER,
   reviewed_id             INTEGER         NOT NULL,
   CONSTRAINT reviewed_id_fk FOREIGN KEY (reviewed_id) REFERENCES games(id)
       ON DELETE CASCADE
       ON UPDATE CASCADE,
   CONSTRAINT reviewer_id_fk FOREIGN KEY (reviewer_id) REFERENCES users(id)
       ON DELETE SET NULL
       ON UPDATE CASCADE
);