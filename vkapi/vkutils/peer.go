package vkutils

func UserToPeer(userID int) int {
	return userID
}

func GroupToPeer(groupID int) int {
	return -groupID
}

func ChatToPeer(chatID int) int {
	return 2e9 + chatID
}

func PeerToUser(peerID int) int {
	return peerID
}

func PeerToGroup(peerID int) int {
	return -peerID
}

func PeerToChat(peerID int) int {
	return peerID - 2e9
}
