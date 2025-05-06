# Pokedex-cli

## Building and running

    go build .
    ./pokedex-cli

Commands:

    help:                    Lists the available commands
    pokedex:                 View all Pokemon in your Pokedex
    map:                     Displays the next 20 map areas
    mapb:                    Displays the previous 20 map areas
    explore <location_area>: Displays pokemon in the location area
    catch <pokemon_name>:    Attempt to catch the specified Pokemon
    inspect <pokemon_name>:  View a caught Pokemon's stats
    exit:                    Exits the Pokedex

## Extending the project TODO:

* Update the CLI to support the "up" arrow to cycle through previous commands
* Simulate battles between pokemon
* Add more unit tests
* Keep pokemon in a "party" and allow them to level up
* Allow for pokemon that are caught to evolve after a set amount of time
* Persist a user's Pokedex to disk so they can save progress between sessions
* Use the PokeAPI to make exploration more interesting. For example, rather than typing the names of areas, maybe you are given choices of areas and just type "left" or "right"
* Random encounters with wild pokemon
* Adding support for different types of balls (Pokeballs, Great Balls, Ultra Balls, etc), which have different chances of catching pokemon
