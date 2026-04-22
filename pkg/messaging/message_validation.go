package messaging

var allowedMessageTypes = map[ServiceType]map[MessageType]bool{
	UserService: {
		AuthTypes.CreatedUser:  true,
		SocialTypes.BlockUser:  true,
		SocialTypes.FollowUser: true,
		"password_reset":       true,
	},
	AuthService: {
		"session_expired": true,
	},
}

func isAllowedMessageType(service ServiceType, messageType MessageType) bool {
	allowedTypes, ok := allowedMessageTypes[service]
	// Eğer servis listede yoksa tüm mesajlara izin ver (mevcut mantığın)
	if !ok {
		return true
	}

	// Doğrudan map üzerinden kontrol (Döngüye gerek kalmaz)
	return allowedTypes[messageType]
}
