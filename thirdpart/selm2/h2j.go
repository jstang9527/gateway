package selm2

// H2j ...
type H2j struct {
	Name     string            `json:"name,omitempty"` // 对应 HTML 标签
	Type     string            `json:"type,omitempty"` // element 或者 text
	Text     string            `json:"text,omitempty"`
	Attrs    map[string]string `json:"attrs,omitempty"`
	Children []*H2j            `json:"children,omitempty"`
}
