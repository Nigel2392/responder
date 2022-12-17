package queryvalues

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"strings"
	"time"
)

type QueryValues struct {
	Values []*QueryValue
}

func (q *QueryValues) Add(name string, value string, typ string) *QueryValues {
	if name == "" || value == "" {
		return q
	}
	for _, qv := range q.Values {
		if qv.Name == name && qv.Type == typ {
			qv.Add(name, value)
			return q
		}
	}
	q.Values = append(q.Values, makeQV(name, value, typ))
	return q
}

func (q *QueryValues) RemoveIndex(name string, typ string, index int) (bool, bool) {
	if index > len(q.Values) || index < 0 {
		return false, false
	} else if len(q.Values) == 0 || len(q.Values) == 1 {
		if len(q.Values) == 1 {
			if q.Values[index].Name == name && q.Values[index].Type == typ {
				q.Values = make([]*QueryValue, 0)
				return true, true
			}
		}
		return false, true
	} else {
		if q.Values[index].Name == name && q.Values[index].Type == typ {
			q.Values = append(q.Values[:index], q.Values[index+1:]...)
			return true, false
		}
	}
	return false, false
}

func (q *QueryValues) GetByHash(h string) *QueryValue {
	for _, qv := range q.Values {
		if qv.HashString() == h {
			return qv
		}
	}
	return nil
}

func (q *QueryValues) DelByHash(h string, index ...int) (bool, bool) {
	if len(index) > 0 {
		for i, qv := range q.Values {
			if qv.HashString() == h {
				if len(qv.V) > 1 && index[0] < len(qv.V) {
					q.Values[i].V = append(qv.V[:index[0]], qv.V[index[0]+1:]...)
					return true, false
				}
			}
		}
	}
	for i, qv := range q.Values {
		if qv.HashString() == h {
			q.Values = append(q.Values[:i], q.Values[i+1:]...)
			return false, true
		}
	}
	return false, false
}

func (q *QueryValues) String() string {
	var b = strings.Builder{}
	for i, qv := range q.Values {
		for j, v := range qv.V {
			b.WriteString(qv.Name)
			b.WriteString("=")
			b.WriteString(v)
			if j < len(qv.V)-1 {
				b.WriteString("&")
			}
		}
		if i < len(q.Values)-1 {
			b.WriteString("&")
		}
	}
	return b.String()
}

func (q *QueryValues) StringByType(typ string) string {
	var b = strings.Builder{}
	for i, qv := range q.Values {
		if qv.Type == typ {
			for j, v := range qv.V {
				b.WriteString(qv.Name)
				b.WriteString("=")
				b.WriteString(v)
				if j < len(qv.V)-1 {
					b.WriteString("&")
				}
			}
			if i < len(q.Values)-1 {
				b.WriteString("&")
			}
			if i < len(q.Values)-1 && q.Values[i+1].Type == typ {
				b.WriteString("&")
			}
		}
	}
	return b.String()
}

type QueryValue struct {
	Name string
	V    []string
	Type string
}

func (q *QueryValue) Hash() uint32 {
	// Create a hash of the query value
	// This is used to identify the query value in the history

	var b = bytes.Buffer{}
	b.WriteString(q.Name)
	b.WriteString("=")
	for _, v := range q.V {
		b.WriteString(v)
	}
	b.WriteString("::")
	b.WriteString(q.Type)
	h := crc32.NewIEEE()
	io.Copy(h, &b)
	return h.Sum32()
}

func (q *QueryValue) HashString() string {
	return fmt.Sprintf("%x", q.Hash())
}

func makeQV(name string, value string, typ string) *QueryValue {
	qv := QueryValue{
		Type: typ,
		Name: name,
		V:    make([]string, 0),
	}
	qv.V = append(qv.V, value)
	return &qv
}

func (q *QueryValue) Add(name string, value string) {
	q.Name = name
	if q.V == nil {
		q.V = make([]string, 0)
	}
	for _, v := range q.V {
		if v == value {
			return
		}
	}
	q.V = append(q.V, value)
}

type History struct {
	Saves        []Save
	Fname        string
	MaximumItems int
}

type Saves History

type Save struct {
	QueryValues *QueryValues
	Method      string
	URL         string
	Timestamp   time.Time
}

func (h *History) Add(qv *QueryValues, method string, url string) {
	h.Saves = append(h.Saves, Save{qv, method, url, time.Now()})
}

func (h *History) Save() {
	fname := h.Fname
	if fname == "" {
		fname = "history.json"
	}
	oldhist := History{}
	oldhist.Load(fname)
	for _, save := range h.Saves {
		var found bool
		for _, oldsave := range oldhist.Saves {
			if save.QueryValues.String() == oldsave.QueryValues.String() &&
				save.Method == oldsave.Method &&
				save.URL == oldsave.URL &&
				save.Timestamp.Sub(oldsave.Timestamp) < time.Second/2 {
				found = true
			}
		}
		if !found {
			oldhist.Saves = append(oldhist.Saves, save)
		}
	}

	if h.MaximumItems > 0 && len(oldhist.Saves) > h.MaximumItems {
		// Keep only the last MaximumItems, not the first 20
		oldhist.Saves = oldhist.Saves[len(oldhist.Saves)-h.MaximumItems:]
	}

	json_bytes, err := json.MarshalIndent(oldhist, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile(fname, json_bytes, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (h *History) Load(fname string) string {
	json_bytes, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(json_bytes, h)
	if err != nil {
		fmt.Println(err)
	}
	return string(json_bytes)
}

func (h *History) Remove(qv *QueryValues, method string, url string, timeStamp time.Time) {
	for i, save := range h.Saves {
		if save.QueryValues.String() == qv.String() &&
			save.Method == method &&
			save.URL == url &&
			save.Timestamp.Sub(timeStamp) < time.Second/2 {
			h.Saves = append(h.Saves[:i], h.Saves[i+1:]...)
			return
		}
	}
}

func WailsEncode(v interface{}) (string, error) {
	json_bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(json_bytes), nil
}

func WailsDecode(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

func WailsEncodeB64(v interface{}) (string, error) {
	json_bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(json_bytes), nil
}

func WailsDecodeB64(data string, v interface{}) error {
	json_bytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(json_bytes, v)
}
