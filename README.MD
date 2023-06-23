<h1 align="center">
  Game Server Selector
</h1>


Game Server Selector is an open-source project that provides a comprehensive list of game servers, allowing users to choose the perfect match based on their preferences. With an easy-to-use interface and customizable parameters, finding 

## Installation

https://github.com/fatihkahveci/game-server-selector/releases download latest release. Configure your .env file and simply run the binary.

Note: If you need SSL support you need to create .env file and set SSL_KEY and SSL_CERT path.

## Features

- Game Engine Agnostic
- Currently supported search types is Match, Equal, Not Equal, Greater Than, Greater Than Or Equal, Less Than, Less Than Or Equal
- Optional HTTPS support
- Custom property support


## Example

Create a new game which is map de_dust_2 and game_mod is ranked with 20 player capacity.

```bash
curl --request POST \
>   --url http://localhost:3000/v1/server/create \
>   --header 'Content-Type: application/json' \
>   --data '{
>     "name":"Dust 2 - Ranked",
>     "ip":"127.0.0.1",
>     "capacity":20,
>     "current_player_count": 1,
>     "port":7070,
>     "custom_data": {
>         "region": "US",
>         "map": "de_dust_2",
>         "game_mod":"ranked"
>     }
> }'
```

Create another de_dust_2 map but this time game_mod is gun_race

```bash
 curl --request POST \
>   --url http://localhost:3000/v1/server/create \
>   --header 'Content-Type: application/json' \
>   --data '{
>     "name":"Dust 2 - GunRace",
>     "ip":"127.0.0.1",
>     "capacity":20,
>     "current_player_count": 1,
>     "port":7070,
>     "custom_data": {
>         "region": "EU",
>         "map": "de_dust_2",
>         "game_mod":"gun_race"
>     }
> }'
```

And finaly search your custom pref.

```bash
curl --request POST \
>   --url http://localhost:3000/v1/server/search \
>   --header 'Content-Type: application/json' \
>   --data '
>     [
>         {
>             "field": "game_mod",
>             "query": {
>                 "operator": "eq",
>                 "value": "gun_race"
>             }
>         },
>         {
>             "field": "capacity",
>             "query": {
>                 "operator": "gt",
>                 "value": 1
>             }
>         }
>     ]
> '
```

## Unity Example

https://github.com/fatihkahveci/gs-selector-unity

## TODOS

- Tests..
- More storage support.

## Contributing

If you would like to contribute to this project, please follow these steps. For example:

1. Fork this project.
2. Create a new branch to add a new feature or fix a bug.
3. Make your changes and submit them as a pull request.

## LICENSE

Game Server Selector is released under MIT license.