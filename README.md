# Pixel world simulation

## How to run

1. Install the [Go 1.19 compiler](https://go.dev/dl/).
2. Run `go run main.go` in the root directory of the project.

## Materials

### Sand

Yellow particle which falls down and fills the space below it.
Sinks in the water.

### Water

Light blue particle which falls down and fills the space below it (+horizontally).
Creates the *steam* when it reaches the *fire*.
Makes the *grass* grow faster.

### Wood

Dark brown particle which doesn't move, but can be destroyed by the *fire*.

### Fire

Light red particle which burns flammable materials (wood, grass, leaves).
Creates the *smoke* while burning.
Replaces any particle under the cursor.

### Grass

Green particle which grows.
Dies when it can't grow anymore.
Grows faster when it consumes the *water*.

### Smoke

Light gray particle which rises up and disappears after some time.

### Steam

Dark blue particle which rises up and disappears after some time creating the *water* particles.

### Metal

Dark grey particle which doesn't move, and it is quite hard to destroy it (is it so?).

### Dirt

Brown particle which doesn't move, but can be destroyed by other heavier particles (including itself).

### Graviton

Purple particle which doesn't move and attracts other particles.
Has a limited range of attraction.

## Controls

### Mouse

#### Left click

Create a new particle (particles if the *circle* tool is selected) at the mouse position / remove.

#### Right click

Select a tool / toggle options on the right side of the screen.

### Keyboard

- `[1 - 9]` - select a *material* to draw with;
- `s` - switch between drawing with a *single dot* and a *circle*;
- `d` - select the *removal tool*;
- `q` - reduce the *circle* tool radius;
- `e` - increase the *circle* tool radius;
- `f` - switch on/off the *apply random force* mode;
- `z` - invert the gravity;

## To try

1. Place multiple single *graviton* particles and pour *sand* or *water* particles near them.
2. Turn on the *circle* tool, select the *fire* material and turn on the *apply random force* mode. Now you can do fireworks.
