package core

/**
 * @Description: 发送通知
 */
type NoticeService interface {

	GetName() string

	/**
	 * @Description: 发送消息
	 */
	Send(room Room, url, key string) bool
}