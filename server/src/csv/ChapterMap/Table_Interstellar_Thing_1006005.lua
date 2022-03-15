local lua_values = 
 { 
    id = {1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,},
    x = {14,12,3,2,9,8,12,13,10,11,9,10,2,1,2,2,},
    y = {2,11,3,12,7,10,12,13,6,6,10,10,2,3,13,14,},
    icon = {'0','ins_monster_01','ins_monster_01','ins_monster_01','ins_monster_02','ins_monster_02','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01',},
    changeicon = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    name = {'0','敌人','敌人','敌人','精锐敌人','精锐敌人','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱',},
    type = {1,3,3,3,4,4,6,6,6,6,6,6,6,6,6,6,},
    initialtype = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    initial = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    delayed = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    trigger = {0,0,0,0,0,0,1,1,1,1,1,1,1,1,1,1,},
    passby = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    show = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    event = {0,500633,500634,500635,500639,500640,605101,605102,605103,605104,605105,605106,605107,605108,605109,605110,},
    priority = {1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,},
    story = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    journal = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    tips = {'0','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力强劲的敌人占据。\n击败他们才能通过此处。','此处被一队实力强劲的敌人占据。\n击败他们才能通过此处。','0','0','0','0','0','0','0','0','0','0',},
    remarks = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    removetype = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    removerelation = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    remove = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    establish = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    modify = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    dispel = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    performance = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    setup = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
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
};

local LuaData = {lua_values = lua_values,lua_idkey = lua_idkey};

function LuaData.GetIds() 
    return {1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,};
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