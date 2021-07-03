# Early game notes

A nice map:
```
>>>eNpjZGBkUAZiIGiwB5EcLMn5iTkQHgRzJecXFKQW6eYXpSILc
yYXlaak6uZnoipOzUvNrdRNSixOhZkIojkyi/Lz0E1gLS7Jz0MVK
SlKTS1GFuEuLUrMyyzNRdfLwFjadudyQ4scAwj/r2dQ+P8fhIGsB
0AbQZiBsQGiEigGBUwSyfl5JUX5ObrFqSUlmXnpVomlFVZJmYnFH
AZ6pgYgIIuuIjc/s7iktCgVVRlrck5mWhoDg4IjEDuBbWNgrBZZ5
/6waoo9I8Q2PQco4wNU5EASTMQTxvBzwCmlAmOYIJljDAafkRgQS
0uAVkBVcTggGBDJFpAkI2Pv260Lvh+7YMf4Z+XHS75JCfaMhq4i7
z4YrbMDSrKDvMAEJ2bNBIGdMK8wwMx8YA+VumnPePYMCLyxZ2QF6
RABEQ4WQOKANzMDowAfkLWgB0goyDDAnGYHM0bEgTENDL7BfPIYx
rhsj+4PYEDYgAyXAxEnQATYQrjLGCFMh34HRgd5mKwkQglQvxEDs
htSED48CbP2MJL9aA7BjAhkf6CJqDhgiQYukIUpcOIFM9w1wPC8w
A7jOcx3YGQGMUCqvgDFIDyQDMwoCC3gAA5umCwkbTDkT7+sDgBeR
rhb<<<
```

## Goal:

Item        | Iron | Copper | Stone
------------|------|--------|-------
Stem engine |   31 |      0 |     0
Boiler      |    4 |      0 |     5
Pump        |    5 |      3 |     0

Total raw:
- 40 Iron (already have 8 from the ship)
- 3 Copper
- 5 Stone

## Tasklist:
- mine the ship
  - `/cleararea {"area":[[-60,-9],[2,20]],"t":"all"}`
  - FIX: It tries to mine the fire and some particles
- mine 32 Iron
  - `/mineresource {"pos":[-92.5,30.5],"name":"iron-ore","amount":40}`
- mine 3 Copper
  - `/mineresource {"pos":[-99.5,-35.5],"name":"copper-ore","amount":3}`
- mine 5 Stone
  - `/mineresource {"pos":[-53.5,38.5],"name":"stone","amount":5}`
- mine 3 Coal
  - `/mineresource {"pos":[-53.5,38.5],"name":"stone","amount":5}`
- Place the furnace
  - `/place {"pos":[-56,12],"item":"stone-furnace"}`
- Put the coal into the furnace
  - `/put {"pos":[-56,12],"amount":3,"item":"coal","slot":1}`
- Smelt iron
  - `/put {"pos":[-56,12],"amount":32,"item":"iron-ore","slot":2}`
- Wait for the iron th be smelted
  - TODO
- Take out the iron
  - `/take {"pos":[-56,12],"amount":32,"item":"iron-plate","slot":3}`
- Smelt copper
  - `/put {"pos":[-56,12],"amount":3,"item":"copper-ore","slot":2}`
- Wait for the copper to be smelted
  - TODO
- Take out the iron
  - `/take {"pos":[-56,12],"amount":3,"item":"copper-plate","slot":3}`
- Craft the Steam engine
  - `/craft {"item":"steam-engine","count":1}`
  - TODO: indicate the crafting process, output the crafting state into state.json
- Craft the Boiler
  - `/craft {"item":"boiler","count":1}`
- Craft the Steam engine
  - `/craft {"item":"offshore-pump","count":1}`
- Mine the furnace
  - `/mine [-56,12]`

