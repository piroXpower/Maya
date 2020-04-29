# Copyright © 2020 Nico Franke
# This file is part of Maya.
#
# Maya is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Maya is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
import io
import math
import os
import random
import string

from PIL import Image
from aiogram.types import ParseMode, InlineKeyboardMarkup, InlineKeyboardButton
from aiogram.utils.exceptions import BadRequest, InvalidStickersSet, InvalidPeerID, BotBlocked

from maya_bot import BOT_USERNAME, decorator, bot
from maya_bot.modules.disable import disablable_dec
from maya_bot.modules.language import get_strings_dec
from maya_bot.modules.users import (aio_get_user)


@decorator.command("kang")
@disablable_dec("kang")
@get_strings_dec('stickers')
async def get_id(message, strings):
    user, txt = await aio_get_user(message, allow_self=True)

    packnum = 0
    packname = "a" + str(user['user_id']) + "_by_" + BOT_USERNAME
    packname_found = 0
    max_stickers = 120
    while packname_found == 0:
        try:
            stickerset = await bot.get_sticker_set(packname)
            if len(stickerset.stickers) >= max_stickers:
                packnum += 1
                packname = "a" + str(packnum) + "_" + str(user.id) + "_by_" + BOT_USERNAME
            else:
                packname_found = 1
        except InvalidStickersSet:
            packname_found = 1

    kangsticker = random_string(8) + ".png"
    if message.reply_to_message:
        if message.reply_to_message.sticker:
            file_id = message.reply_to_message.sticker.file_id
        elif message.reply_to_message.photo:
            file_id = message.reply_to_message.photo[-1].file_id
        elif message.reply_to_message.document:
            file_id = message.reply_to_message.document.file_id
        else:
            await message.reply(strings['cant_kang'])
            return
        kang_file = await bot.get_file(file_id)
        await kang_file.download("temp/" + kangsticker)
        if message.get_args():
            sticker_emoji = message.get_args()[0]
        elif message.reply_to_message.sticker and message.reply_to_message.sticker.emoji:
            sticker_emoji = message.reply_to_message.sticker.emoji
        else:
            sticker_emoji = "🤔"
        try:
            im = Image.open("temp/" + kangsticker)
            maxsize = (512, 512)

            if (im.width and im.height) < 512:
                size1 = im.width
                size2 = im.height
                if im.width > im.height:
                    scale = 512 / size1
                    size1new = 512
                    size2new = size2 * scale
                else:
                    scale = 512 / size2
                    size1new = size1 * scale
                    size2new = 512
                size1new = math.floor(size1new)
                size2new = math.floor(size2new)
                sizenew = (size1new, size2new)
                im = im.resize(sizenew)
            else:
                im.thumbnail(maxsize)
            if not message.reply_to_message.sticker:
                im.save("temp/" + kangsticker, "PNG")

            png_sticker = io.BytesIO()
            sticker = Image.open("temp/" + kangsticker, mode='r')
            sticker.save(png_sticker, format="PNG")
            png_sticker = png_sticker.getvalue()

            await bot.add_sticker_to_set(user_id=user["user_id"], name=packname,
                                         png_sticker=png_sticker, emojis=sticker_emoji)
            text = strings["sticker_added"].format(packname=packname)
            os.remove("temp/" + kangsticker)
            await message.reply(text, parse_mode=ParseMode.MARKDOWN)
            return
        except OSError:
            await message.reply(strings["only_images"])
            os.remove("temp/" + kangsticker)
            return
        except InvalidStickersSet:
            await makepack_internal(message, user, png_sticker, sticker_emoji, packname, packnum, strings)
            os.remove("temp/" + kangsticker)
        except BadRequest as e:
            if str(e) == "Invalid sticker emojis":
                await message.reply(strings["invalid_emoji"])
            os.remove("temp/" + kangsticker)
    else:
        await message.reply(strings["wrong_reply"])
        return


async def makepack_internal(message, user, png_sticker, emoji, packname, packnum, strings):
    name = user["first_name"][:50]
    try:
        extra_version = ""
        if packnum > 0:
            extra_version = " " + str(packnum)
        print(packname)
        success = await bot.create_new_sticker_set(user_id=user["user_id"], name=packname,
                                                   title=f"{name}s kang pack" + extra_version,
                                                   png_sticker=png_sticker, emojis=emoji)
    except (InvalidPeerID, BotBlocked):
        buttons = InlineKeyboardMarkup()
        buttons.add(InlineKeyboardButton("Start", url=f"https://t.me/{BOT_USERNAME}"))
        await message.reply(strings["contact_first"], reply_markup=buttons)
        return
    if success:
        text = strings["pack_created"].format(packname=packname)
        await message.reply(text, parse_mode=ParseMode.MARKDOWN)
    else:
        await message.reply(strings["failed_create"])


def random_string(stringLength):
    letters = string.ascii_letters
    return ''.join(random.choice(letters) for i in range(stringLength))
