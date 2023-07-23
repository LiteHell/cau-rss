package server

type CauWebsite struct {
	Hidden   bool
	Name     string
	Url      string
	Key      string
	Children []CauWebsite
}

func GetSupportedSites() []CauWebsite {
	return []CauWebsite{
		{
			Name: "소프트웨어대학/SW교육 관련",
			Url:  "https://sw.cau.ac.kr",
			Children: []CauWebsite{
				{
					Key:  "cse",
					Name: "소프트웨어학부",
					Url:  "https://cse.cau.ac.kr",
				},
				{
					Key:  "ai",
					Name: "AI학과",
					Url:  "https://ai.cau.ac.kr",
				},
				{
					Key:  "swedu",
					Name: "다빈치SW교육원",
					Url:  "https://swedu.cau.ac.kr",
				},
			},
		},
		{
			Name: "공학인증 관련",
			Children: []CauWebsite{
				{
					Key:  "abeek",
					Name: "공학교육혁신센터",
					Url:  "https://abeek.cau.ac.kr",
				},
			},
		},
		{
			Name: "서울캠퍼스 기숙사",
			Children: []CauWebsite{
				{
					Key:  "dormitory/seoul/bluemir",
					Name: "블루미르홀",
					Url:  "https://dormitory.cau.ac.kr",
				},
				{
					Key:  "dormitory/seoul/future_house",
					Name: "퓨처하우스",
					Url:  "https://dormitory.cau.ac.kr",
				},
				{
					Key:  "dormitory/seoul/global_house",
					Name: "글로벌하우스",
					Url:  "https://dormitory.cau.ac.kr",
				},
				{
					Hidden: true,
					Key:    "dormitory/seoul/all",
					Name:   "블루미르홀/퓨처하우스/글로벌하우스",
					Url:    "https://dormitory.cau.ac.kr",
				},
			},
		},
		{
			Name: "다빈치캠퍼스 기숙사",
			Url:  "https://dorm.cau.ac.kr",
			Key:  "dormitory/davinci",
		},
	}
}

func LoopForAllSites(fn func(*CauWebsite)) {
	for _, site := range GetSupportedSites() {
		if site.Key != "" {
			fn(&site)
		}
		if site.Children != nil {
			loopForSites(&site.Children, fn)
		}
	}
}
func loopForSites(sites *[]CauWebsite, fn func(*CauWebsite)) {
	for _, site := range *sites {
		if site.Key != "" {
			fn(&site)
		}
		if site.Children != nil {
			loopForSites(&site.Children, fn)
		}
	}
}