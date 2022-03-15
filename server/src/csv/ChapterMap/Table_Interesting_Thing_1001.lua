local lua_values = 
 { 
    id = {1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,},
    x = {1,1,1,3,9,7,1,3,9,7,2,9,8,2,9,8,10,2,4,7,10,8,8,1,10,10,2,2,9,3,4,9,9,},
    y = {1,1,2,5,8,8,2,5,8,8,8,3,9,8,3,9,9,9,7,5,2,8,10,10,1,10,4,5,9,6,5,2,10,},
    icon = {'0','0','ins_monster_01','ins_monster_01','ins_monster_01','ins_monster_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_monster_02','ins_monster_02','ins_monster_02','ins_relics_01','ins_relics_01','ins_relics_01','ins_monster_03','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_02','ins_case_02','ins_case_02','ins_barracks_01','ins_fountain_01','ins_resurrection_01','0','0','ins_barracks_01','ins_castle_01',},
    changeicon = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    name = {'0','剧情','敌人','敌人','敌人','敌人','选择遗物','选择遗物','选择遗物','选择遗物','精锐敌人','精锐敌人','精锐敌人','选择遗物','选择遗物','选择遗物','首领级敌人','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','高级宝箱','高级宝箱','高级宝箱','佣兵营地','魔法泉水','复苏祭坛','剧情','剧情','佣兵营地','城堡',},
    type = {1,2,3,3,3,3,8,8,8,8,4,4,4,8,8,8,5,6,6,6,6,6,6,7,7,7,9,10,11,16,16,9,12,},
    initialtype = {0,0,0,0,0,0,2,2,2,2,0,0,0,2,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    initial = {'0','0','0','0','0','0','3','4','5','6','0','0','0','11','12','13','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    delayed = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    trigger = {0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,1,1,1,1,1,1,1,1,0,0,0,1,1,0,0,},
    passby = {0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,1,0,0,},
    show = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,},
    event = {0,0,400101,400102,400104,400103,2,2,2,2,400107,400106,400105,3,3,3,400108,1001101,1001102,1001103,1001104,1001105,1001106,1001201,1001202,1001203,0,0,0,0,0,0,1,},
    priority = {1,2,1,1,1,1,2,2,2,2,1,1,1,2,2,2,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,},
    story = {0,101001,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    journal = {'0','边境突然出现大量可疑的暗影生物，这或许与神器失踪有关系，先尝试击败他们搜寻线索。','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','前面的暗影生物数量变多，实力也增强了不少，一定是在看守着什么重要的东西。','前面的暗影生物数量变多，实力也增强了不少，一定是在看守着什么重要的东西。','0','0',},
    tips = {'0','0','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','0','0','0','0','此处被一队实力强劲的敌人占据。\n击败他们才能通过此处。','此处被一队实力强劲的敌人占据。\n击败他们才能通过此处。','此处被一队实力强劲的敌人占据。\n击败他们才能通过此处。','0','0','0','此处被一队实力极其强大的敌人占据。\n击败他们才能通过此处。','0','0','0','0','0','0','0','0','0','邀请一位英雄，协助完成本次探索。','聚集魔法能量的温泉，似乎可以缓解战斗带来的伤痛。','神秘的祭坛，似乎有生命的法则环绕。','0','0','邀请一位英雄，协助完成本次探索。','0',},
    remarks = {'0','0','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','0','0','0','0','掉落精英或史诗级别的遗物','掉落精英或史诗级别的遗物','掉落精英或史诗级别的遗物','0','0','0','0','0','0','0','0','0','0','0','0','0','0','所有存活英雄恢复50%生命','随机复活一名已经阵亡的英雄','0','0','0','0',},
    removetype = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    removerelation = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    remove = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','31','30','0','0',},
    establish = {'0','0','7','8','9','10','0','0','0','0','14','15','16','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    modify = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    dispel = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','7|8=8|9=10|9=9|8','7|8=8|9=10|9=9|8','0','0',},
    performance = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,10101,10101,0,0,},
    setup = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
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
    [20] = 20,
    [21] = 21,
    [22] = 22,
    [23] = 23,
    [24] = 24,
    [25] = 25,
    [26] = 26,
    [27] = 27,
    [28] = 28,
    [29] = 29,
    [30] = 30,
    [31] = 31,
    [32] = 32,
    [33] = 33,
};

local LuaData = {lua_values = lua_values,lua_idkey = lua_idkey};

function LuaData.GetIds() 
    return {1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,};
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