package test

import (
	"net/http"
	"net/url"
	"testing"
)

func isRedirected(from string, to string) (bool, string, error) {
	baseUrl, err := url.Parse("http://127.0.0.1:8080")

	fromRelativeUrl, err := url.Parse(from)
	if err != nil {
		return false, "", err
	}

	toUrl, err := url.Parse(to)
	if err != nil {
		return false, "", err
	}
	toUrl = baseUrl.ResolveReference(toUrl)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(baseUrl.ResolveReference(fromRelativeUrl).String())
	if err != nil {
		return false, "", err
	}

	isRedirectCode := resp.StatusCode == 301 || resp.StatusCode == 302 || resp.StatusCode == 307 || resp.StatusCode == 308
	if !isRedirectCode {
		return false, "", nil
	}

	locationUrl, err := resp.Location()
	if err != nil {
		return false, resp.Request.URL.String(), err
	}
	return *locationUrl == *toUrl, locationUrl.String(), nil

}

func testRedirect(from, to string, t *testing.T) {
	redirected, location, err := isRedirected("/index.html", "/")

	if err != nil {
		t.Error(err)
	} else if !redirected {
		t.Errorf("Redirect failure (from %s, expected: %s, real: %s)", from, to, location)
	}
}

func TestRedirect(t *testing.T) {
	runServer()

	testRedirect("/index.html", "/", t)
	testRedirect("/cau/notice", "https://www.cau.ac.kr/cms/FR_PRO_CON/BoardRss.do?pageNo=1&pagePerCnt=15&MENU_ID=100&SITE_NO=2&BOARD_SEQ=4&S_CATE_SEQ=&BOARD_TYPE=C0301&BOARD_CATEGORY_NO=&P_TAB_NO=&TAB_NO=&P_CATE_SEQ=&CATE_SEQ=&SEARCH_FLD=SUBJECT&SEARCH=", t)
	for _, feedType := range []string{"rss", "atom", "json"} {
		testRedirect("/cau/sw/"+feedType, "/cau/swedu/"+feedType, t)
		testRedirect("/cau/dormitory/"+feedType, "/cau/dormitory/seoul/all/"+feedType, t)
	}
}
