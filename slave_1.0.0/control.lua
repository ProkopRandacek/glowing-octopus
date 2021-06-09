require("util")

local walking_state = {walking = false}
local walking_to = {x = 0, y = 0}
local p;
local surface;

local function has_value(tab, val)
	for index, value in ipairs(tab) do
		if value == val then return true end
	end
	return false
end

local function write_world()
	local surface = game.get_surface("nauvis")
	local result = ""

	-- == water == --
	--local tiles = surface.find_tiles_filtered{name="water"}
	local tiles = surface.find_tiles_filtered{}
	local output = {}
	existingEntities = {}
	for i, t in pairs(tiles) do
		if not has_value(existingEntities, t.name) then
			table.insert(existingEntities, t.name)
		end
		--table.insert(output, { t.position.x, t.position.y })
	end
	--game.write_file("water" .. ".json", game.table_to_json(output))
	game.print("found resources: " .. game.table_to_json(existingEntities))
	--result = ""

	-- == ores == --
	local entities = surface.find_entities_filtered{type="resource"}
	output = {}
	existingEntities = {}
	for i, e in pairs(entities) do
		if not has_value(existingEntities, e.name) then
			output[e.name] = {}
			table.insert(existingEntities, e.name)
		end
		table.insert(output[e.name], {e.position.x, e.position.y, e.amount})
	end
	game.write_file("resources" .. ".json", game.table_to_json(output))
	game.print("found resources: " .. game.table_to_json(existingEntities))
end

function Split(s, delimiter)
	result = {};
	for match in (s..delimiter):gmatch("(.-)"..delimiter) do
		table.insert(result, match);
	end
	return result;
end

local function getdir(x, y)
	local delta_x = x - p.position.x
	local delta_y = y - p.position.y

	if delta_x > 0.2 then
		-- Easterly
		if     delta_y >  0.2 then return {walking = true, direction = defines.direction.southeast}
		elseif delta_y < -0.2 then return {walking = true, direction = defines.direction.northeast}
		else                       return {walking = true, direction = defines.direction.east}
		end
	elseif delta_x < -0.2 then
		-- Westerly
		if     delta_y >  0.2 then return {walking = true, direction = defines.direction.southwest}
		elseif delta_y < -0.2 then return {walking = true, direction = defines.direction.northwest}
		else                       return {walking = true, direction = defines.direction.west}
		end
	else
		-- Vertically
		if     delta_y >  0.2 then return {walking = true,  direction = defines.direction.south}
		elseif delta_y < -0.2 then return {walking = true,  direction = defines.direction.north}
		else                       return {walking = false, direction = defines.direction.north}
		end
	end
end

commands.add_command("walkto", nil, function(command) end)

commands.add_command("writeworld", nil, function(command)
	write_world()
end)

script.on_load(function()
end)

script.on_event(defines.events.on_script_path_request_finished, function(event) 
	p.print(game.table_to_json(event))
end)

script.on_event(defines.events.on_game_created_from_scenario, function()
	remote.call("freeplay", "set_skip_intro", true)
	speed(1)
end)

script.on_event(defines.events.on_tick, function(event)
	p = game.players[1];
	surface = game.surfaces[1];

	if walking_state.walking then
		p.walking_state = walking_state
		walking_state = getdir(walking_to.x, walking_to.y)
	end
end)
