from telethon.sync import TelegramClient
from telethon.tl.functions.channels import GetChannelsRequest
from telethon.tl.functions.contacts import ResolveUsernameRequest
from telethon.tl.types import InputChannel
from telethon import events

api_id = 'your_api_id'  # 你的API ID
api_hash = 'your_api_hash'  # 你的API Hash
username = 'your_telegram_username'  # 你的Telegram用户名

client = TelegramClient(username, api_id, api_hash)

async def main():
    channel = await client.get_entity('channel_username')  # 你想要监听的频道的用户名
    print('Connected to channel: ', channel.title)

    @client.on(events.NewMessage(chats=channel))
    async def handler(event):
        print(event.message.text)  # 打印从频道接收的消息

    await client.run_until_disconnected()

with client:
    client.start()
    client.loop.run_until_complete(main())