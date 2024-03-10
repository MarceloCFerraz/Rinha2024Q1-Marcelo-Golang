package models

type Client struct {
	// undercase = private
	// uppercase = public
	Id     int   `json:"id"`
	Saldo  int32 `json:"saldo"`
	Limite int32 `json:"limite"`
}

func (c *Client) IsInvalid() bool {
	// validIds := [5]int{1, 2, 3, 4, 5}
	return (c.Id < 1 || c.Id > 5)
}
