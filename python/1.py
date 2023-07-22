
from telethon import TelegramClient, events

api_id = 26875732
api_hash = '3e9b6bccde2480cf6823ae785710f462'
bot_token = '6344282196:AAFUDFLir2uH4z16umoEURJzLvFUXVbrIKI'


bot = TelegramClient('bot', api_id, api_hash,proxy=("socks5", '127.0.0.1', 15235)).start(bot_token=bot_token)

@bot.on(events.NewMessage(pattern='/start'))
async def start(event):
    await event.reply('Hi, I am your bot!')

@bot.on(events.NewMessage(pattern='/hello'))
async def hello(event):
    await event.reply('Hello!')

print('Bot started!')
bot.run_until_disconnected()