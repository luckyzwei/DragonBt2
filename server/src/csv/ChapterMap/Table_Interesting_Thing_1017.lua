local lua_values = 
 { 
    id = {1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100,101,102,103,104,105,106,107,108,109,110,111,112,113,114,115,116,117,118,119,120,121,122,},
    x = {1,1,4,7,8,14,15,15,4,7,8,14,15,15,2,2,3,4,5,6,7,8,10,12,12,13,14,15,17,2,2,3,4,5,6,7,8,10,12,12,13,14,15,17,9,11,13,18,9,11,13,18,1,1,1,2,4,5,5,6,7,8,9,10,10,10,10,11,11,12,13,13,13,13,16,18,19,19,8,14,20,5,9,12,14,1,11,1,1,2,1,12,20,1,1,2,1,12,20,6,10,12,15,3,3,1,15,5,14,7,13,7,2,6,8,6,3,4,16,18,2,18,},
    y = {13,13,13,19,21,20,4,7,13,19,21,20,4,7,10,15,21,10,6,21,4,5,11,9,16,4,23,19,4,10,15,21,10,6,21,4,5,11,9,16,4,23,19,4,14,14,12,21,14,14,12,21,4,9,17,22,5,10,24,1,24,1,16,4,5,13,22,11,16,3,3,10,13,22,17,2,3,20,15,13,20,8,13,22,8,14,12,8,22,3,18,24,4,8,22,3,18,24,4,15,18,7,17,12,13,12,2,14,19,8,6,2,9,6,19,24,5,20,21,20,16,5,},
    icon = {'0','0','ins_monster_01','ins_monster_01','ins_monster_01','ins_monster_01','ins_monster_01','ins_monster_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_monster_02','ins_monster_02','ins_monster_02','ins_monster_02','ins_monster_02','ins_monster_02','ins_monster_02','ins_monster_02','ins_monster_02','ins_monster_02','ins_monster_02','ins_monster_02','ins_monster_02','ins_monster_02','ins_monster_02','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_monster_03','ins_monster_03','ins_monster_03','ins_monster_03','ins_relics_01','ins_relics_01','ins_relics_01','ins_relics_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_01','ins_case_02','ins_case_02','ins_case_02','ins_barracks_01','ins_fountain_01','ins_fountain_01','ins_fountain_01','ins_resurrection_01','ins_resurrection_01','ins_mechanism_03','ins_mechanism_03','ins_mechanism_03','ins_mechanism_01','ins_mechanism_01','ins_mechanism_01','ins_mechanism_99','ins_mechanism_99','ins_mechanism_99','ins_mechanism_99','ins_mechanism_99','ins_mechanism_99','ins_portal_01','ins_portal_02','ins_portal_03','ins_portal_04','ins_portal_05','ins_portal_06','ins_sun_console_03','ins_sun_console_03','ins_sun_decline_03','ins_sun_decline_03','ins_moon_decline_03','ins_moon_decline_03','ins_sun_console_01','ins_sun_decline_01','ins_sun_decline_01','ins_moon_decline_01','ins_sun_console_02','ins_sun_decline_02','ins_moon_decline_02','ins_moon_decline_02','ins_sun_console_04','ins_sun_decline_04','ins_moon_decline_04',},
    changeicon = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','ins_moon_console_03','ins_moon_console_03','ins_sun_rise_03','ins_sun_rise_03','ins_moon_rise_03','ins_moon_rise_03','ins_moon_console_01','ins_sun_rise_01','ins_sun_rise_01','ins_moon_rise_01','ins_moon_console_02','ins_sun_rise_02','ins_moon_rise_02','ins_moon_rise_02','ins_moon_console_04','ins_sun_rise_04','ins_moon_rise_04',},
    name = {'0','剧情','敌人','敌人','敌人','敌人','敌人','敌人','选择遗物','选择遗物','选择遗物','选择遗物','选择遗物','选择遗物','精锐敌人','精锐敌人','精锐敌人','精锐敌人','精锐敌人','精锐敌人','精锐敌人','精锐敌人','精锐敌人','精锐敌人','精锐敌人','精锐敌人','精锐敌人','精锐敌人','精锐敌人','选择遗物','选择遗物','选择遗物','选择遗物','选择遗物','选择遗物','选择遗物','选择遗物','选择遗物','选择遗物','选择遗物','选择遗物','选择遗物','选择遗物','选择遗物','首领级敌人','首领级敌人','首领级敌人','首领级敌人','选择遗物','选择遗物','选择遗物','选择遗物','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','普通宝箱','高级宝箱','高级宝箱','高级宝箱','佣兵营地','魔法泉水','魔法泉水','魔法泉水','复苏祭坛','复苏祭坛','蓝机关','蓝机关','蓝机关','红机关','红机关','红机关','踩下机关','踩下机关','踩下机关','踩下机关','踩下机关','踩下机关','红传送点','绿传送点','蓝传送点','黄传送点','青传送点','紫传送点','蓝控制台','蓝控制台','蓝太阳','蓝太阳','蓝月亮','蓝月亮','红控制台','红太阳','红太阳','红月亮','绿控制台','绿太阳','绿月亮','绿月亮','黄控制台','黄太阳','黄月亮',},
    type = {1,2,3,3,3,3,3,3,8,8,8,8,8,8,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,8,8,8,8,8,8,8,8,8,8,8,8,8,8,8,5,5,5,5,8,8,8,8,6,6,6,6,6,6,6,6,6,6,6,6,6,6,6,6,6,6,6,6,6,6,6,6,6,6,7,7,7,9,10,10,10,11,11,16,16,16,16,16,16,19,19,19,19,19,19,17,17,17,17,17,17,20,20,100,100,101,101,20,100,100,101,20,100,101,101,20,100,101,},
    initialtype = {0,0,0,0,0,0,0,0,2,2,2,2,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,0,0,0,0,2,2,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,1,1,1,1,1,0,0,0,0,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    initial = {'0','0','0','0','0','0','0','0','3','4','5','6','7','8','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','15','16','17','18','19','20','21','22','23','24','25','26','27','28','29','0','0','0','0','45','46','47','48','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','88','89','90','91','92','93','0','0','0','0','88|89|90','91|92|93','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    delayed = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    trigger = {0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,0,0,0,0,0,0,1,1,1,1,1,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    passby = {0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,1,1,1,1,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    show = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,1,1,1,1,1,0,0,0,0,0,0,1,1,1,1,0,0,0,0,0,0,1,1,1,0,0,1,1,0,1,1,0,0,1,},
    event = {0,0,401001,401001,401001,401001,401001,401001,1,1,1,1,1,1,401001,401001,401001,401001,401001,401001,401001,401001,401001,401001,401001,401001,401001,401001,401001,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,401001,401001,401001,401001,3,3,3,3,1017101,1017102,1017103,1017104,1017105,1017106,1017107,1017108,1017109,1017110,1017111,1017112,1017113,1017114,1017115,1017116,1017117,1017118,1017119,1017120,1017121,1017122,1017123,1017124,1017125,1017126,1017201,1017202,1017203,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,1,1,0,1,1,0,0,1,1,0,0,1,},
    priority = {1,2,1,1,1,1,1,1,2,2,2,2,2,2,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,1,1,1,1,2,2,2,2,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,2,2,2,2,2,2,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,},
    story = {0,1701001,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    journal = {'0','这里被黑暗和魔法的力量扭曲，形成互相影响的区域，要小心那些机关！','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    tips = {'0','0','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','0','0','0','0','0','0','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','此处被一队实力普通的敌人占据。\n击败他们才能通过此处。','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','邀请一位英雄，协助完成本次探索。','聚集魔法能量的温泉，似乎可以缓解战斗带来的伤痛。','聚集魔法能量的温泉，似乎可以缓解战斗带来的伤痛。','神秘的祭坛，似乎有生命的法则环绕。','神秘的祭坛，似乎有生命的法则环绕。','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','一个魔法机关，可以改变所有蓝色石柱的状态。','一个魔法机关，可以改变所有蓝色石柱的状态。','一个地面升起的石柱，上面有“太阳”的图案。','一个地面升起的石柱，上面有“太阳”的图案。','一个地面升起的石柱，上面有“月亮”的图案。','一个地面升起的石柱，上面有“月亮”的图案。','一个魔法机关，可以改变所有红色石柱的状态。','一个地面升起的石柱，上面有“太阳”的图案。','一个地面升起的石柱，上面有“太阳”的图案。','一个地面升起的石柱，上面有“太阳”的图案。','一个魔法机关，可以改变所有绿色石柱的状态。','一个地面升起的石柱，上面有“太阳”的图案。','一个地面升起的石柱，上面有“太阳”的图案。','一个地面升起的石柱，上面有“太阳”的图案。','一个魔法机关，可以改变所有黄色石柱的状态。','一个地面升起的石柱，上面有“月亮”的图案。','一个地面升起的石柱，上面有“月亮”的图案。',},
    remarks = {'0','0','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','0','0','0','0','0','0','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','掉落稀有或精英级别的遗物','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','所有存活英雄恢复50%生命','所有存活英雄恢复50%生命','随机复活一名已经阵亡的英雄','随机复活一名已经阵亡的英雄','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    removetype = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    removerelation = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    remove = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    establish = {'0','0','9','10','11','12','13','14','0','0','0','0','0','0','30','31','32','33','34','35','36','37','38','39','40','41','42','43','44','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','49','50','51','52','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','94|104','95|104','96|104','97|105','98|105','99|105','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    modify = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','12|18','5|12','7|17','15|9','8|10','14|15','107|108|109|110|111','106|108|109|110|111','106|107|109|110|111','106|107|108|110|111','106|107|108|109|111','106|107|108|109|110','113|114|115','112|114|115','112|113|115','112|113|114','117|118|119','116|118|119','116|117|119','116|117|118','121|122','120|122','120|121',},
    dispel = {'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0',},
    performance = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
    setup = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,170101,170201,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,},
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
    [34] = 34,
    [35] = 35,
    [36] = 36,
    [37] = 37,
    [38] = 38,
    [39] = 39,
    [40] = 40,
    [41] = 41,
    [42] = 42,
    [43] = 43,
    [44] = 44,
    [45] = 45,
    [46] = 46,
    [47] = 47,
    [48] = 48,
    [49] = 49,
    [50] = 50,
    [51] = 51,
    [52] = 52,
    [53] = 53,
    [54] = 54,
    [55] = 55,
    [56] = 56,
    [57] = 57,
    [58] = 58,
    [59] = 59,
    [60] = 60,
    [61] = 61,
    [62] = 62,
    [63] = 63,
    [64] = 64,
    [65] = 65,
    [66] = 66,
    [67] = 67,
    [68] = 68,
    [69] = 69,
    [70] = 70,
    [71] = 71,
    [72] = 72,
    [73] = 73,
    [74] = 74,
    [75] = 75,
    [76] = 76,
    [77] = 77,
    [78] = 78,
    [79] = 79,
    [80] = 80,
    [81] = 81,
    [82] = 82,
    [83] = 83,
    [84] = 84,
    [85] = 85,
    [86] = 86,
    [87] = 87,
    [88] = 88,
    [89] = 89,
    [90] = 90,
    [91] = 91,
    [92] = 92,
    [93] = 93,
    [94] = 94,
    [95] = 95,
    [96] = 96,
    [97] = 97,
    [98] = 98,
    [99] = 99,
    [100] = 100,
    [101] = 101,
    [102] = 102,
    [103] = 103,
    [104] = 104,
    [105] = 105,
    [106] = 106,
    [107] = 107,
    [108] = 108,
    [109] = 109,
    [110] = 110,
    [111] = 111,
    [112] = 112,
    [113] = 113,
    [114] = 114,
    [115] = 115,
    [116] = 116,
    [117] = 117,
    [118] = 118,
    [119] = 119,
    [120] = 120,
    [121] = 121,
    [122] = 122,
};

local LuaData = {lua_values = lua_values,lua_idkey = lua_idkey};

function LuaData.GetIds() 
    return {1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100,101,102,103,104,105,106,107,108,109,110,111,112,113,114,115,116,117,118,119,120,121,122,};
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