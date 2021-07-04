# fbot
Autonomous Factorio bot.

## Mod side:
- [x] Map export
  - [ ] Fast map export
- [x] Walk to x y
- [x] Drawbox
- [x] Place thing
- [x] Craft thing
- [x] Set recipe
- [x] Mine resource
  - [ ] Handle when resource path runs out
- [x] Put item into a thing
- [ ] Take item from a thing
- [x] Mine rocks & trees in an area
  - [ ] Split the area into smaller chunks to optimize the path when clearing
- [ ] Explore area
- [ ] Sometimes tries to mine fire

## Go side:
- [x] RCON
 - [x] Map export
 - [x] Walk to
 - [x] Draw box
 - [x] Craft thing
- [ ] Mapper
  - [x] Read exported map and generate lookup maps
  - [ ] Check for patch intersections and somehow resolve them
  - [x] Area allocating / deallocating
  - [x] Belt path finding
- [ ] Tasks
  - [ ] Task lists
  - [ ] Task dependencies
  - [ ] ...
- [ ] Building
- [ ] Unit tests
- [ ] Command shell

## Bugs:
- Starting the bot while the cutscene is running and then skipping the cutscene causes desync error.
