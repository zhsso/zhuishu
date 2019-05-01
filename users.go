package main

import "strconv"

// UserManager 用户管理
type UserManager struct {
	users map[int64]*User
	umkey string
}

func newUserManager() *UserManager {
	return &UserManager{
		users: make(map[int64]*User),
		umkey: "userIDS",
	}
}

func (u *UserManager) getUser(id int64) *User {
	if user, ok := u.users[id]; ok {
		return user
	}
	return u.addUser(id)
}

func (u *UserManager) loadAll() {
	userIDStrs := redisClient.SMembers(u.umkey)
	if userIDStrs != nil {
		for _, userIDStr := range userIDStrs.Val() {
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err == nil {
				user := newUser(userID)
				user.loadBooks()
				u.users[userID] = user
			}
		}
	}
}

func (u *UserManager) addUser(userID int64) *User {
	user := newUser(userID)
	user.loadBooks()
	u.users[userID] = user
	redisClient.SAdd(u.umkey, userID)
	return user
}

func (u *UserManager) deleteUser(userID int64) {
	if user, ok := u.users[userID]; ok {
		user.delete()
	}
	delete(u.users, userID)
	redisClient.SRem(u.umkey, userID)
}
