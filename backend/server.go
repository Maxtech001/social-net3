package main

import (
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/app/handlers"

	_ "social-network/pkg/db/sqlite"
)

const PORT = "5000"

func main() {
	manager := handlers.NewManager()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./pkg/static"))))
	http.HandleFunc("/ws", manager.HandleWS)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/getCurrentUser", handlers.CurrentUserHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/getUser", handlers.GetUserHandler)
	http.HandleFunc("/follow", handlers.FollowUserHandler)
	http.HandleFunc("/unfollow", handlers.UnfollowUserHandler)
	http.HandleFunc("/post", handlers.PostHandler)
	http.HandleFunc("/getAllPosts", handlers.GetAllPostsHandler)
	http.HandleFunc("/getAllGroupPosts", handlers.GetAllGroupPostsHandler)
	http.HandleFunc("/getUserPosts", handlers.GetUserPostsHandler)
	http.HandleFunc("/getComments", handlers.GetCommentsHandler)
	http.HandleFunc("/comment", handlers.CommentHandler)
	http.HandleFunc("/getMutualFollowers", handlers.GetMutualFollowersHandler)
	http.HandleFunc("/createGroup", handlers.CreateGroupHandler)
	http.HandleFunc("/getGroup", handlers.GetGroupHandler)
	http.HandleFunc("/createEvent", handlers.CreateEventHandler)
	http.HandleFunc("/getUpcomingEvents", handlers.GetUpcomingEventsHandler)
	http.HandleFunc("/getPastEvents", handlers.GetPastEventsHandler)
	http.HandleFunc("/getEvent", handlers.GetEventHandler)
	http.HandleFunc("/setGoing", handlers.SetGoingHandler)
	http.HandleFunc("/setNotGoing", handlers.SetNotGoingHandler)
	http.HandleFunc("/getFollowers", handlers.GetFollowersHandler)
	http.HandleFunc("/getFollowed", handlers.GetFollowedHandler)
	http.HandleFunc("/getFollowRequests", handlers.GetFollowRequestsHandler)
	http.HandleFunc("/acceptRequest", handlers.AcceptRequestHandler)
	http.HandleFunc("/discardRequest", handlers.DiscardRequestHandler)
	http.HandleFunc("/getUserGroups", handlers.GetUserGroupsHandler)
	http.HandleFunc("/getGroupsByName", handlers.GetGroupsByNameHandler)
	http.HandleFunc("/getGroupPosts", handlers.GetGroupPostsHandler)
	http.HandleFunc("/joinGroup", handlers.JoinGroupHandler)
	http.HandleFunc("/leaveGroup", handlers.LeaveGroupHandler)
	http.HandleFunc("/cancelGroupRequest", handlers.CancelGroupRequestHandler)
	http.HandleFunc("/getGroupRequests", handlers.GetGroupRequestsHandler)
	http.HandleFunc("/acceptGroupRequest", handlers.AcceptGroupRequestHandler)
	http.HandleFunc("/discardGroupRequest", handlers.DiscardGroupRequestHandler)
	http.HandleFunc("/storeInvites", handlers.StoreInvitesHandler)
	http.HandleFunc("/getMutualFollowersNotInGroup", handlers.GetMutualFollowersNotInGroupHandler)
	http.HandleFunc("/getFollowersNotInGroup", handlers.GetFollowersNotInGroupHandler)
	http.HandleFunc("/getGroupInvites", handlers.GetgroupInvitesHandler)
	http.HandleFunc("/acceptGroupInvite", handlers.AcceptGroupInviteHandler)
	http.HandleFunc("/discardGroupInvite", handlers.DiscardGroupInviteHandler)
	http.HandleFunc("/setPublic", handlers.SetPublicHandler)
	http.HandleFunc("/setPrivate", handlers.SetPrivateHandler)
	http.HandleFunc("/getAllUserNotifications", handlers.GetAllUserNotificationsHandler)
	http.HandleFunc("/setNotificationIsRead", handlers.SetNotificationIsReadHandler)
	http.HandleFunc("/getUsersByName", handlers.GetUsersByNameHandler)
	http.HandleFunc("/getConversation", handlers.GetConversationHandler)
	http.HandleFunc("/getUserConversations", handlers.GetUserConversationsHandler)
	http.HandleFunc("/getContacts", handlers.GetContactsHandler)
	http.HandleFunc("/getUserFollowers", handlers.GetUserFollowersHandler)
	http.HandleFunc("/getUserFollowing", handlers.GetUserFollowingHandler)
	http.HandleFunc("/setMessageIsRead", handlers.SetMessageIsReadHandler)
	fmt.Println("Starting server at: http://localhost:" + PORT)
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Fatal(err)
	}
}
