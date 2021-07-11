require("util")

local water_chunk_size = 32
local dont_mine = {"player", "corpse", "character", "flying-text", "resource", "fish", "smoke-with-trigger", "fire", "explosion"}

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
local placing_ugbt = nil -- underground belt type

local puting = false
local puting_pos = nil
local puting_item = nil
local puting_amount = nil
local puting_slot = nil

local taking = false
local taking_pos = nil
local taking_item = nil
local taking_amount = nil
local taking_slot = nil

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
	if x1 > x2 or y1 > y2 then
		game.print("wrong corners")
	end
	rendering.draw_rectangle{
		color = color,
		filled = false,
		left_top = {x1, y1},
		right_bottom = {x2, y2},
		time_to_live = 600,
		surface = surface,
	}
end

function mark_pos(color, text, x, y)
	rendering.draw_circle{
		color = color,
		filled = true,
		radius = 0.3,
		target = {x, y},
		time_to_live = 600,
		surface = surface,
	}
end

function write_resrc(area)
	output = { ["iron-ore"] = {}, ["copper-ore"] = {}, ["coal"] = {}, ["stone"] = {}, ["uranium-ore"] = {}, ["crude-oil"] = {},
			   ["iron-ore-a"] = {}, ["copper-ore-a"] = {}, ["coal-a"] = {}, ["stone-a"] = {}, ["uranium-ore-a"] = {}, ["crude-oil-a"] = {}}
	for i, e in pairs(game.surfaces[1].find_entities_filtered{area = area, type="resource"}) do
		table.insert(output[e.name], {x = e.position.x, y = e.position.y})
		table.insert(output[e.name .. "-a"], e.amount)
	end
	for k, v in pairs(output) do -- remove empty keys
		if #v == 0 then
			output[k] = nil
		end
	end
	game.write_file("resrc.json", game.table_to_json(output))
	game.print("resources export done")
end

function write_water(area)
	output = {}
	tl_chunk = {math.floor(area[1][1] / water_chunk_size) * water_chunk_size, math.floor(area[1][2] / water_chunk_size) * water_chunk_size}
	br_chunk = {math.ceil(area[2][1] / water_chunk_size) * water_chunk_size, math.ceil(area[2][2] / water_chunk_size) * water_chunk_size}
	game.print(game.table_to_json(tl_chunk))
	game.print(game.table_to_json(br_chunk))

	for x = tl_chunk[1],br_chunk[1],water_chunk_size do
		for y = tl_chunk[2],br_chunk[2],water_chunk_size do
			if #surface.find_tiles_filtered{area={{x,y},{x+water_chunk_size,y+water_chunk_size}},name={"water", "deepwater"},limit=1} >= 1 then
				table.insert(output, {x = x, y = y})
			end
		end
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
	bot.get_main_inventory().sort_and_merge()
	state = {
		["position"] = bot.position,
		["walking_state"] = walking_state.walking,
		["mining_state"] = mining,
		["mining_resource_state"] = resource_mining,
		["placing_state"] = placing,
		["puting_state"] = puting,
		["taking_state"] = taking,
		["clearing_state"] = clearing,
		["building_state"] = building,
		["inv"] = bot.get_main_inventory().get_contents()
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
	mining_target = surface.find_entities_filtered{limit=1, position=pos, radius=1.0, type=dont_mine, invert=true}[1]
	if mining_target == nil then
		game.print("mine target entity not found")
		return
	end
end

function mine_resource(pos, amount, name)
	game.print("mine resource at " .. game.table_to_json(pos))
	resource_mining_target = surface.find_entities_filtered{limit=1, position=pos, radius=0.2, type="resource"}[1]
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
		clearing_targets = surface.find_entities_filtered{area=clearing_area, type=dont_mine, invert=true, limit=1}
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

function place(pos, item, dir, recipe, ugbt) -- just place it
	-- Check if we can actually place the item at this tile
	local placed = false

	if surface.can_place_entity{name=item, position=pos, direction=direction, type=ugbt} then
		if surface.can_fast_replace{name=item, position=pos, direction=dir, force="player", type=ugbt} then
			placed = surface.create_entity{name=item, position=pos, direction=dir, recipe=recipe, force="player", fast_replace=true, player=bot, type=ugbt}
		else
			placed = surface.create_entity{name=item, position=pos, direction=dir, recipe=recipe, force="player", type=ugbt}
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

function place_safe(pos, item, dir, recipe, ugbt) -- walk there and then place it
	if not game.item_prototypes[item] == nil then
		game.print("'" .. item .. "' is not valid item name")
		return
	end
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
	placing_ugbt = ugbt
end

-- slot: https://lua-api.factorio.com/latest/defines.html#defines.inventory
function putin(position, item, amount, slot)
	bot.update_selected_entity(position)

	local amountininventory = bot.get_item_count(item)
	local otherinv = bot.selected.get_inventory(slot)
	local toinsert = math.min(amountininventory, amount)

	if toinsert == 0 then
		game.print("nothing to insert cuz toinsert == 0")
		return 0
	end
	if not otherinv then
		game.print("no slot")
		return false
	end

	local inserted = otherinv.insert{name=item, count=toinsert}
	--if we already failed for trying to insert no items, then if no items were inserted, it must be because it is full
	if inserted == 0 then
		game.print("nothing to insert cuz inserted == 0")
		return 0
	end

	bot.remove_item{name=item, count=inserted}
	return inserted
end

function takeout(position, item, amount, slot)
	bot.update_selected_entity(position)

	local otherinv = bot.selected.get_inventory(slot)
	local amountintarget = otherinv.get_item_count(item)
	local totake = math.min(amountintarget, amount)

	if not otherinv then
		game.print("no slot")
		return
	end
	if totake == 0 then
		game.print("waiting for smelter")
		return 0
	end

	local taken = bot.insert{name=item, count=totake}

	otherinv.remove{name=item, count=taken}
	return taken
end

function putin_safe(pos, item, amount, slot)
	walkto(pos)

	puting = true
	puting_pos = pos
	puting_item = item
	puting_amount = amount
	puting_slot = slot
end

function takeout_safe(pos, item, amount, slot)
	walkto(pos)

	taking = true
	taking_pos = pos
	taking_item = item
	taking_amount = amount
	taking_slot = slot
end

--- ================== ---
--- ===  COMMANDS  === ---
--- ================== ---

commands.add_command("writeresrc", nil, function(command) write_resrc(game.json_to_table(command.parameter)) end)
commands.add_command("writewater", nil, function(command) write_water(game.json_to_table(command.parameter)) end)
commands.add_command("writetrees", nil, function(command) write_trees(game.json_to_table(command.parameter)) end)
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
		place_safe(a.pos, a.item, a.dir, a.ugbt)
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
	local a = game.json_to_table(command.parameter)
	game.print("drawbox!")
	mark_area(a.color, nil, a.x1, a.y1, a.x2, a.y2)
end)

commands.add_command("drawpoint", nil, function(command)
	local a = game.json_to_table(command.parameter)
	game.print("drawbox!")
	mark_pos(a.color, nil, a.x, a.y)
end)

commands.add_command("craft", nil, function(command)
	if (check_busy()) then
		game.print("im busy")
	else
		local a = game.json_to_table(command.parameter)
		game.print("started crafting " .. a.count .. " " .. a.recipe)
		bot.begin_crafting{recipe=a.item, count=a.count}
	end
end)

commands.add_command("put", nil, function(command)
	if (check_busy()) then
		game.print("im busy")
	else
		local a = game.json_to_table(command.parameter)
		putin_safe(a.pos, a.item, a.amount, a.slot)
	end
end)

commands.add_command("take", nil, function(command)
	if (check_busy()) then
		game.print("im busy")
	else
		local a = game.json_to_table(command.parameter)
		takeout_safe(a.pos, a.item, a.amount, a.slot)
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
		local put = putin(puting_pos, puting_item, puting_amount, puting_slot)
		-- if puting into furnace source slot and stuff same as taking
		if put ~= puting_amount and puting_slot == 2 then
			puting = true
			puting_amount = puting_amount - put
		else
			puting = false
		end
	elseif taking then
		if taking_amount == 0 then
			taking = false
			return
		end
		local taken = takeout(taking_pos, taking_item, taking_amount, taking_slot)
		-- if taking from furnace result slot && did not take enough yet, wait a tick and try to take again
		if taken ~= taking_amount and taking_slot == 3 then
			taking = true
			taking_amount = taking_amount - taken
		else
			taking = false
		end
	elseif resource_mining then
		i_have = bot.get_item_count(resource_mining_name)
		i_need = resource_mining_amount
		if i_have < i_need then
			game.print(i_have .. " / " .. i_need .. " " .. resource_mining_name .. " at " .. game.table_to_json(resource_mining_target.position))
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
			bot.mining_state = { mining = false }
			game.print("mining done")
		end
	elseif placing then
		place(placing_pos, placing_item, placing_dir, placing_recipe, placing_ugbt)
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
				if b.rotation == nil then b.rotation = 0 end
				if b.ugbt == nil then b.ugbt = "input" end
				place_safe({b.pos.x, b.pos.y}, b.name, b.rotation, b.recipe, b.ugbt)
			end
		end
	elseif clearing then
		if clearing_type == "all" then
			game.print("clearing all")
			clearing_target = surface.find_entities_filtered{area=clearing_area, type=dont_mine, invert=true, limit=1}[1]
		elseif clearing_type == "nature" then
			clearing_target = surface.find_entities_filtered{area=clearing_area, type="tree", limit=1}[1]
		end

		if clearing_target == nil then
			clearing = false
			mining = false
			game.print("clearing done")
		else
			game.print("found next thing to clear. Its a " .. clearing_target.type .. " " .. clearing_target.name)
			mine({clearing_target.position.x, clearing_target.position.y})
		end
	end
end)

script.on_nth_tick(60, function(event) -- update state once per second
	if inited then write_state() end
end)

