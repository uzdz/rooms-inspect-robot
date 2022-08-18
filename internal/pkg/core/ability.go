package core

/**
 * @Description: 不同平台所需能力
 */
type AbilityService interface {

	/**
	 * @Description: 获取平台信息
	 */
	GetPlatform() string

	/**
	 * @Description: 解析房源HTTP请求URL，设置请求URL模版
	 * @param requestUrl: 请求地址
	 */
	Validation()

	/**
	 * @Description: 解析请求总页数
	 * @param requestUrl: 请求地址
	 * @return int: 总页数
	 */
	TotalPage() int

	/**
	 * @Description: 获取指定分页最新房源
	 * @param page 总页数
	 * @return []core.Room 最新房源列表
	 */
	ObtainRefreshRooms(page int) []Room

	/**
	 * @Description: 计算房间差值
	 * @param oldRooms 已缓存房源
	 * @param refreshRooms 最新房源
	 * @param []core.Room 新房源
	 */
	Calculation(refreshRooms []Room) []Room
}
