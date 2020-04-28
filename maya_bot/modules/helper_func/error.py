# Copyright © 2018, 2019 MrYacha
# This file is part of SophieBot.
#
# SophieBot is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# SophieBot is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License

import sys

from maya_bot import dp


@dp.errors_handler()
async def all_errors_handler(message, dp):
    await report_error(message)


async def report_error(event, telethon=False):
    error = str(sys.exc_info()[1])
    class_error = sys.exc_info()[0].__name__

    if telethon is True:
        msg = event
    else:
        msg = event.message

    if class_error == 'ChatWriteForbiddenError':
        # This error mean bot is muted in chat
        return
    elif class_error == 'BadRequest' and error == 'Have no rights to send a message':
        return
    elif class_error == 'RetryAfter':
        return

    elif class_error == 'ServerSelectionTimeoutError':  # Redis server error
        return
    elif class_error == 'NotMasterError':  # MongoDB server error
        return

    text = "I'm sorry, I get a trouble. This incident was reported!"
    text += f'{class_error}: {error}'
    
    try:
        await msg.reply(text)
    except Exception:
        pass
    return
