package models

import (
	"database/sql"
)

type User struct {
	Id          int    `json:"userId"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	BirthDate   string `json:"birthDate"`
	AvatarPath  string `json:"avatar"`
	Nickname    string `json:"nickname"`
	AboutMe     string `json:"aboutMe"`
	Followers   int    `json:"followers"`
	IsPublic    bool   `json:"isPublic"`
	IsFollowed  bool   `json:"isFollowed"`
	IsRequested bool   `json:"isRequested"`
	NameSearch  string `json:"nameSearch"`
}

type Post struct {
	Id             int          `json:"postId"`
	AuthorId       int          `json:"authorId"`
	GroupId        int          `json:"groupId"`
	Group          Group        `json:"group"`
	AuthorName     string       `json:"authorName"`
	AuthorAvatar   string       `json:"authorAvatar"`
	Content        string       `json:"content"`
	TotalComments  int          `json:"totalComments"`
	CreatedAt      string       `json:"createdAt"`
	ImageArray     []*PostImage `json:"imageArray"`
	PostType       int          `json:"postType"` // 0 - private, 1 - public, 2 - specified users
	SpecifiedUsers []*User      `json:"specifiedUsers"`
}

type PostImage struct {
	PostId    int    `json:"postId"`
	ImagePath string `json:"imagePath"`
}

type PostUser struct {
	PostId int `json:"postId"`
	UserId int `json:"userId"`
}

type Comment struct {
	Id           int             `json:"commentId"`
	PostId       int             `json:"postId"`
	AuthorId     int             `json:"authorId"`
	AuthorName   string          `json:"authorName"`
	AuthorAvatar string          `json:"authorAvatar"`
	Content      string          `json:"content"`
	CreatedAt    string          `json:"createdAt"`
	ImageArray   []*CommentImage `json:"imageArray"`
}

type CommentImage struct {
	CommentId int    `json:"commentId"`
	ImagePath string `json:"imagePath"`
}

type Group struct {
	Id          int     `json:"groupId"`
	AdminId     int     `json:"adminId"`
	Admin       User    `json:"admin"`
	GroupName   string  `json:"groupName"`
	Description string  `json:"description"`
	MemberCount int     `json:"memberCount"`
	CreatedAt   string  `json:"createdAt"`
	Members     []*User `json:"members"`
	IsMember    bool    `json:"isMember"`
	IsRequested bool    `json:"isRequested"`
}

type Event struct {
	Id              int              `json:"eventId"`
	AuthorId        int              `json:"authorId"`
	Author          User             `json:"author"`
	Title           string           `json:"title"`
	Description     string           `json:"description"`
	GroupId         int              `json:"groupId"`
	Group           Group            `json:"group"`
	CreatedAt       string           `json:"createdAt"`
	StartDate       string           `json:"startDate"`
	EndDate         string           `json:"endDate"`
	InvitedMembers  []*InvitedMember `json:"invitedMembers"`
	GoingMembers    []*InvitedMember `json:"goingMembers"`
	NotGoingMembers []*InvitedMember `json:"notGoingMembers"`
	UserAnswer      int              `json:"userAnswer"`
}

type InvitedMember struct {
	EventId  int  `json:"eventId"`
	MemberId int  `json:"memberId"`
	Member   User `json:"member"`
	Answer   int  `json:"answer"`
}

type Session struct {
	Id        string
	AccountId int
	Expiry    string
}

type FollowConnection struct {
	FollowedId int `json:"followedId"`
	FollowerId int `json:"followerId"`
}

type GroupMemberConnection struct {
	GroupId  int `json:"groupId"`
	MemberId int `json:"memberId"`
}

type GroupRequest struct {
	GroupId     int   `json:"groupId"`
	Group       Group `json:"group"`
	RequesterId int   `json:"requesterId"`
	User        User  `json:"user"`
}

type GroupInvite struct {
	GroupId     int   `json:"groupId"`
	Group       Group `json:"group"`
	InvitedId   int   `json:"invitedId"`
	InvitedUser User  `json:"invitedUser"`
	InviterId   int   `json:"inviterId"`
	InviterUser User  `json:"inviterUser"`
}

type Notification struct {
	Id               int    `json:"notificationId"`
	UserId           int    `json:"userId"`
	User             User   `json:"user"`
	Content          string `json:"content"`
	NotificationType int    `json:"notificationType"`
	IsRead           bool   `json:"isRead"`
	CreatedAt        string `json:"createdAt"`
	FollowerId       int    `json:"followerId"`
	NullFollowerId   sql.NullInt32
	Follower         *User `json:"follower"`
	GroupId          int   `json:"groupId"`
	NullGroupId      sql.NullInt32
	Group            *Group `json:"group"`
	EventId          int    `json:"eventId"`
	NullEventId      sql.NullInt32
	Event            *Event `json:"event"`
}

/*Conversation constructs*/
type Conversation struct {
	Id          int     `json:"conversationId"`
	User1ID     int     `json:"user1Id"`
	User1Name   string  `json:"user1Name"`
	User2ID     int     `json:"user2Id"`
	User2Name   string  `json:"user2Name"`
	CreatedAt   string  `json:"createdAt"`
	GroupId     int     `json:"groupId"`
	GroupObj    *Group  `json:"groupObj"` // receiver
	UserObj     *User   `json:"userObj"`  // receiver
	LastMessage Message `json:"lastMessage"`
}

/*Message constructs*/
type Message struct {
	Id             int    `json:"messageId"`
	ConversationId int    `json:"conversationId"`
	SenderId       int    `json:"senderId"`
	SenderName     string `json:"senderName"`
	Message        string `json:"message"`
	CreatedAt      string `json:"createdAt"`
	IsRead         bool   `json:"isRead"`
	SenderObj      User   `json:"senderObj"`
}

type Contact struct {
	User  *User
	Group *Group
}
