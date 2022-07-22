# Outline Bot
Telegram bot for your Outline server

## Client
<row>
    <img src="docs/client_start.jpg" alt="client start" width="200"/>
    <img src="docs/client_pay.jpg" alt="client pay" width="200"/>
    <img src="docs/client_help.jpg" alt="client help" width="200"/>
</row>

## Admin
<row>
    <img src="docs/admin_help.jpg" alt="client start" width="200"/>
    <img src="docs/admin_after_help.jpg" alt="client pay" width="200"/>
</row>

## Server:
Export `Telegram bot token`, `Outline management API URL`, `card number`, `tonnames url` to environment
```
export TGTOKEN=<tbot token>
export OUTLINEAPI=<management API URL>
export CARD=<card number>
export TON=<tonnames url>
```
Run server
```
>> outlinebot run
```
Give admin rights by Telegram user ID
```
>> outlinebot admin -u <user_id>
```
Deprive admin rights by Telegram user ID
```
>> outlinebot admin -u <user_id> -d
```