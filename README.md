# Maya
[![Dockerhub Status](https://img.shields.io/docker/cloud/build/zernico/maya)](https://hub.docker.com/r/zernico/maya)

Telegram bot written in Go. Currently in alpha. Contributions are welcome.

A modular group management bot, written with the purpose of being highly concurrent.

You can [find me on telegram](https://t.me/NicoFranke)! I'm usually online, so I can hopefully answer any questions you may have.



## Setting up the bot (Important! Please go through once):

### Configuration
The preferred method is to create a dotenv file named `.env`, as it makes it much easier to see all your configuration settings grouped together. A sample dotenv file called `sample.env` has been included for convenience.

The available fields for the .env file are as follows:
* `BOT_API_KEY` :  Your bot token, as a string
* `BOT_NAME` : The name of your bot, as it appears on telegram
* `OWNER_USERNAME` : Your Telegram username, without the `@`
* `OWNER_ID` : Your Telegram ID
* `DATABASE_URI`: Self explanatory (postgres)
* `SUDO_USERS`: A list of userIDs, separated by spaces, who should have sudo access to the bot
* `HEROKU`: Setting this to **anything** will activate it. Use if you're using a heroku database
* `ENV`: Setting this will allow you to run the bot without a .env file
* `DEBUG`: Setting this to **anything** will activate it. Use it if you're debugging something.


## Starting the bot
Download the latest binary for your machine's OS and architecture from the releases page. Put it in the same directory as the .env file, and execute it.

It's that simple.

## Credits
Thanks to [ATechnoHazard](https://github.com/ATechnoHazard) the original creator of [ginko](https://github.com/ATechnoHazard/ginko).  
Thanks to [PaulSonOfLars](https://github.com/PaulSonOfLars) for his [telegram go library](https://github.com/PaulSonOfLars/gotgbot).

## Download source
Contributions to this project are welcome.
To download the source, get it like any other Go project:
 `go get -u github.com/ZerNico/Maya`.
