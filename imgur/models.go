package imgur

type Image struct {
	Id   string `json:"id"`
	Link string `json:"link"`
	Nsfw bool   `json:"nsfw"`
}

type Album struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	ImagesCount int      `json:"images_count"`
	Images      []*Image `json:"images"`
}

type Gallery struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	ImagesCount int      `json:"images_count"`
	Images      []*Image `json:"images"`
}

type ImgurResponse struct {
	Data interface{} `json:"data"`
}
