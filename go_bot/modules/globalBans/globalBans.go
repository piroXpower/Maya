/*
 *   Copyright 2019 Nico Franke  <nico.franke01@gmail.com>
 *
 *   Permission is hereby granted, free of charge, to any person obtaining a copy
 *   of this software and associated documentation files (the "Software"), to deal
 *   in the Software without restriction, including without limitation the rights
 *   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *   copies of the Software, and to permit persons to whom the Software is
 *   furnished to do so, subject to the following conditions:
 *
 *   The above copyright notice and this permission notice shall be included in all
 *   copies or substantial portions of the Software.
 *
 *   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *   SOFTWARE.
 */

package globalBans

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
	"github.com/ZerNico/Maya/go_bot"
	"github.com/ZerNico/Maya/go_bot/modules/sql"
	"github.com/ZerNico/Maya/go_bot/modules/utils/chat_status"
	"github.com/ZerNico/Maya/go_bot/modules/utils/error_handling"
	"github.com/ZerNico/Maya/go_bot/modules/utils/extraction"
	"github.com/ZerNico/Maya/go_bot/modules/utils/helpers"
	"log"
	"strconv"
	"strings"
)

var gbanErrors = map[string]struct{}{
	"Bad Request: Chat not found": {},
	"Bad Request: not enough rights to restrict/unrestrict chat member": {},
	"Bad Request: User_not_participant": {},
	"Bad Request: Peer_id_invalid": {},
	"Bad Request: Group chat was deactivated": {},
	"Bad Request: Need to be inviter of a user to kick it from a basic group": {},
	"Bad Request: Chat_admin_required": {},
	"Bad Request: Only the creator of a basic group can kick group administrators": {},
	"Bad Request: Channel_private": {},
	"Bad Request: Not in the chat": {},
}

var unGbanErrors = map[string]struct{}{
	"Bad Request: User is an administrator of the chat": {},
	"Bad Request: Chat not found": {},
	"Bad Request: Not enough rights to restrict/unrestrict chat member": {},
	"Bad Request: User_not_participant": {},
	"Bad Request: Method is available for supergroup and channel chats only": {},
	"Bad Request: Not in the chat": {},
	"Bad Request: Channel_private": {},
	"Bad Request: Chat_admin_required": {},
	"Bad Request: Peer_id_invalid": {},
}

func globalBan(bot ext.Bot, u *gotgbot.Update, args []string) error {
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	userId, reason := extraction.ExtractUserAndText(msg, args)

	if userId == 0 {
		_, err := msg.ReplyText("Try targeting a user next time bud.")
		return err
	}

	sudos := go_bot.BotConfig.SudoUsers
	sudos = append(sudos, strconv.Itoa(go_bot.BotConfig.OwnerId))

	if !helpers.Contains(sudos, strconv.Itoa(user.Id)) {
		_, err := msg.ReplyText("You don't have permissions to gban!")
		return err
	}

	if userId == go_bot.BotConfig.OwnerId {
		_, err := msg.ReplyText("I'm not gbanning my owner!")
		return err
	}

	for _, id := range go_bot.BotConfig.SudoUsers {
		sudoId, _ := strconv.Atoi(id)
		if userId == sudoId {
			_, err := msg.ReplyText("I'm not gbanning a sudo user!")
			return err
		}
	}

	gbannedUser := sql.GetGbanUser(strconv.Itoa(userId))

	if reason == "" {
		reason = "No reason."
	}

	go sql.GbanUser(strconv.Itoa(userId), reason)
	member, _ := bot.GetChat(userId)

	if gbannedUser == nil {
		_, err := msg.ReplyHTMLf("%v has been globally banned.", helpers.MentionHtml(member.Id, member.FirstName))
		error_handling.HandleErr(err)
		go func(bot ext.Bot, userId int, gbanErrors map[string]struct{}) {
			for _, chatIdStr := range sql.AllChats() {

				chatId, err := strconv.Atoi(chatIdStr)
				error_handling.HandleErr(err)

				chat, err := bot.GetChat(chatId)
				error_handling.HandleErr(err)
				if chat_status.IsBotAdmin(chat, nil) &&
					!chat_status.IsUserAdmin(chat, userId) &&
					sql.DoesChatGban(chatIdStr) {
					_, err = bot.KickChatMember(chatId, userId)
					if err != nil {
						if _, contains := gbanErrors[err.Error()]; !contains {
							error_handling.HandleErr(err)
						}
					}
				}
			}
		}(bot, userId, gbanErrors)

		go func (bot ext.Bot, user *ext.User, member *ext.Chat, reason string, sudos []string) {
			for _, sudo := range sudos {
				sudoId, _ := strconv.Atoi(sudo)
				_, _ = bot.SendMessageHTML(sudoId, fmt.Sprintf("<b>New GBan</b>"+
					"\n<b>Banner</b>: %v"+
					"\n<b>User</b>: %v"+
					"\n<b>User ID</b>: <code>%v</code>"+
					"\n<b>Reason</b>: %v", helpers.MentionHtml(user.Id, user.FirstName), helpers.MentionHtml(member.Id, member.FirstName),
					member.Id, reason))
			}
		}(bot, user, member, reason, sudos)

	} else {
		_, err := msg.ReplyHTMLf("Reason of %v's GBan has been changed to %v", helpers.MentionHtml(member.Id, member.FirstName), reason)
		return err
	}
	return nil
}

func unGlobalBan(bot ext.Bot, u *gotgbot.Update, args []string) error {
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	userId, _ := extraction.ExtractUserAndText(msg, args)

	if userId == 0 {
		_, err := msg.ReplyText("Try targeting a user next time bud.")
		return err
	}

	sudos := go_bot.BotConfig.SudoUsers
	sudos = append(sudos, strconv.Itoa(go_bot.BotConfig.OwnerId))

	if !helpers.Contains(sudos, strconv.Itoa(user.Id)) {
		_, err := msg.ReplyText("You don't have permissions to ungban!")
		return err
	}

	if userId == go_bot.BotConfig.OwnerId {
		_, err := msg.ReplyText("I'm not gbanning my owner!")
		return err
	}

	for _, id := range go_bot.BotConfig.SudoUsers {
		sudoId, _ := strconv.Atoi(id)
		if userId == sudoId {
			_, err := msg.ReplyText("I'm not gbanning a sudo user!")
			return err
		}
	}

	gbannedUser := sql.GetGbanUser(strconv.Itoa(userId))

	if gbannedUser == nil {
		_, err := msg.ReplyText("This user isn't gbanned")
		return err
	}

	go sql.UnGbanUser(strconv.Itoa(userId))

	member, _ := bot.GetChat(userId)

	go func(bot ext.Bot, user int, unGbanErrors map[string]struct{}) {
		for _, chatIdStr := range sql.AllChats() {

			chatId, err := strconv.Atoi(chatIdStr)
			error_handling.HandleErr(err)

			chat, err := bot.GetChat(chatId)
			error_handling.HandleErr(err)
			if chat_status.IsBotAdmin(chat, nil) &&
				!chat_status.IsUserAdmin(chat, u.EffectiveUser.Id) &&
				sql.DoesChatGban(chatIdStr) {
				_, err = bot.UnbanChatMember(chatId, userId)
				if err != nil {
					if _, contains := unGbanErrors[err.Error()]; !contains {
						error_handling.HandleErr(err)
					}
				}
			}
		}
	}(bot, userId, unGbanErrors)

	_, err := msg.ReplyHTMLf("%v has been globally unbanned.", helpers.MentionHtml(member.Id, member.FirstName))
	return err
}

func globalCheckBan(bot ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	chatId := strconv.Itoa(chat.Id)

	if chat_status.IsUserAdmin(chat, u.EffectiveUser.Id) {
		return gotgbot.ContinueGroups{}
	}

	if !chat_status.IsBotAdmin(chat, nil) {
		return gotgbot.ContinueGroups{}
	}

	if !sql.DoesChatGban(chatId) {
		return gotgbot.ContinueGroups{}
	}

	member := sql.GetGbanUser(strconv.Itoa(user.Id))

	if member != nil {
		_, err := msg.Delete()
		error_handling.HandleErr(err)
		_, err = bot.KickChatMember(chat.Id, user.Id)
		return err
	}
	return gotgbot.ContinueGroups{}
}

func gbanStat(bot ext.Bot, u *gotgbot.Update, args []string) error {

	chat := u.EffectiveChat

	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chat_status.IsUserAdmin(chat, u.EffectiveUser.Id) {
		_, err := u.EffectiveMessage.ReplyText("You must be an admin to perform this action!")
		return err
	}

	if len(args) == 0 {
		gbanPrefs := sql.DoesChatGban(strconv.Itoa(chat.Id))
		_, err := u.EffectiveMessage.ReplyHTMLf("I am currently gbanning people in this group: <code>%v</code>",
			gbanPrefs)
		return err

	} else if len(args) >= 1 {
		switch strings.ToLower(args[0]){
		case "on", "yes":
			go sql.SetGBanPref(strconv.Itoa(chat.Id), true)
			_, err := u.EffectiveMessage.ReplyText("I'll gban users in this chat now.")
			return err
		case "off", "no":
			go sql.SetGBanPref(strconv.Itoa(chat.Id), false)
			_, err := u.EffectiveMessage.ReplyText("I'll not gban users in this chat anymore.")
			return err
		default:
			_, err := u.EffectiveMessage.ReplyText("I understand 'on/yes' or 'off/no' only!")
			return err
		}
	}
	return nil
}

func LoadGlobalBans(u *gotgbot.Updater) {
	defer log.Println("Loading module global_bans")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("gban", []rune{'/', '!'}, globalBan))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("ungban", []rune{'/', '!'}, unGlobalBan))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("gbanstat", []rune{'/', '!'}, gbanStat))

	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, globalCheckBan))
}
