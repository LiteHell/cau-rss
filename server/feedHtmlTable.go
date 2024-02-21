package server

type HTMLCell struct {
	IsFeedLinks bool
	Colspan     int
	Rowspan     int
	FeedKey     string
	Text        string
}

type HTMLRow struct {
	Cells []HTMLCell
}

func processSiteForFeedHtmlTable(site *CauWebsite, rows *[]HTMLRow) {
	var row HTMLRow
	var childRows []HTMLRow

	if len(site.Children) > 0 {
		parentHasKey := site.Key != ""

		removeFirstChild := false
		firstChild := *site
		rowspan := len(site.Children)
		if !parentHasKey {
			firstChild = site.Children[0]
			removeFirstChild = true
		} else {
			rowspan++
		}

		row.Cells = []HTMLCell{
			{
				IsFeedLinks: false,
				Rowspan:     rowspan,
				Text:        site.Name,
			},
			{
				Text: firstChild.Name,
			},
			{
				IsFeedLinks: true,
				FeedKey:     firstChild.Key,
			},
		}

		if removeFirstChild {
			site.Children = site.Children[1:]
		}

		for _, childSite := range site.Children {
			if childSite.Hidden {
				rowspan--
				continue
			}

			childRow := HTMLRow{
				Cells: []HTMLCell{
					{
						Text: childSite.Name,
					},
					{
						IsFeedLinks: true,
						FeedKey:     childSite.Key,
					},
				},
			}

			childRows = append(childRows, childRow)
		}

		row.Cells[0].Rowspan = rowspan
	} else {
		row.Cells = []HTMLCell{
			{
				Colspan: 2,
				Text:    site.Name,
			},
			{
				IsFeedLinks: true,
				FeedKey:     site.Key,
			},
		}
	}

	*rows = append(*rows, row)
	*rows = append(*rows, childRows...)
}

func GetFeedHtmlTable() []HTMLRow {
	rows := []HTMLRow{}
	sites := GetSupportedSites()

	for _, site := range sites {
		if site.Hidden || len(site.Children) > 0 {
			continue
		}
		processSiteForFeedHtmlTable(&site, &rows)
	}

	for _, site := range sites {
		if site.Hidden || len(site.Children) == 0 {
			continue
		}
		processSiteForFeedHtmlTable(&site, &rows)
	}
	return rows
}
