A Discord bot to fetch Cookie Run Braverse TCG card information. Database used is a PostGreSQL database from Supabase

Be sure to set the following environment variables:
    Powershell users:
    $env:BOT_TOKEN = <discord bot token>
    $env:CONNECTION_STR = <postgres database connection string>

To start the bot:

    go run main.go

Bot commands:

| Command  | Description |
| ------------- | ------------- |
| !help  | Displays all bot commands and their descriptions.  |
| !fetchEN [card name in English] | Searches for a card by its English name  |
| !fetchKR [card name in Korean] | Searches for a card by its Korean name |
