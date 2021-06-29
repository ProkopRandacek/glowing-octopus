#!/usr/bin/lua
luna = require 'lunajson'

-- fake data object for the recipe file
data = {table}
-- with fake extend function
function data:extend(table)
	self.table = table
end

-- run the recipes file from factorio. It writes to the data object
require("./factorio.data.base.prototypes.recipe")

recipes = {}
items = {}

fluids = {
    ["water"] = true,
    ["steam"] = true,
    ["sulfuric-acid"] = true,
    ["crude-oil"] = true,
    ["heavy-oil"] = true,
    ["light-oil"] = true,
    ["petroleum-gas"] = true,
    ["lubricant"] = true
}

-- parse the recipes to a better format
for i, t in pairs(data.table) do
	items[i] = t.name

	recipes[t.name] = {}
	if fluids[t.name] == true then
		recipes[t.name].fluid = true
	else
		recipes[t.name].fluid = false
	end
	if t.result_count == nil then t.result_count = 1 end
	if t.energy_required == nil then t.energy_required = 0.5 end
	recipes[t.name].craftTime = t.energy_required / t.result_count

	recipes[t.name].deps = {}

	has_recipe = true
	if t.ingredients == nil then
		t.ingredients = t.normal.ingredients
	end
	for ii, d in pairs(t.ingredients) do
		recipes[t.name].deps[ii] = {}
		recipes[t.name].deps[ii].name = d[1]
		recipes[t.name].deps[ii].count = d[2]
	end
end

local recipesjson = luna.encode(recipes)
local itemsjson = luna.encode(items)

file = io.open("master/recipes.json", "w+")
io.output(file)
io.write(recipesjson)
io.close(file)

file = io.open("master/items.json", "w+")
io.output(file)
io.write(itemsjson)
io.close(file)

print("done")
