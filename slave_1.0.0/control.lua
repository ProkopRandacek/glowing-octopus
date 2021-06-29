require("util")

local instant_mine = false

local walking_state = {walking = false} --  what direction we are walking and if were walking
local walking_to = {} -- the position that were walking to

local mining = false
local mining_target = nil

local clearing = false
local clearing_target = nil
local clearing_area = nil

local building = false
local building_item = nil
local building_pos = nil
local building_dir = nil

local bot;
local surface;

--- ============ ---
--- === UTIL === ---
--- ============ ---

function has_value(tab, val) -- check if key is in a table
	for index, value in ipairs(tab) do
		if value == val then return true end
	end
	return false
end

function get_dir(x, y, eps) -- returns the direction from player position to {x. y}
	local delta_x = x - bot.position.x
	local delta_y = y - bot.position.y
	if delta_x > eps then
		if     delta_y >  eps then return {walking = true, direction = defines.direction.southeast}
		elseif delta_y < -eps then return {walking = true, direction = defines.direction.northeast}
		else                       return {walking = true, direction = defines.direction.east}
		end
	elseif delta_x < -eps then
		if     delta_y >  eps then return {walking = true, direction = defines.direction.southwest}
		elseif delta_y < -eps then return {walking = true, direction = defines.direction.northwest}
		else                       return {walking = true, direction = defines.direction.west}
		end
	else
		if     delta_y >  eps then return {walking = true,  direction = defines.direction.south}
		elseif delta_y < -eps then return {walking = true,  direction = defines.direction.north}
		else                       return {walking = false, direction = defines.direction.north}
		end
	end
end

--- =================== ---
--- ===  FUNCTIONS  === ---
--- =================== ---

function mark_area(color, text, x1, y1, x2, y2)
	rendering.draw_rectangle{
		color = color,
		filled = true,
		left_top = {x1, y1},
		right_bottom = {x2, y2},
		time_to_live = 600,
		surface = surface,
	}
end

function write_resrc(area)
	ids = { ["iron-ore"] = 1, ["copper-ore"] = 2, ["coal"] = 3, ["stone"] = 4, ["uranium-ore"] = 5, ["crude-oil"] = 6}
	output = {{},{},{},{},{},{}}
	for i, e in pairs(game.surfaces[1].find_entities_filtered{area = area, type="resource"}) do
		table.insert(output[ids[e.name]], {x = e.position.x, y = e.position.y, a = e.amount})
	end
	game.write_file("resrc.json", game.table_to_json(output))
	game.print("resources export done")
end

function write_water(area)
	output = {}
	for i, t in pairs(game.surfaces[1].find_tiles_filtered{area = area, name = {"water", "deepwater"}}) do
		output[i] = { t.position.x, t.position.y }
	end
	game.write_file("water.json", game.table_to_json(output))
	game.print("water export done")
end

function write_trees(area)
	output = {}
	for i, e in pairs(game.surfaces[1].find_entities_filtered{area = area, type="tree"}) do
		output[i] = {e.position.x, e.position.y, e.prototype.mineable_properties.products[1].amount}
	end
	game.write_file("trees.json", game.table_to_json(output))
	game.print("trees export done")
end

function write_rocks() -- the ammount of rocks is tiny, this function can export all rocks in huge world in a fraction of a second
	output = {}
	for i, e in pairs(game.surfaces[1].find_entities_filtered{name={"rock-huge", "sand-rock-big", "rock-big"}}) do
		products = e.prototype.mineable_properties.products
		for i, e in pairs(products) do
			e["type"] = nil
			e["probability"] = nil
		end
		output[i] = {e.position.x, e.position.y, products}
	end
	game.write_file("rocks.json", game.table_to_json(output))
	game.print("rocks export done")
end

function write_state()
	state = {
		["walking_state"] = walking_state.walking,
		["mining_state"] = mining,
		--["position"] = game.players[1].position
	}
	game.write_file("state.json", game.table_to_json(state) .. "\n")
end

function walkto(pos)
	walking_state = get_dir(pos[1], pos[2], 0.2) -- this sets the walking state to true
	if walking_state.walking then -- maybe were already there
		game.print("walk to " .. game.table_to_json(pos))
		write_state()
		walking_to = pos
	end
end

function mine(pos)
	game.print("mine at " .. game.table_to_json(pos))
	walkto({pos[1] + 1.5, pos[2]})
	mining = true
	mining_target = surface.find_entities_filtered{limit=1, position=pos, radius=1.0, type={"player", "corpse", "character", "flying-text"}, invert=true}[1]
	if mining_target == nil then
		game.print("mine target entity not found")
		return
	end
end

function clear(area)
	clearing_area = area

	clearing_targets = bot.surface.find_entities_filtered{area=clearing_area, type="tree", limit=2}
	if #clearing_targets == 0 then
		game.print("nothing to clear")
		return
	end

	clearing_target = clearing_targets[1]

	clearing = true
end

function place(pos, item, dir) -- just place it
	-- Check if we can actually place the item at this tile
	local placed = false

	if bot.can_place_entity{name=item, position=pos, direction=direction} then
		if bot.surface.can_fast_replace{name=item, position=pos, direction=dir, force="player"} then
			placed = bot.surface.create_entity{name=item, position=pos, direction=dir, force="player", fast_replace=true, player=p}
		else
			placed = bot.surface.create_entity{name=item, position=pos, direction=dir, force="player"}
		end
	else
		game.print("cannot place: " .. item .. " at " .. game.table_to_json(pos))
		return false
	end
	if placed then
		bot.remove_item({name = item, count = 1})
		game.print("placed " .. item .. " at " .. game.table_to_json(pos))
	end
	return true
end

function build(pos, item, dir) -- walk there and then place it
	if bot.get_item_count(item) == 0 then
		game.print("cant build " .. item .. " because i dont have it")
		return
	end

	walkto(pos)

	building = true
	building_item = item
	building_pos = pos
	building_dir = dir
end

--- ================== ---
--- ===  COMMANDS  === ---
--- ================== ---

commands.add_command("walkto", nil, function(command)
	local args = game.json_to_table(command.parameter)
	walkto(args)
end)

commands.add_command("mine", nil, function(command)
	local a = game.json_to_table(command.parameter)
	mine(a)
end)

commands.add_command("cleararea", nil, function(command)
	local a = game.json_to_table(command.parameter)
	mark_area({r = 1, a = 0.05}, nil, a[1][1], a[1][2], a[2][1], a[2][2])
	clear(a)
end)

commands.add_command("build", nil, function(command)
	local a = game.json_to_table(command.parameter)
	if a.dir == nil then a.dir = 0 end
	build(a.pos, a.item, a.dir)
end)

commands.add_command("writeresrc", nil, function(command) write_resrc(game.json_to_table(command.parameter)) end)
commands.add_command("writewater", nil, function(command) write_water(game.json_to_table(command.parameter)) end)
commands.add_command("writetrees", nil, function(command) write_water(game.json_to_table(command.parameter)) end)
commands.add_command("writerocks", nil, function(command) write_rocks(                                     ) end)

commands.add_command("drawbox", nil, function(command)
	local a = game.json_to_table(command.parameter)
	mark_area(a.color, nil, a.x1, a.y1, a.x2, a.y2)
end)

commands.add_command("craft", nil, function(command)
	local a = game.json_to_table(command.parameter)
	bot.begin_crafting{recipe=a.recipe, count=a.count}
end)

--- ================ ---
--- ===  EVENTS  === ---
--- ================ ---

script.on_event(defines.events.on_tick, function(event)
	bot = game.players[1];
	surface = game.surfaces[1];

	if walking_state.walking then -- update player walking state
		bot.walking_state = walking_state
		walking_state = get_dir(walking_to[1], walking_to[2], 0.2)
		if not walking_state.walking then
			bot.walking_state = walking_state -- stop the player
			game.print("walking done")
		end
	elseif mining then
		if instant_mine then
			result = bot.mine_entity(mining_target, true)
			if result == false then game.print("failed to mine " .. mining_target.type .. " " .. mining_target.name) end
		else
			bot.update_selected_entity(mining_target.position)
			bot.mining_state = {mining = true, position = mining_target.position}
		end
	elseif building then
		place(building_pos, building_item, building_dir)
		building = false
	elseif clearing then
		clearing_target = game.surfaces[1].find_entities_filtered{area = clearing_area, type="tree", limit=1}[1]

		if clearing_target == nil then -- if still no tree then clearing done
			clearing = false
			game.print("clearing done")
		else
			game.print("found next tree to clear")
			mine({clearing_target.position.x, clearing_target.position.y})
		end
	end
end)

script.on_nth_tick(60, function(event) -- update state once per second
	bot = game.players[1];
	write_state()
end)

script.on_event(defines.events.on_player_mined_entity, function(event)
	mining = false
	game.print("mining done")
end)

