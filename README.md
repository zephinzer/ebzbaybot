# EbzBot

This is an unofficial Ebisus Bay Bot for easily monitoring floor price changes on Telegram.

Talk to the bot at https://t.me/ebzbaybot for a customised experience. Otherwise, pre-set channels for various collections include (basically collections I watch):

1. [Mad Meerkats](https://t.me/ebzbay_mmkfp)
2. [MM Treehouses](https://t.me/ebzbay_mmtfp)
3. [Cronos Primate Business](https://t.me/ebzbay_cpbfp)
4. [Nifty Nanas](https://t.me/ebzbay_nanasfp)
5. [Lolita & Friends](https://t.me/ebzbay_lolitasfp)

# Contributing

## Running this locally

1. Run `make env` to bring up the POstgres database with Docker Compose.
2. Run `make deps` to pull in the dependencies.
3. Make a copy of `.envrc.sample` as `.envrc` and insert your Telegram bot token (get it from [The BotFather](https://t.me/BotFather)).
4. Run `direnv allow .` to enable the `.envrc` file to be loaded
5. Run `go run. start` to start the bot.

## Adding to the collection whitelist

The whitelisted collection list can be found at [`./pkg/constants/data.json`](./pkg/constants/data.json). That should be where you're raising a pull request to whitelist a collection.

## Releasing this service

Run `make release release_tag=$(git rev-parse HEAD | head -c 8)` to create and push the artifact.

## Deploying this service

This service's tooling was created around deploying on Digital Ocean (DO) App Platform because it's simpler and more affordable.

To create the required resources, go to the DO App Platform and use a Docker image as the type of application. Use `zephinzer/ebzbaybot` (find the latest image tag on [DockerHub](https://hub.docker.com/r/zephinzer/ebzbaybot) as the image source). Also, create and attach a database to it and change the default injected environment variable from `DATABASE_URL` to `POSTGRES_URL`.

To deploy a new release, navigate to the application in DO and change the image repository tag. The application should automagically deploy itself.


# Licensing

Code is licensed [under GPLv3](https://www.gnu.org/licenses/gpl-3.0) which basically means you can deploy this software AS-IS by yourself anywhere if you so wish. Also, modifications are allowed but must be open-sourced as well.
