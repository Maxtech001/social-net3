-- +migrate Up
CREATE TABLE post_user_connections (
    POSTID INTEGER NOT NULL,
    USERID INTEGER NOT NULL,
    FOREIGN KEY(POSTID) REFERENCES posts(ID),
    FOREIGN KEY(USERID) REFERENCES users(ID)
);