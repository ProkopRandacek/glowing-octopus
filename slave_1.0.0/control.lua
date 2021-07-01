require("util")

local tick = 0
local inited = false

local walking_state = {walking = false} --  what direction we are walking and if were walking
local walking_to = {} -- the position that were walking to

local mining = false
local mining_target = nil

local resource_mining = false
local resource_mining_name = nil
local resource_mining_target = nil
local resource_mining_amount = nil

local clearing = false
local clearing_target = nil
local clearing_area = nil
local clearing_type = "nature" -- "nature" or "all"

local placing = false
local placing_item = nil
local placing_pos = nil
local placing_dir = nil

local puting = false
local puting_pos = nil
local puting_item = nil
local puting_amount = nil
local puting_slot = nil

local building = false
local building_queue = nil

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

function check_busy()
	return not (inited and not (walking_state.walking or mining or clearing or placing or building))
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

function write_rocks() -- the amount of rocks is tiny, this function can export all rocks in huge world in a fraction of a second
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
		["position"] = bot.position,
		["walking_state"] = walking_state.walking,
		["mining_state"] = mining,
		["placing_state"] = placing,
		["clearing_state"] = clearing,
		["building_state"] = building
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
	mining_target = surface.find_entities_filtered{limit=1, position=pos, radius=1.0, type={"player", "corpse", "character", "flying-text", "resource", "fish"}, invert=true}[1]
	if mining_target == nil then
		game.print("mine target entity not found")
		return
	end
end

function mine_resource(pos, amount, name)
	game.print("mine resource at " .. game.table_to_json(pos))
	resource_mining_target = surface.find_entities_filtered{limit=1, position=pos, radius=1.0, type="resource"}[1]
	resource_mining_amount = amount
	resource_mining_name = name
	if resource_mining_target == nil then
		game.print("resource not found")
		return
	elseif bot.get_item_count(resource_mining_name) >= resource_mining_amount then
		game.print("aleady have the resources")
	else -- only walk if resource found
		walkto(pos)
		resource_mining = true
	end
end

function clear(area, t)
	clearing_area = area
	clearing_type = t

	if clearing_type == "all" then
		clearing_targets = surface.find_entities_filtered{area=clearing_area, type={"player", "corpse", "character", "flying-text", "resource", "fish"}, invert=true, limit=1}
	elseif clearing_type == "nature" then
		clearing_targets = surface.find_entities_filtered{area=clearing_area, type={"tree", "simple-entity"}, limit=1}
	end
	if #clearing_targets == 0 then
		game.print("nothing to clear")
		return
	end

	clearing_target = clearing_targets[1]

	clearing = true
end

function place(pos, item, dir, recipe) -- just place it
	-- Check if we can actually place the item at this tile
	local placed = false

	if surface.can_place_entity{name=item, position=pos, direction=direction} then
		if surface.can_fast_replace{name=item, position=pos, direction=dir, force="player"} then
			placed = surface.create_entity{name=item, position=pos, direction=dir, recipe=recipe, force="player", fast_replace=true, player=bot}
		else
			placed = surface.create_entity{name=item, position=pos, direction=dir, recipe=recipe, force="player"}
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

function place_safe(pos, item, dir, recipe) -- walk there and then place it
	if bot.get_item_count(item) == 0 then
		game.print("cant place " .. item .. " because i dont have it")
		return
	end

	walkto(pos)

	placing = true
	placing_item = item
	placing_pos = pos
	placing_dir = dir
	placing_recipe = recipe
end

-- slot: https://lua-api.factorio.com/latest/defines.html#defines.inventory
function putin(position, item, amount, slot)
	bot.update_selected_entity(position)

	local amountininventory = bot.get_item_count(item)
	local otherinv = bot.selected.get_inventory(slot)
	local toinsert = math.min(amountininventory, amount)

	if toinsert == 0 then
		game.print("nothing to insert cuz toinsert == 0")
		return true
	end
	if not otherinv then
		game.print("no slot")
		return false
	end

	local inserted = otherinv.insert{name=item, count=toinsert}
	--if we already failed for trying to insert no items, then if no items were inserted, it must be because it is full
	if inserted == 0 then
		game.print("nothing to insert cuz inserted == 0")
		return true
	end

	bot.remove_item{name=item, count=inserted}
	return true
end

function putin_safe(pos, item, amount, slot)
	walkto(pos)

	puting = true
	puting_pos = pos
	puting_item = item
	puting_amount = amount
	puting_slot = slot
end

--- ================== ---
--- ===  COMMANDS  === ---
--- ================== ---

commands.add_command("writeresrc", nil, function(command) write_resrc(game.json_to_table(command.parameter)) end)
commands.add_command("writewater", nil, function(command) write_water(game.json_to_table(command.parameter)) end)
commands.add_command("writetrees", nil, function(command) write_water(game.json_to_table(command.parameter)) end)
commands.add_command("writerocks", nil, function(command) write_rocks(                                     ) end)

commands.add_command("walkto", nil, function(command)
	if (check_busy()) then
		game.print("im busy")
	else
		local args = game.json_to_table(command.parameter)
		walkto(args)
	end
end)

commands.add_command("mineresource", nil, function(command)
	if (check_busy()) then
		game.print("im busy")
	else
		local a = game.json_to_table(command.parameter)
		mine_resource(a.pos, a.amount, a.name)
	end
end)

commands.add_command("mine", nil, function(command)
	if (check_busy()) then
		game.print("im busy")
	else
		local a = game.json_to_table(command.parameter)
		mine(a)
	end
end)

commands.add_command("cleararea", nil, function(command)
	if (check_busy()) then
		game.print("im busy")
	else
		local a = game.json_to_table(command.parameter)
		mark_area({r = 1, a = 0.05}, nil, a.area[1][1], a.area[1][2], a.area[2][1], a.area[2][2])
		clear(a.area, a.t)
	end
end)

commands.add_command("place", nil, function(command)
	if (check_busy()) then
		game.print("im busy")
	else
		local a = game.json_to_table(command.parameter)
		if a.dir == nil then a.dir = 0 end
		place_safe(a.pos, a.item, a.dir)
	end
end)

commands.add_command("build", nil, function(command)
	if (check_busy()) then
		game.print("im busy")
	else
		building_queue = game.json_to_table(command.parameter)
		building = true
	end
end)

commands.add_command("drawbox", nil, function(command)
	if (check_busy()) then
		game.print("im busy")
	else
		local a = game.json_to_table(command.parameter)
		mark_area(a.color, nil, a.x1, a.y1, a.x2, a.y2)
	end
end)

commands.add_command("craft", nil, function(command)
	if (check_busy()) then
		game.print("im busy")
	else
		local a = game.json_to_table(command.parameter)
		bot.begin_crafting{recipe=a.recipe, count=a.count}
	end
end)

commands.add_command("put", nil, function(command)
	if (check_busy()) then
		game.print("im busy")
	else
		local a = game.json_to_table(command.parameter)
		game.print(command.parameter)
		putin_safe(a.pos, a.item, a.amount, a.slot)
	end
end)

--- ================ ---
--- ===  EVENTS  === ---
--- ================ ---

function init()
	surface = game.surfaces[1]

	if global.fbot == nil then
		game.print("no fbot found. Imma create it")
		bot = surface.create_entity{
			name="character",
			position={3, 2},
			direction=3,
			force="player",
			fast_replace=true
		}
		bot.insert({name="burner-mining-drill", count=1}) -- the starting inventory
		bot.insert({name="stone-furnace", count=1})
		bot.color = {r=1, g=1, b=1}
		global.fbot = bot
	else
		game.print("fbot already exists imma use it")
		bot = global.fbot
	end
	game.print("init done")
	inited = true
end

script.on_event(defines.events.on_tick, function(event)
	if tick == 0 then init() end
	tick = tick + 1
	if not inited then return end

	--[[thing = surface.find_entities_filtered{position=game.players[1].position, radius=1, type={"character"}, invert=true, limit=1}[1]
	if thing ~= nil then
		game.print(thing.type .. " " .. thing.name)
	end]]
	bot.color = {
		r=(math.sin(tick / 10.0) * 0.5 + 0.5) * 255.0,
		g=(math.cos(tick / 20.0) * 0.5 + 0.5) * 255.0,
		b=(math.sin(tick / 30.0) * 0.5 + 0.5) * 255.0
	}

	if walking_state.walking then -- update player walking state
		bot.walking_state = walking_state
		walking_state = get_dir(walking_to[1], walking_to[2], 0.2)
		if not walking_state.walking then
			bot.walking_state = walking_state -- stop the player
			game.print("walking done")
		end
	elseif puting then
		putin(puting_pos, puting_item, puting_amount, puting_slot)
		puting = false
	elseif resource_mining then
		i_have = bot.get_item_count(resource_mining_name)
		i_need = resource_mining_amount
		if i_have < i_need then
			game.print(i_have .. " / " .. i_need .. " " .. resource_mining_name)
			bot.update_selected_entity(resource_mining_target.position)
			bot.mining_state = {mining = true, position = resource_mining_target.position}
		else
			resource_mining = false
			bot.mining_state = { mining = false }
			game.print("resource mining done")
		end
	elseif mining then
		if mining_target ~= nil and mining_target.valid then
			bot.update_selected_entity(mining_target.position)
			bot.mining_state = { mining = true, position = mining_target.position }
		else
			mining = false
			game.print("mining done")
		end
	elseif placing then
		place(placing_pos, placing_item, placing_dir, placing_recipe)
		placing = false
	elseif building then
		if #building_queue == 0 then
			building = false
			game.print("building done")
		else
			if type(building_queue) == 4 then
				game.print("building_queue is a string instead of a stable O.o") -- this happens for some reason
			else
				b = table.remove(building_queue, 1)
				place_safe({b.pos.x, b.pos.y}, b.name, b.rotation, b.recipe)
			end
		end
	elseif clearing then
		if clearing_type == "all" then
			game.print("clearing all")
			clearing_target = surface.find_entities_filtered{area=clearing_area, type={"player", "corpse", "character", "flying-text", "resource", "fish"}, invert=true, limit=1}[1]
		elseif clearing_type == "nature" then
			clearing_target = surface.find_entities_filtered{area=clearing_area, type="tree", limit=1}[1]
		end

		if clearing_target == nil then
			clearing = false
			game.print("clearing done")
		else
			game.print("found next thing to clear. Its a " .. clearing_target.name)
			mine({clearing_target.position.x, clearing_target.position.y})
		end
	end
end)

script.on_nth_tick(60, function(event) -- update state once per second
	if inited then write_state() end
end)

