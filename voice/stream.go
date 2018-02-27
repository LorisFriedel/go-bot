package voice

import dg "github.com/bwmarrin/discordgo"

func test() {
	var v dg.VoiceConnection
	v.AddHandler(func(vc *dg.VoiceConnection, vs *dg.VoiceSpeakingUpdate) {

	})
}
