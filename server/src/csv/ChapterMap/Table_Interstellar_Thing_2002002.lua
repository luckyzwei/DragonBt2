local lua_values = 
 { 
    id = {1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,},
    x = {6,13,13,10,10,8,8,14,13,14,13,13,12,13,12,5,5,4,4,},
    y = {8,7,10,2,14,1,15,5,6,11,11,2,2,15,14,2,15,3,13,},
    icon = {'0','ins_monster_01','ins_monster_01','ins_monster_02','ins_monster_02','ins_monster_03','ins_monster_03','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_02','ins_case_02',},
    changeicon = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    name = {'0','敌人','敌人','精锐敌人','精锐敌人','首领级敌人','首领级敌人','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','高级宝箱','高级宝箱',},
    type = {1,3,3,4,4,5,5,6,6,6,6,6,6,6,6,6,6,7,7,},
    initialtype = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    initial = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    delayed = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    trigger = {0,0,0,0,0,0,0,1,1,1,1,1,1,1,1,1,1,1,1,},
    passby = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    show = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    event = {0,600213,600214,600219,600220,600223,600224,2002101,2002102,2002103,2002104,2002105,2002106,2002107,2002108,2002109,2002110,2002201,2002202,},
    priority = {1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,},
    story = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    journal = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    tips = {'0','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力强劲的敌人占据。\n击败他们才能通过此处。','此处被一队实力强劲的敌人占据。\n击败他们才能通过此处。','此处被一队实力极其强大的敌人占据。\n击败他们才能通过此处。','此处被一队实力极其强大的敌人占据。\n击败他们才能通过此处。','0','0','0','0','0','0','0','0','0','0','0','0',},
    remarks = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    removetype = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    removerelation = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    remove = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    establish = {'0','0','0','0','0','9','10','0','0','0','0','0','0','0','0','0','0','0','0',},
    modify = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    dispel = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    performance = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    setup = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
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
    [11] = 11,
    [12] = 12,
    [13] = 13,
    [14] = 14,
    [15] = 15,
    [16] = 16,
    [17] = 17,
    [18] = 18,
    [19] = 19,
};

local LuaData = {lua_values = lua_values,lua_idkey = lua_idkey};

function LuaData.GetIds() 
    return {1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,};
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