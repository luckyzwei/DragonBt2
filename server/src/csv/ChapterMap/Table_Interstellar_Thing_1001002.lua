local lua_values = 
 { 
    id = {1,2,3,4,5,6,7,8,9,10,},
    x = {9,7,4,8,9,10,3,2,1,9,},
    y = {2,5,6,5,6,6,6,5,5,2,},
    icon = {'0','ins_monster_01','ins_monster_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','0',},
    changeicon = {'0','0','0','0','0','0','0','0','0','0',},
    name = {'0','敌人','敌人','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','剧情',},
    type = {1,3,3,6,6,6,6,6,6,2,},
    initialtype = {0,0,0,0,0,0,0,0,0,0,},
    initial = {'0','0','0','0','0','0','0','0','0','0',},
    delayed = {0,0,0,0,0,0,0,0,0,0,},
    trigger = {0,0,0,1,1,1,1,1,1,1,},
    passby = {0,0,0,0,0,0,0,0,0,1,},
    show = {0,0,0,0,0,0,0,0,0,0,},
    event = {0,500107,500108,102101,102102,102103,102104,102105,102106,0,},
    priority = {1,1,1,1,1,1,1,1,1,2,},
    story = {0,0,0,0,0,0,0,0,0,0,},
    journal = {'0','0','0','0','0','0','0','0','0','击败所有守卫获取奖励！',},
    tips = {'0','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','0','0','0','0','0','0','0',},
    remarks = {'0','0','0','0','0','0','0','0','0','0',},
    removetype = {0,0,0,0,0,0,0,0,0,0,},
    removerelation = {'0','0','0','0','0','0','0','0','0','0',},
    remove = {'0','0','0','0','0','0','0','0','0','0',},
    establish = {'0','0','0','0','0','0','0','0','0','0',},
    modify = {'0','0','0','0','0','0','0','0','0','0',},
    dispel = {'0','0','0','0','0','0','0','0','0','0',},
    performance = {0,0,0,0,0,0,0,0,0,0,},
    setup = {0,0,0,0,0,0,0,0,0,0,},
}; 

local lua_idkey = 
 { 
    [1] = 1,
    [2] = 2,
    [3] = 3,
    [4] = 4,
    [5] = 5,
    [6] = 6,
    [7] = 7,
    [8] = 8,
    [9] = 9,
    [10] = 10,
};

local LuaData = {lua_values = lua_values,lua_idkey = lua_idkey};

function LuaData.GetIds() 
    return {1,2,3,4,5,6,7,8,9,10,};
end

function LuaData.GetKeys() 
    return {'id','x','y','icon','changeicon','name','type','initialtype','initial','delayed','trigger','passby','show','event','priority','story','journal','tips','remarks','removetype','removerelation','remove','establish','modify','dispel','performance','setup',};
end

function LuaData.GetIndex(id) 
   local index = lua_idkey[id];
   if (index == nil) then
       return nil;
   end
   return index
end

function LuaData.GetValueTable(id) 
    if (id == nil) then 
        return nil;
    end
    local index = lua_idkey[id];
    if (index == nil) then
        return nil;
    end
    return {id = lua_values.id[index], x = lua_values.x[index], y = lua_values.y[index], icon = lua_values.icon[index], changeicon = lua_values.changeicon[index], name = lua_values.name[index], type = lua_values.type[index], initialtype = lua_values.initialtype[index], initial = lua_values.initial[index], delayed = lua_values.delayed[index], trigger = lua_values.trigger[index], passby = lua_values.passby[index], show = lua_values.show[index], event = lua_values.event[index], priority = lua_values.priority[index], story = lua_values.story[index], journal = lua_values.journal[index], tips = lua_values.tips[index], remarks = lua_values.remarks[index], removetype = lua_values.removetype[index], removerelation = lua_values.removerelation[index], remove = lua_values.remove[index], establish = lua_values.establish[index], modify = lua_values.modify[index], dispel = lua_values.dispel[index], performance = lua_values.performance[index], setup = lua_values.setup[index], }
end

function LuaData.GetValue(id, key) 
   if (id == nil) then
       return nil;
   end
   local index = lua_idkey[id];
   if (index == nil) then
       return nil;
   end
   if (lua_values[key] == nil) then
       return nil;
   end
   return lua_values[key][index];
end

function LuaData.GetColValues(key)
   if (lua_values[key] == nil) then
       return nil;
   end
   return lua_values[key];
end

function LuaData.IsIdExist(id)
   if (id == nil) then
       return false;
   end
   local index = lua_idkey[id];
   if (index == nil) then
       return false;
   end
   return true;
end

function LuaData.IsKeyExist(key)
   if (lua_values[key] == nil) then
       return false;
   end
   return true;
end

return LuaData;