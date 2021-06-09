function has_value(tab, val) -- check if key is in a table
	for index, value in ipairs(tab) do
		if value == val then return true end
	end
	return false
end

function get_dir(x, y) -- returns the direction from player position to {x. y}
	local delta_x = x - p.position.x
	local delta_y = y - p.position.y
	if delta_x > 0.2 then
		if     delta_y >  0.2 then return {walking = true, direction = defines.direction.southeast}
		elseif delta_y < -0.2 then return {walking = true, direction = defines.direction.northeast}
		else                       return {walking = true, direction = defines.direction.east}
		end
	elseif delta_x < -0.2 then
		if     delta_y >  0.2 then return {walking = true, direction = defines.direction.southwest}
		elseif delta_y < -0.2 then return {walking = true, direction = defines.direction.northwest}
		else                       return {walking = true, direction = defines.direction.west}
		end
	else
		if     delta_y >  0.2 then return {walking = true,  direction = defines.direction.south}
		elseif delta_y < -0.2 then return {walking = true,  direction = defines.direction.north}
		else                       return {walking = false, direction = defines.direction.north}
		end
	end
end
