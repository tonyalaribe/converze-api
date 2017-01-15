package models

import (
	"math/rand"
	"time"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

//Page carries pagination info to aid in knowing whether any given page has a
//next or previous page, and to know its page number
type Page struct {
	Prev    bool
	PrevVal int

	Next    bool
	NextVal int

	NextURL string

	pages int
	Pages []string
	Total int
	Count int
	Skip  int
}

//SearchPagination returns a page strict which carries details about the
//pagination of any given search result or pag/Users/SMILECS/Downloads/go-oddjobs-master/functions.goe
func SearchPagination(count int, page int, perPage int) Page {
	var pg Page
	var total int

	if count%perPage != 0 {
		total = count/perPage + 1
	} else {
		total = count / perPage
	}

	if total < page {
		page = total
	}

	if page == 1 {
		pg.Prev = false
		pg.Next = true
	}

	if page != 1 {
		pg.Prev = true
	}

	if total > page {
		pg.Next = true
	}

	if total == page {
		pg.Next = false
	}

	var pgs = make([]string, total)

	//The number of number of documents to skip
	skip := perPage * (page - 1)

	pg.Total = total
	pg.Skip = skip
	pg.Count = count
	pg.NextVal = page + 1
	pg.PrevVal = page - 1
	pg.Pages = pgs

	return pg
}

// Slug replaces each run of characters which are not unicode letters or
// numbers with a single hyphen, except for leading or trailing runs. Letters
// will be stripped of diacritical marks and lowercased. Letter or number
// codepoints that do not have combining marks or a lower-cased variant will
// be passed through unaltered.
func Slug(s string) string {
	var lat = []*unicode.RangeTable{unicode.Letter, unicode.Number}
	var nop = []*unicode.RangeTable{unicode.Mark, unicode.Sk, unicode.Lm}

	buf := make([]rune, 0, len(s))
	dash := false
	for _, r := range norm.NFKD.String(s) {
		switch {
		// unicode 'letters' like mandarin characters pass through
		case unicode.IsOneOf(lat, r):
			buf = append(buf, unicode.ToLower(r))
			dash = true
		case unicode.IsOneOf(nop, r):
			// skip
		case dash:
			buf = append(buf, '-')
			dash = false
		}
	}
	if i := len(buf) - 1; i >= 0 && buf[i] == '-' {
		buf = buf[:i]
	}
	return string(buf)
}

//RandomString Generates random alphanummeric string of given length
func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
