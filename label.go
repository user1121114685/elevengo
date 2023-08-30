package elevengo

import (
	"github.com/deadblue/elevengo/internal/api"
	"github.com/deadblue/elevengo/internal/api/errors"
)

type LabelColor int

const (
	LabelColorBlank LabelColor = iota
	LabelColorRed
	LabelColorOrange
	LabelColorYellow
	LabelColorGreen
	LabelColorBlue
	LabelColorPurple
	LabelColorGray
)

var (
	labelColorMap = map[LabelColor]string{
		LabelColorBlank:  api.LabelColorBlank,
		LabelColorRed:    api.LabelColorRed,
		LabelColorOrange: api.LabelColorOrange,
		LabelColorYellow: api.LabelColorYellow,
		LabelColorGreen:  api.LabelColorGreen,
		LabelColorBlue:   api.LabelColorBlue,
		LabelColorPurple: api.LabelColorPurple,
		LabelColorGray:   api.LabelColorGray,
	}

	labelColorRevMap = map[string]LabelColor{
		api.LabelColorBlank:  LabelColorBlank,
		api.LabelColorRed:    LabelColorRed,
		api.LabelColorOrange: LabelColorOrange,
		api.LabelColorYellow: LabelColorYellow,
		api.LabelColorGreen:  LabelColorGreen,
		api.LabelColorBlue:   LabelColorBlue,
		api.LabelColorPurple: LabelColorPurple,
		api.LabelColorGray:   LabelColorGray,
	}
)

type Label struct {
	Id    string
	Name  string
	Color LabelColor
}

// labelIterator implements Iterator[Label].
type labelIterator struct {
	// Offset
	offset int
	// Total count
	count int

	// Cached labels
	labels []*api.LabelInfo
	// Cache index
	index int
	// Cache size
	size int

	// Update function
	uf func(*labelIterator) error
}

func (i *labelIterator) Next() (err error) {
	i.index += 1
	if i.index < i.size {
		return
	}
	i.offset += i.size
	if i.offset >= i.count {
		return errors.ErrReachEnd
	}
	return i.uf(i)
}

func (i *labelIterator) Index() int {
	return i.offset + i.index
}

func (i *labelIterator) Get(label *Label) error {
	if i.index >= i.size {
		return errors.ErrReachEnd
	}
	l := i.labels[i.index]
	label.Id = l.Id
	label.Name = l.Name
	label.Color = labelColorRevMap[l.Color]
	return nil
}

func (i *labelIterator) Count() int {
	return i.count
}

func (a *Agent) LabelIterate() (it Iterator[Label], err error) {
	li := &labelIterator{
		uf: a.labelIterateInternal,
	}
	if err = a.labelIterateInternal(li); err == nil {
		it = li
	}
	return
}

func (a *Agent) labelIterateInternal(i *labelIterator) (err error) {
	spec := (&api.LabelListSpec{}).Init(i.offset)
	if err = a.pc.ExecuteApi(spec); err != nil {
		return
	}
	i.count = spec.Result.Total
	i.index, i.size = 0, len(spec.Result.List)
	i.labels = make([]*api.LabelInfo, i.size)
	copy(i.labels, spec.Result.List)
	return
}

// LabelFind finds label whose name is name, and returns it.
// func (a *Agent) LabelFind(name string, label *Label) (err error) {
// 	qs := protocol.Params{}.
// 		With("keyword", name).
// 		WithInt("limit", 10)
// 	resp := &webapi.BasicResponse{}
// 	if err = a.pc.CallJsonApi(webapi.ApiLabelList, qs, nil, resp); err != nil {
// 		return
// 	}
// 	data := &webapi.LabelListData{}
// 	if err = resp.Decode(data); err != nil {
// 		return
// 	}
// 	if data.Total == 0 || data.List[0].Name != name {
// 		err = webapi.ErrNotExist
// 	} else {
// 		label.Id = data.List[0].Id
// 		label.Name = data.List[0].Name
// 		label.Color = LabelColor(webapi.LabelColorMap[data.List[0].Color])
// 	}
// 	return
// }

// LabelCreate creates a label with name and color, returns its ID.
func (a *Agent) LabelCreate(name string, color LabelColor) (labelId string, err error) {
	colorName, ok := labelColorMap[color]
	if !ok {
		colorName = api.LabelColorBlank
	}
	spec := (&api.LabelCreateSpec{}).Init(
		name, colorName,
	)
	if err = a.pc.ExecuteApi(spec); err != nil {
		return
	}
	if len(spec.Result) > 0 {
		labelId = spec.Result[0].Id
	}
	return
}

// LabelUpdate updates label's name or color.
func (a *Agent) LabelUpdate(label *Label) (err error) {
	if label == nil || label.Id == "" {
		return
	}
	colorName, ok := labelColorMap[label.Color]
	if !ok {
		colorName = api.LabelColorBlank
	}
	spec := (&api.LabelEditSpec{}).Init(
		label.Id, label.Name, colorName,
	)
	return a.pc.ExecuteApi(spec)
}

// LabelDelete deletes a label whose ID is labelId.
func (a *Agent) LabelDelete(labelId string) (err error) {
	if labelId == "" {
		return
	}
	spec := (&api.LabelDeleteSpec{}).Init(labelId)
	return a.pc.ExecuteApi(spec)
}

func (a *Agent) LabelSetOrder(labelId string, order FileOrder, asc bool) (err error) {
	spec := (&api.LabelSetOrderSpec{}).Init(
		labelId, getOrderName(order), asc,
	)
	return a.pc.ExecuteApi(spec)
}

// FileSetLabels sets labels for a file, you can also remove all labels from it
// by not passing any labelId.
func (a *Agent) FileSetLabels(fileId string, labelIds ...string) (err error) {
	spec := (&api.FileLabelSpec{}).Init(fileId, labelIds)
	return a.pc.ExecuteApi(spec)
}
