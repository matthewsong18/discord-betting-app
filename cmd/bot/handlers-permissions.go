package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func doesNotHaveManageMemberPerm(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	if (i.Member.Permissions & discordgo.PermissionManageMessages) != discordgo.PermissionManageMessages {
		log.Printf("User \"%s\" does not have permission to end polls", i.Member.User.GlobalName)
		sendInteractionResponse(s, i, "You do not have permission to edit polls")
		return true
	}
	return false
}
