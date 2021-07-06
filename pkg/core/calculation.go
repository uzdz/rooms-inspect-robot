package core

// 将本次循环所有的房间进行计算清洗
func Calculation(rooms []Room) {

	testCount := TestForNotice

	if rooms == nil || len(rooms) <= 0 {
		return
	}

	begin := false
	if LastTimeData == nil || len(LastTimeData) <= 0 {
		begin = true
	}

	for i := 0; i < len(rooms); i++ {
		if begin == false {
			if _, ok := LastTimeData[rooms[i].Url]; ok {
				// 存在，不预警（测试除外）
				if testCount > 0 {
					DingToInfo(rooms[i])
					testCount = testCount - 1
				}
			} else {
				// 通知dingding
				DingToInfo(rooms[i])
			}
		}

		LastTimeData[rooms[i].Url] = rooms[i]
	}
}
