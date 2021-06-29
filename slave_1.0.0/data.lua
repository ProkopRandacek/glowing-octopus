local mplayer = util.table.deepcopy(data.raw["character"]["character"])
mplayer.collision_mask = {}
data:extend{mplayer}
