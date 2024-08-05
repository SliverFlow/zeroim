package session

import (
	"strconv"
	"strings"
)

type Session string // 会话

// NewSession 新建 session
func NewSession(name string, token string, id uint64) Session {
	if len(name) == 0 || len(token) == 0 {
		panic("name or token is empty")
	}

	idStr := strconv.FormatUint(id, 10)
	return Session(name + ":" + token + ":" + idStr)
}

// FormatString 格式化字符串
func FormatString(s string) Session {
	return Session(s)
}

// Name 名称
func (s Session) Name() string {
	arr := strings.Split(string(s), ":")
	if len(arr) != 3 {
		panic("invalid session")
	}

	return arr[0]
}

// Token 令牌
func (s Session) Token() string {
	arr := strings.Split(string(s), ":")
	if len(arr) != 3 {
		panic("invalid session")
	}

	return arr[1]
}

// Id 会话 id
func (s Session) Id() uint64 {
	arr := strings.Split(string(s), ":")
	if len(arr) != 3 {
		panic("invalid session")
	}

	id, err := strconv.ParseUint(arr[2], 10, 64)
	if err != nil {
		panic("invalid id")
	}

	return id
}

// Info 会话信息
func (s Session) Info() (string, string, uint64) {
	arr := strings.Split(string(s), ":")
	if len(arr) != 3 {
		panic("invalid session")
	}

	id, err := strconv.ParseUint(arr[2], 10, 64)
	if err != nil {
		panic("invalid id")
	}

	return arr[0], arr[1], id
}

// String
func (s Session) String() string {
	return string(s)
}
