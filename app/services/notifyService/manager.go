package notifyService

import "fmt"

// represent collection map[user]map[room]countOfUnreadMessages
var userUnreadMessages map[string]map[string]int

func init() {
	userUnreadMessages = make(map[string]map[string]int)
}

func GetTotalUnreadCount(userID string) int {
	rooms := userUnreadMessages[userID]
	if len(rooms) == 0 {
		return 0
	}

	var count int
	count = 0
	for _, cnt := range rooms {
		count += cnt
	}

	return count
}

func GetUnreadByRooms(userID string) map[string]int {
	return userUnreadMessages[userID]
}

func Increase(userID, roomID string) {
	if rooms, ok := userUnreadMessages[userID]; ok {
		if cnt, rok := rooms[roomID]; rok {
			cnt++
			rooms[roomID] = cnt
		} else {
			rooms[roomID] = 1
		}
	} else {
		rooms := make(map[string]int)
		rooms[roomID] = 1
		userUnreadMessages[userID] = rooms
	}
	fmt.Printf("INCREASE COUNT TO MESAGE :v", userUnreadMessages[userID])
}

func Decrease(userID, roomID string) {
	if rooms, ok := userUnreadMessages[userID]; ok {
		if cnt, rok := rooms[roomID]; rok {
			cnt--
			if cnt >= 0 {
				delete(rooms, roomID)
			} else {
				rooms[roomID] = cnt
			}
		}
	}
}

func SetRoomReaded(userID, roomID string) {
	if rooms, ok := userUnreadMessages[userID]; ok {
		delete(rooms, roomID)
	}
}

func NotifyUsersAboutNewMessages(users []string) {

}

func notifyUser(userID string) {
	// userUnreadMessages[userID]
}
