CREATE TABLE IF NOT EXISTS group_members (
    groupId INTEGER NOT NULL,
    memberId INTEGER NOT NULL,
    FOREIGN KEY (groupId) REFERENCES `groups` (groupId) ON DELETE CASCADE,
    FOREIGN KEY (memberId) REFERENCES users (userId) ON DELETE CASCADE
)