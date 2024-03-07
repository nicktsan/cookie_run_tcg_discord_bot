A Discord bot to fetch Cookie Run Braverse TCG card information. Database used is a PostGreSQL database from Supabase

## Prerequsistes:
- Go installed
- Supabase account, project, and database table
- Discord Developer mode enabled

## How to make your own Discord bot:
1. Turn on “Developer mode” in your Discord account.
2. Click on “Discord API”.
3. In the Developer portal, click on “Applications”. Log in again and then, back in the “Applications” menu, click on “New  Application”.
4. Name the bot and then click “Create”.
5. Go to the “Bot” menu and generate a token using “Add Bot”.
6. Enable "Message Content Intent".
7. Program your bot using the bot token and save the file.
8. Define other details for your bot under “General Information”.
9. Click on “OAuth2”, activate “bot”, set the permissions, and then Copy the generated url
10. Paste the generated url into your web browser.
11. Select your server to add your bot to it.


## Be sure to set the following environment variables:

    $env:BOT_TOKEN = <discord bot token>
    $env:CONNECTION_STR = <postgres database connection string>

## To start the bot:

    go run main.go

## Bot commands:

| Command  | Description |
| ------------- | ------------- |
| !help  | Displays all bot commands and their descriptions.  |
| !fetchEN [card name in English] | Searches for a card by its English name.  |
| !fetchKR [card name in Korean] | Searches for a card by its Korean name. |
