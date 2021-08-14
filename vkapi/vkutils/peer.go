package vkutils

func ToPeer(id int) int {
	if id < 0 {
		id = -id
	} else if id > 2e9 {
		id -= 2e9
	}
	return id
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
