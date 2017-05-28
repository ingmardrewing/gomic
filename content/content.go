package content

/**
 * struct for the comic page
 */

type Page struct {
	Id, PageNumber                                  int
	Title, Description, Path, ImgUrl, DisqusId, Act string
}

type Pages struct {
	Pages []Page
}

func EmptyPage() Page {
	return Page{0, 0, "", "", "", "", "", ""}
}
