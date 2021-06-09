local mplayer = util.table.deepcopy(data.raw["character"]["character"])
mplayer.collision_box = {{0,0},{0,0}}
data:extend{mplayer}
