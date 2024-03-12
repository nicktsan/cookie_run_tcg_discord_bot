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

    $ go run main.go

## Building the docker image
Navigate to directory with the Dockerfile. Then, containerize the app with the command: 

    $ docker compose build

This should build a local Docker image.

## Running the container
First, edit the docker-compose.yml environment variables BOT_TOKEN and CONNECTION_STR.
Then run the service:
    $ docker compose up 

## Deploying to Fly.io
1. Create a Fly.io account at https://fly.io/
2. Install the fly CLI. Instructions can be found at https://fly.io/docs/hands-on/install-flyctl/
3. In the project's root directory, run in the terminal
```
$ flyctl launch
```
Follow the instructions in the terminal to set up your application
4. Set secrets in fly for BOT_TOKEN and CONNECTION_STR with the following commands:
```
flyctl secrets set BOT_TOKEN=<your discord bot token>
flyctl secrets set CONNECTION_STR=<your supabase db connection string>
```
5. After editting flyctl secrets, navigate to the project's root directory and deploy the project with:
```
flyctl deploy
```

## Redeploying to Fly.io
1. After making your changes, simply navigate to the project's root directory and run:
```
flyctl deploy
```

## Bot commands:

| Command  | Description |
| ------------- | ------------- |
| !help  | Displays all bot commands and their descriptions.  |
| !fetch [card name in English or Korean] | Searches for a card by its English or Korean name.  |
