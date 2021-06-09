require("util")
require("helper")

walking_state = {walking = false} --  what direction we are walking and if were walking
walking_to = {} -- the position that were walking to
p = {}; -- player
surface = {};

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

function write_world()
	-- == water == --
	local tiles = surface.find_tiles_filtered{name = {"water", "deepwater"}}
	local output = {}
	for i, t in pairs(tiles) do
		output[i] = { t.position.x, t.position.y }
	end
	game.write_file("water.json", game.table_to_json(output))

	-- == resources == --
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
	game.write_file("resources.json", game.table_to_json(output))
	game.print("world export done")
end

--- ================== ---
--- ===  COMMANDS  === ---
--- ================== ---

commands.add_command("walkto", nil, function(command)
	local args = game.json_to_table(command.parameter)
	if args == nil then -- check that the argument is valid
		game.print("wrong input: " .. command.parameter)
		return
	end
	walking_state = get_dir(args.x, args.y) -- this sets the walking state to true
	walking_to = args
end)

commands.add_command("writeworld", nil, function(command)
	write_world()
end)

commands.add_command("drawbox", nil, function(command)
	local a = game.json_to_table(command.parameter)
	if a == nil then -- check that the argument is valid
		game.print("wrong input: " .. command.parameter)
		return
	end
	mark_area(a.color, nil, a.x1, a.y1, a.x2, a.y2)
end)

--- ================ ---
--- ===  EVENTS  === ---
--- ================ ---

script.on_event(defines.events.on_tick, function(event)
	p = game.players[1];
	surface = game.get_surface("nauvis")

	if walking_state.walking then
		p.walking_state = walking_state
		walking_state = get_dir(walking_to.x, walking_to.y)
		if not walking_state.walking then game.print("walking done") end
	end
end)
