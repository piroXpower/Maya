# MayaBot

MayaBot can be ran in 2 ways.

Using Docker is always recommended, as its setup is automated, and needs a very little knowledge of Linux or Command Line.
You also gain an advantage of isolating your server from the Maya Bot.

You need to configure this bot a bit before it can be used, don't worry, its easy!

# Requirements

+ Install git, python 3.7+ and docker(for docker method) from your package manager
+ You need to know how to clone this repo


# Docker Way

## Cloning this repo
    git clone https://github.com/ZerNico/Maya

## Setting config

+ Go to MayaBot/data
+ Rename bot_conf.json.example to bot_conf.json
+ Open in text editor
+ Set mongo_conn to "mongo-server"
+ Set redis_conn to "redis-server"
+ Set other configs as needed

## Creating bridge
    docker network create mayabot-net

## Running Redis and MongoDB
    docker run -d --rm --name redis-server --network mayabot-net redis:alpine
    docker run -d --rm --name mongo-server --network mayabot-net mongo:latest

## Start a MayaBot
    docker run -d -v /home/nico/MayaBot/data/:/opt/maya_bot/data --network mayabot-net maya 


# I am an old man, I like to go the manual way...


## Cloning this repo
    git clone https://github.com/ZerNico/MayaBot


## Setting config

+ Go to MayaBot/data
+ Rename bot_conf.json.example to bot_conf.json
+ Open in text editor
+ Set configs as needed

## Installing requirements
    cd MayaBot
    sudo pip3 install -r requirements.txt
   
    redis and mongodb from your package manager

## Running

    cd MayaBot
    python3 -m maya_bot
