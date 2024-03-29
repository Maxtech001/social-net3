-- +migrate Up
CREATE TABLE follow_requests (
    FOLLOWEDID INTEGER,
    FOLLOWERID INTEGER,
    PRIMARY KEY (FOLLOWEDID, FOLLOWERID),
    FOREIGN KEY(FOLLOWEDID) REFERENCES users(ID),
    FOREIGN KEY(FOLLOWERID) REFERENCES users(ID)
);