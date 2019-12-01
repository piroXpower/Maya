/*
 *   Copyright 2019 ATechnoHazard  <amolele@gmail.com>
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

package help

import (
	"fmt"
	"github.com/ZerNico/Maya/go_bot/modules/utils/error_handling"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"html"
	"log"
	"regexp"
)

var markup ext.InlineKeyboardMarkup
var markdownHelpText string

func initMarkdownHelp() {
	markdownHelpText = "You can use markdown to make your messages more expressive. This is the markdown currently " +
		"supported:\n\n" +
		"<code>`code words`</code>: backticks allow you to wrap your words in monospace fonts.\n" +
		"<code>*bold*</code>: wrapping text with '*' will produce bold text\n" +
		"<code>_italics_</code>: wrapping text with '_' will produce italic text\n" +
		"<code>[hyperlink](example.com)</code>: this will create a link - the message will just show " +
		"<code>hyperlink</code>, and tapping on it will open the page at <code>example.com</code>\n\n" +
		"<code>[buttontext](buttonurl:example.com)</code>: this is a special enhancement to allow users to have " +
		"telegram buttons in their markdown. <code>buttontext</code> will be what is displayed on the button, and " +
		"<code>example.com</code> will be the url which is opened.\n\n" +
		"If you want multiple buttons on the same line, use :same, as such:\n" +
		"<code>[one](buttonurl://github.com)</code>\n" +
		"<code>[two](buttonurl://google.com:same)</code>\n" +
		"This will create two buttons on a single line, instead of one button per line.\n\n" +
		"Keep in mind that your message MUST contain some text other than just a button!"

}

func initHelpButtons() {
	helpButtons := [][]ext.InlineKeyboardButton{make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 2),
		make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 2)}

	// First column
	helpButtons[0][0] = ext.InlineKeyboardButton{
		Text:         "Admin",
		CallbackData: fmt.Sprintf("help(%v)", "admin"),
	}
	helpButtons[1][0] = ext.InlineKeyboardButton{
		Text:         "Bans",
		CallbackData: fmt.Sprintf("help(%v)", "bans"),
	}
	helpButtons[2][0] = ext.InlineKeyboardButton{
		Text:         "Blacklist",
		CallbackData: fmt.Sprintf("help(%v)", "blacklist"),
	}
	helpButtons[3][0] = ext.InlineKeyboardButton{
		Text:         "Deleting",
		CallbackData: fmt.Sprintf("help(%v)", "deleting"),
	}
	helpButtons[4][0] = ext.InlineKeyboardButton{
		Text:         "Federations",
		CallbackData: fmt.Sprintf("help(%v)", "feds"),
	}

	// Second column
	helpButtons[0][1] = ext.InlineKeyboardButton{
		Text:         "Misc",
		CallbackData: fmt.Sprintf("help(%v)", "misc"),
	}
	helpButtons[1][1] = ext.InlineKeyboardButton{
		Text:         "Muting",
		CallbackData: fmt.Sprintf("help(%v)", "muting"),
	}
	helpButtons[2][1] = ext.InlineKeyboardButton{
		Text:         "Notes",
		CallbackData: fmt.Sprintf("help(%v)", "notes"),
	}
	helpButtons[3][1] = ext.InlineKeyboardButton{
		Text:         "Global bans",
		CallbackData: fmt.Sprintf("help(%v)", "globalbans"),
	}
	helpButtons[4][1] = ext.InlineKeyboardButton{
		Text:         "Warns",
		CallbackData: fmt.Sprintf("help(%v)", "warns"),
	}

	markup = ext.InlineKeyboardMarkup{InlineKeyboard: &helpButtons}
}

func help(b ext.Bot, u *gotgbot.Update) error {
	msg := b.NewSendableMessage(u.EffectiveChat.Id, "Hey there! I'm Maya, a group management bot written in Go,"+
		"here to help you manage your groups!\n" +
		"I have a ton of useful features, such as a note keeping system, administration, filters and even a warn system.\n\n"+
		"Commands are preceded with a slash (/) or an exclamation mark (!)\n\n"+
		"Some basic commands:\n"+
		" - /start: duh, you already know what this does\n"+
		" - /help: for info on how to use me\n"+
		" - /donate: info on who made me and how you can support them\n\n"+
		"If you have any bugs reports, questions or suggestions you can message me (@NicoFranke).")
	msg.ParseMode = parsemode.Html
	msg.ReplyToMessageId = u.EffectiveMessage.MessageId
	msg.ReplyMarkup = &markup
	_, err := msg.Send()
	if err != nil {
		msg.ReplyToMessageId = 0
		_, err = msg.Send()
	}
	return err
}

func markdownHelp(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	if chat.Type != "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in PM!")
		return err
	}

	_, err := u.EffectiveMessage.ReplyHTML(markdownHelpText)
	return err
}

func buttonHandler(b ext.Bot, u *gotgbot.Update) error {
	query := u.CallbackQuery
	pattern, _ := regexp.Compile(`help\((.+?)\)`)

	if pattern.MatchString(query.Data) {
		module := pattern.FindStringSubmatch(query.Data)[1]
		chat := u.EffectiveChat
		msg := b.NewSendableEditMessageText(chat.Id, u.EffectiveMessage.MessageId, "placeholder")
		msg.ParseMode = parsemode.Html
		backButton := [][]ext.InlineKeyboardButton{{ext.InlineKeyboardButton{
			Text:         "Back",
			CallbackData: "help(back)",
		}}}
		backKeyboard := ext.InlineKeyboardMarkup{InlineKeyboard: &backButton}
		msg.ReplyMarkup = &backKeyboard

		switch module {
		case "admin":
			msg.Text = "Here is the help for the <b>Admin</b> module:\n\n" +
				" - /adminlist: Lists all admins in the chat.\n\n" +
				"<b>Admin only:</b>" +
				html.EscapeString("- /pin: Silently pins the message replied to - add 'loud' or 'notify' to give notifs to users.\n"+
					" - /unpin: Unpins the currently pinned message.\n"+
					" - /invitelink: Gets the groups invitelink.\n"+
					" - /promote: Promotes the user replied to.\n"+
					" - /demote: Demotes the user replied to.\n")
			break
		case "bans":
			msg.Text = "Here is the help for the <b>Bans</b> module:\n\n" +
				" - /kickme: Kicks the user who issued the command.\n\n" +
				"<b>Admin only</b>:\n" +
				html.EscapeString(" - /ban <userhandle>: Bans a user (via handle, or reply).\n"+
					" - /tban <userhandle> x(m/h/d): Bans a user for x time (via handle, or reply). m = minutes, h = hours,"+
					" d = days.\n"+
					" - /unban <userhandle>: Unbans a user (via handle, or reply)."+
					" - /kick <userhandle>: Kicks a user (via handle, or reply).")

			break
		case "blacklist":
			msg.Text = "Here is the help for the <b>Word Blacklists</b> module:\n\n" +
				"Blacklists are used to stop certain triggers from being said in a group. Any time the trigger is " +
				"mentioned, the message will immediately be deleted. A good combo is sometimes to pair this up with " +
				"warn filters!\n\n" +
				"<b>NOTE:</b> Blacklists do not affect group admins.\n\n" +
				" - /blacklist: View the current blacklisted words.\n\n" +
				"<b>Admin only:</b>\n" +
				html.EscapeString(" - /addblacklist <triggers>: Adds a trigger to the blacklist. Each line is "+
					"considered one trigger, so using different lines will allow you to add multiple triggers.\n"+
					" - /unblacklist <triggers>: Removes triggers from the blacklist. Same newline logic applies here, "+
					"so you can remove multiple triggers at once.\n"+
					" - /rmblacklist <triggers>: Same as /unblacklist.")
			break
		case "deleting":
			msg.Text = "Here is the help for the <b>Purges</b> module:\n\n" +
				"<b>Admin only:</b>\n" +
				" - /del: Deletes the message you replied to.\n" +
				" - /purge: Deletes all messages between this and the replied to message.\n"
			break
		case "feds":
			msg.Text = "Here is the help for the <b>Federations</b> module:\n\n" +
				html.EscapeString("Ah, group management. It's all fun and games, until you start getting spammers in, and you need to ban them. " + "" +
					"Then you need to start banning more, and more, and it gets painful. " +
					"But then you have multiple groups, and you don't want these spammers in any of your groups - how can you deal? " +
					"Do you have to ban them manually, in all your groups?\n\n" +
					"No more! With federations, you can make a ban in one chat overlap to all your other chats." +
					"You can even appoint federation admins, so that your trustworthiest admins can ban across all the chats that you want to protect.\n\n") +

				html.EscapeString(" - /newfed <fedname>: Creates a new federation with the given name. Users are only allowed to own one federation. " +
					"Using this method when you already have a fed will simply change the federation name.\n" +
					" - /delfed: Deletes your federation, and any information relating to it. Will not unban any banned users.\n" +
					" - /fedinfo <fedID>: Shows information about the specified federation.\n" +
					" - /joinfed <fedID>: Joins the current chat to the federation. Each chat can only be in one federation. " +
					"Only chat owners can do this.\n" +
					" - /leavefed: Leaves the current federation. Only chat owners can do this.\n" +
					" - /fedadmins <fedID>: Lists the admins in a federation.\n" +
					" - /fedstat: Lists all the federations you've been banned from.\n" +
					" - /fedstat <userhandle>: Lists all the federations the specified user has been banned from (also works with username, mention, and replies).\n" +
					" - /fedstat <userhandle> <fedID>: Gives information on the specified user's ban reason in that federation. " +
					"If no user is specified, checks the sender.\n" +
					" - /chatfed: Gets information about the federation that the current chat is in.\n\n") +
				"<b>Fed admin only:</b>\n" +
				html.EscapeString(" - /fban <user>: Bans a user from the current chat's federation.\n" +
					" - /unfban <user>: Unbans a user the current chat's federation.\n\n") +
				"<b>Fed owner only:</b>\n" +
				html.EscapeString(" - /fedpromote <userhandle>: Promotes the user to fed admin in your fed.\n" +
					" - /feddemote <userhandle>: Demotes the user from fed admin to normal user, in your fed.")
			break
		case "misc":
			msg.Text = "Here is the help for the <b>Misc</b> module:\n\n" +
				html.EscapeString(" - /id <userhandle>: Gets the ID of a user or group.\n" +
					" - /info <userhandle>: Displays info about a user.\n" +
					" - /ping: Shows the ping of the bot.")
			break
		case "muting":
			msg.Text = "Here is the help for the <b>Muting</b> module:\n\n" +
				"<b>Admin only:</b>\n" +
				html.EscapeString(" - /mute <userhandle>: Silences a user. Can also be used as a reply, muting the "+
					"replied to user.\n"+
					" - /tmute <userhandle> x(m/h/d): Mutes a user for x time (via handle, or reply). m = minutes, h = "+
					"hours, d = days.\n"+
					" - /unmute <userhandle>: Unmutes a user. Can also be used as a reply, muting the replied to user.")
			break
		case "notes":
			msg.Text = "Here is the help for the <b>Notes</b> module:\n\n" +
				html.EscapeString(" - /get <notename>: Gets the note with this notename.\n"+
					" - #<notename>: Same as /get.\n"+
					" - /notes or /saved: Lists all saved notes in this chat.\n\n"+
					"If you would like to retrieve the contents of a note without any formatting, use /get"+
					" <notename> noformat. This can be useful when updating a current note.\n\n") +
				"<b>Admin only:</b>\n" +
				html.EscapeString(" - /save <notename> <notedata>: Saves notedata as a note with name notename.\n"+
					"A button can be added to a note by using standard markdown link syntax - the link should just "+
					"be prepended with a buttonurl: section, as such: [somelink](buttonurl:example.com). Check "+
					"/markdownhelp for more info.\n"+
					" - /save <notename>: Saves the replied-to message as a note with name notename.\n"+
					" - /clear <notename>: Clears note with this name.")
			break
		case "users":
			break
		case "warns":
			msg.Text = "Here is the help for the <b>Warnings</b> module:\n\n" +
				html.EscapeString(" - /warns <userhandle>: Gets a user's number, and reason, of warnings.\n"+
					" - /warnlist: Gets a list of all current warning filters.\n\n") +
				"<b>Admin only:</b>\n" +
				html.EscapeString(" - /warn <userhandle>: Warns a user. After the warn limit, the user will be banned from the group. "+
					"Can also be used as a reply.\n"+
					" - /resetwarn <userhandle>: Resets the warnings for a user. Can also be used as a reply.\n"+
					" - /addwarn <keyword> <reply message>: Sets a warning filter on a certain keyword. If you want your "+
					"keyword to be a sentence, encompass it with quotes, as such: /addwarn \"very angry\" "+
					"This is an angry user.\n"+
					" - /nowarn <keyword>: Stops a warning filter\n"+
					" - /warnlimit <num>: Sets the warning limit\n"+
					" - /strongwarn <on/yes/off/no>: If set to on, exceeding the warn limit will result in a ban. "+
					"Else, will just kick.\n")
			break
		case "globalbans":
			msg.Text = "Here is the help for the <b>Global Bans</b> module:\n\n" +
				"<b>Admin only:</b>\n" +
				html.EscapeString(" - /gbanstat <on/off/yes/no>: Will disable the effect of global bans on your group, or return your current settings.\n\n") +
				"<b>Sudo only:</b>\n" +
				html.EscapeString(" - /gban <userhandle>: globally bans a user. (via handle, or reply) \n" +
					" - /unban <userhandle>: unbans a globally banned user. (via handle, or reply) \n\n" +
					"Gbans, also known as global bans, are used by the bot owners to ban spammers across all groups. " +
					"This helps protect you and your groups by removing spam flooders as quickly as possible. " +
					"They can be disabled for you group by calling /gbanstat")
			break
		case "back":
			msg.Text = "Hey there! I'm Maya, a group management bot written in Go,"+
			"here to help you manage your groups!\n" +
				"I have a ton of useful features, such as a note keeping system, administration, filters and even a warn system.\n\n"+
				"Commands are preceded with a slash (/) or an exclamation mark (!)\n\n"+
				"Some basic commands:\n"+
				" - /start: duh, you already know what this does\n"+
				" - /help: for info on how to use me\n"+
				" - /donate: info on who made me and how you can support them\n\n"+
				"If you have any bugs reports, questions or suggestions you can message me (@NicoFranke)."
			msg.ReplyMarkup = &markup
			break
		}

		_, err := msg.Send()
		error_handling.HandleErr(err)
		_, err = b.AnswerCallbackQuery(query.Id)
		return err
	}
	return nil
}

func LoadHelp(u *gotgbot.Updater) {
	defer log.Println("Loading module help")
	initHelpButtons()
	initMarkdownHelp()
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("help", []rune{'/', '!'}, help))
	u.Dispatcher.AddHandler(handlers.NewCallback("help", buttonHandler))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("markdownhelp", []rune{'/', '!'}, markdownHelp))
}
