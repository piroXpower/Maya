# Copyright (C) 2018 - 2020 MrYacha. All rights reserved. Source code available under the AGPL.
# Copyright (C) 2019 Aiogram

#
# This file is part of SophieBot.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.

# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
import html
import io
import math

from PIL import Image
from aiogram.types import InlineKeyboardMarkup, InlineKeyboardButton, ParseMode
from aiogram.types.input_file import InputFile
from aiogram.utils.exceptions import InvalidStickersSet, BadRequest, InvalidPeerID, BotBlocked

from maya_bot import bot, BOT_USERNAME
from maya_bot.decorator import register
from .utils.disable import disableable_dec
from .utils.language import get_strings_dec


@register(cmds='getsticker')
@disableable_dec('getsticker')
@get_strings_dec('stickers')
async def get_sticker(message, strings):
    if 'reply_to_message' not in message or 'sticker' not in message.reply_to_message:
        await message.reply(strings['rpl_to_sticker'])
        return

    sticker = message.reply_to_message.sticker
    file_id = sticker.file_id
    text = strings['ur_sticker'].format(emoji=sticker.emoji, id=file_id)

    sticker_file = await bot.download_file_by_id(file_id, io.BytesIO())

    await message.reply_document(
        InputFile(sticker_file, filename=f'{sticker.set_name}_{sticker.file_id[:5]}.png'),
        text
    )


@register(cmds='kang')
@disableable_dec('kang')
@get_strings_dec('stickers')
async def get_sticker(message, strings):
    user_id = message.from_user.id
    first_name = html.escape(message.from_user.first_name)

    if message.reply_to_message:
        if message.reply_to_message.sticker:
            file_id = message.reply_to_message.sticker.file_id
        elif message.reply_to_message.photo:
            file_id = message.reply_to_message.photo[-1].file_id
        elif message.reply_to_message.document:
            file_id = message.reply_to_message.document.file_id
        else:
            await message.reply(strings['rpl_to_sticker_image'])
            return
    else:
        await message.reply(strings['rpl_to_sticker_image'])
        return

    if message.get_args():
        sticker_emoji = message.get_args()[0]
    elif message.reply_to_message.sticker and message.reply_to_message.sticker.emoji:
        sticker_emoji = message.reply_to_message.sticker.emoji
    else:
        sticker_emoji = "🤔"

    kang_file = await bot.download_file_by_id(file_id, io.BytesIO())

    if not message.reply_to_message.sticker:
        try:
            im = Image.open(kang_file)
            maxsize = 512
            if (im.width and im.height) < maxsize:
                w = im.width
                h = im.height
                if im.width > im.height:
                    scale = maxsize / w
                    neww = maxsize
                    newh = h * scale
                else:
                    scale = maxsize / h
                    neww = w * scale
                    newh = maxsize
                sizenew = (math.floor(neww), math.floor(newh))
                im = im.resize(sizenew)
            else:
                im.thumbnail((maxsize, maxsize))
            sticker_file = io.BytesIO()
            im.save(sticker_file, "PNG")
            sticker_file.seek(0)
        except OSError:
            await message.reply(strings['rpl_to_sticker_image'])
            return
    else:
        sticker_file = kang_file

    packnum = 0
    packname = "a" + str(user_id) + "_by_" + BOT_USERNAME
    packname_found = False
    max_stickers = 120
    while not packname_found:
        try:
            stickerset = await bot.get_sticker_set(packname)
            if len(stickerset.stickers) >= max_stickers:
                packnum += 1
                packname = "a" + str(packnum) + "_" + str(user_id) + "_by_" + BOT_USERNAME
            else:
                packname_found = True
        except InvalidStickersSet:
            packname_found = True
    try:
        try:
            await bot.add_sticker_to_set(user_id=user_id, name=packname, png_sticker=sticker_file.getvalue(),
                                         emojis=sticker_emoji)
        except InvalidStickersSet:
            if packnum > 0:
                extra_version = " " + str(packnum)
            else:
                extra_version = ""
            await bot.create_new_sticker_set(user_id=user_id, name=packname, png_sticker=sticker_file.getvalue(),
                                             emojis=sticker_emoji, title=f"{first_name}s kang pack" + extra_version)
    except (InvalidPeerID, BotBlocked):
        await message.reply("Message bot first")
        buttons = InlineKeyboardMarkup()
        buttons.add(InlineKeyboardButton("Start", url=f"https://t.me/{BOT_USERNAME}"))
        await message.reply(strings["msg_first"], reply_markup=buttons)
        return
    except BadRequest as e:
        print(e)
        await message.reply(strings["invalid_emoji"])
        return

    text = strings["sticker_added"].format(packname=packname)
    await message.reply(text, parse_mode=ParseMode.MARKDOWN)
