package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

type ItemJson struct {
	Email    string
	Name     string
	Browsers []string
}

func FastSearch(out io.Writer) {

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	// VARIABLES
	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""

	// READ FILES TEXT
	if err != nil {
		panic(err)
	}
	// HEADER
	fmt.Fprintf(out,"found users:\n")
	// READER
	reader := bufio.NewScanner(file)

	// MARSHAL
	for i:=0; reader.Scan(); i++{
		user := ItemJson{}
		user.UnmarshalJSON(reader.Bytes())

		isAndroid := false
		isMSIE := false
		browsers := user.Browsers

		for _, browserRaw := range browsers {
			browser := browserRaw

			if ok := strings.Contains(browser, "Android"); ok && err == nil {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		for _, browserRaw := range browsers {
			browser := browserRaw
			if ok := strings.Contains(browser, "MSIE"); ok && err == nil {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		email := strings.Replace(user.Email, "@", " [at] ", -1)
		foundUsers = fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
		fmt.Fprintf(out,foundUsers)
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}

//////////////////////////////////////////////////////////////////////////////////////
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson3486653aDecodeHw3Bench(in *jlexer.Lexer, out *ItemJson) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func easyjson3486653aEncodeHw3Bench(out *jwriter.Writer, in ItemJson) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"browsers\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browsers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ItemJson) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3486653aEncodeHw3Bench(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ItemJson) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3486653aEncodeHw3Bench(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ItemJson) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3486653aDecodeHw3Bench(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ItemJson) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3486653aDecodeHw3Bench(l, v)
}
