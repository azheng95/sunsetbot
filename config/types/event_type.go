package types

type EventType int

const (
	Evening EventType = iota // 晚霞
	Morning                  // 朝霞
)

func (e EventType) String() string {
	return [...]string{"晚霞", "朝霞"}[e]
}

func (e EventType) Params() string {
	switch e {
	case Evening:
		return "set_1" // 今天日落
	case Morning:
		return "rise_2" // 明天日出
	}
	return ""
}
