package pager

/*

This package is Go implemntation of Perl5's Data::Page and Data::Page::Navigation
see. https://metacpan.org/pod/Data::Page

*/

import (
	"net/url"
	"strconv"
)

type Pager struct {
	totalEntries       int
	entriesPerPage     int
	currentPage        int
	pagesPerNavigation int
	url                *url.URL
}

func New(args ...Option) Pager {
	opt := options{
		totalEntries:       0,
		entriesPerPage:     10,
		currentPage:        1,
		pagesPerNavigation: 10,
	}
	for _, o := range args {
		o(&opt)
	}

	return Pager{
		totalEntries:       opt.totalEntries,
		entriesPerPage:     opt.entriesPerPage,
		currentPage:        opt.currentPage,
		pagesPerNavigation: opt.pagesPerNavigation,
		url:                opt.url,
	}
}

func (p Pager) SetTotalEntries(te int) {
	if te >= 0 {
		p.totalEntries = te
	}
}

func (p Pager) SetEntriesPerPage(epp int) {
	if epp > 0 {
		p.entriesPerPage = epp
	}
}

func (p Pager) SetCurrentPage(cp int) {
	if cp > 0 {
		p.currentPage = cp
	}
}

func (p Pager) TotalEntries() int {
	return p.totalEntries
}

func (p Pager) EntriesPerPage() int {
	return p.entriesPerPage
}

func (p Pager) CurrentPage() int {
	if p.currentPage < p.FirstPage() {
		return p.FirstPage()
	}
	if p.currentPage > p.LastPage() {
		return p.LastPage()
	}
	return p.currentPage
}

func (p Pager) EntriesOnThisPage() int {
	if p.totalEntries == 0 {
		return 0
	}
	return p.Last() - p.First() + 1
}

func (p Pager) FirstPage() int {
	return 1
}

func (p Pager) LastPage() int {
	return 1 + (p.totalEntries / p.entriesPerPage)
}

func (p Pager) First() int {
	if p.totalEntries == 0 {
		return 0
	}
	return ((p.currentPage - 1) * p.entriesPerPage) + 1
}

func (p Pager) Last() int {
	if p.currentPage == p.LastPage() {
		return p.totalEntries
	}
	return p.currentPage * p.entriesPerPage
}

func (p Pager) PreviousPage() int {
	if p.currentPage > 1 {
		return p.currentPage - 1
	}
	return 0
}

func (p Pager) NextPage() int {
	if p.currentPage < p.LastPage() {
		return p.currentPage + 1
	}
	return 0
}

// FIXME
func (p Pager) Select(list []interface{}) []interface{} {
	top := len(list)
	if top > p.Last() {
		top = p.Last()
	}
	if top == 0 {
		return list[:0]
	}
	return list[p.First()-1 : top-1]
}

func (p Pager) Skipped() int {
	skipped := p.First() - 1
	if skipped < 0 {
		return 0
	}
	return skipped
}

func (p Pager) ChangeEntriesPerPage(epp int) {
	if epp < 1 {
		return
	}
	cp := 1 + (p.First() / epp)
	p.entriesPerPage = epp
	p.currentPage = cp
}

func (p Pager) PagesInNavigation() []int {
	lp := p.LastPage()
	ppn := p.pagesPerNavigation
	if ppn >= lp {
		ret := make([]int, lp)
		for i := 0; i < lp; i++ {
			ret[i] = p.FirstPage() + i
		}
		return ret
	}

	prev := p.PreviousPage()
	next := p.NextPage()
	ret := []int{}
	i := 0
	for {
		if len(ret) < ppn {
			break
		}
		if i%2 != 0 {
			if p.FirstPage() <= prev {
				ret = append([]int{prev}, ret...)
			}
			prev--
		} else {
			if lp >= next {
				ret = append(ret, next)
				next++
			}
		}
		i++
	}

	return ret
}

func (p Pager) PageLink(i int) string {
	if p.url == nil {
		return ""
	}
	q, err := url.ParseQuery(p.url.RawQuery)
	if err != nil {
		return ""
	}
	q.Set("p", strconv.Itoa(i))
	return p.url.EscapedPath() + "?" + q.Encode()
}

//
// options
//

type options struct {
	totalEntries       int
	entriesPerPage     int
	currentPage        int
	pagesPerNavigation int
	url                *url.URL
}

type Option func(*options)

func WithTotalEntries(te int) Option {
	return func(opts *options) {
		if te >= 0 {
			opts.totalEntries = te
		}
	}
}

func WithEntriesPerPage(epp int) Option {
	return func(opts *options) {
		if epp > 0 {
			opts.entriesPerPage = epp
		}
	}
}

func WithCurrentPage(cp int) Option {
	return func(opts *options) {
		if cp > 0 {
			opts.currentPage = cp
		}
	}
}

func WithPagesPerNavigation(ppn int) Option {
	return func(opts *options) {
		if ppn > 0 {
			opts.pagesPerNavigation = ppn
		}
	}
}

func WithURL(u *url.URL) Option {
	return func(opts *options) {
		if u != nil {
			opts.url = u
		}
	}
}
