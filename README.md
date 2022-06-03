# Project 92

I don't know what this will be... Maybe a some kind of capture the flag

- Build "stuff" with attack/strength
- "attack other things?
- stratego based ?
- does terrain do something?

Running:

    go run .

# Keyboard usage

```
Keys:
    <space>     does a sonar ping (i don't know what this will do)
    r           refreshes the screen
    arrows      move something, depending on move mode
    tab         change move mode
```

# Move modes

There are 3 move modes:

- player      moves the player
- map         moves the map
- menu        browse through menu (not implemented yet)

You see the current movemode at the left top of the screen



  - We have to separate the map itself from the viewport / screen

  - When we do things like ping or add players, they should be added on top of 
     the terrain. Maybe have a separate field in the struct:
```
    {
       // color
       // rune
       // type (of the terrain)
       // elements (placed on top of the terrain)
    }
```

  - Other stuff




screen structure:
    height
    width
    viewport of map
    viewport of menu ?
    other elements


game:
    map
    player
    stuff


draw(game, screen)
