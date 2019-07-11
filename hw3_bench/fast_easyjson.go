// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.
//
package main
//
//import (
//	json "encoding/json"
//	easyjson "github.com/mailru/easyjson"
//	jlexer "github.com/mailru/easyjson/jlexer"
//	jwriter "github.com/mailru/easyjson/jwriter"
//)
//
//// suppress unused package warning
//var (
//	_ *json.RawMessage
//	_ *jlexer.Lexer
//	_ *jwriter.Writer
//	_ easyjson.Marshaler
//)
//
//func easyjson3486653aDecodeHw3Bench(in *jlexer.Lexer, out *ItemJson) {
//	isTopLevel := in.IsStart()
//	if in.IsNull() {
//		if isTopLevel {
//			in.Consumed()
//		}
//		in.Skip()
//		return
//	}
//	in.Delim('{')
//	for !in.IsDelim('}') {
//		key := in.UnsafeString()
//		in.WantColon()
//		if in.IsNull() {
//			in.Skip()
//			in.WantComma()
//			continue
//		}
//		switch key {
//		case "Email":
//			out.Email = string(in.String())
//		case "Name":
//			out.Name = string(in.String())
//		case "Browsers":
//			if in.IsNull() {
//				in.Skip()
//				out.Browsers = nil
//			} else {
//				in.Delim('[')
//				if out.Browsers == nil {
//					if !in.IsDelim(']') {
//						out.Browsers = make([]string, 0, 4)
//					} else {
//						out.Browsers = []string{}
//					}
//				} else {
//					out.Browsers = (out.Browsers)[:0]
//				}
//				for !in.IsDelim(']') {
//					var v1 string
//					v1 = string(in.String())
//					out.Browsers = append(out.Browsers, v1)
//					in.WantComma()
//				}
//				in.Delim(']')
//			}
//		default:
//			in.SkipRecursive()
//		}
//		in.WantComma()
//	}
//	in.Delim('}')
//	if isTopLevel {
//		in.Consumed()
//	}
//}
//func easyjson3486653aEncodeHw3Bench(out *jwriter.Writer, in ItemJson) {
//	out.RawByte('{')
//	first := true
//	_ = first
//	{
//		const prefix string = ",\"Email\":"
//		if first {
//			first = false
//			out.RawString(prefix[1:])
//		} else {
//			out.RawString(prefix)
//		}
//		out.String(string(in.Email))
//	}
//	{
//		const prefix string = ",\"Name\":"
//		if first {
//			first = false
//			out.RawString(prefix[1:])
//		} else {
//			out.RawString(prefix)
//		}
//		out.String(string(in.Name))
//	}
//	{
//		const prefix string = ",\"Browsers\":"
//		if first {
//			first = false
//			out.RawString(prefix[1:])
//		} else {
//			out.RawString(prefix)
//		}
//		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
//			out.RawString("null")
//		} else {
//			out.RawByte('[')
//			for v2, v3 := range in.Browsers {
//				if v2 > 0 {
//					out.RawByte(',')
//				}
//				out.String(string(v3))
//			}
//			out.RawByte(']')
//		}
//	}
//	out.RawByte('}')
//}
//
//// MarshalJSON supports json.Marshaler interface
//func (v ItemJson) MarshalJSON() ([]byte, error) {
//	w := jwriter.Writer{}
//	easyjson3486653aEncodeHw3Bench(&w, v)
//	return w.Buffer.BuildBytes(), w.Error
//}
//
//// MarshalEasyJSON supports easyjson.Marshaler interface
//func (v ItemJson) MarshalEasyJSON(w *jwriter.Writer) {
//	easyjson3486653aEncodeHw3Bench(w, v)
//}
//
//// UnmarshalJSON supports json.Unmarshaler interface
//func (v *ItemJson) UnmarshalJSON(data []byte) error {
//	r := jlexer.Lexer{Data: data}
//	easyjson3486653aDecodeHw3Bench(&r, v)
//	return r.Error()
//}
//
//// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
//func (v *ItemJson) UnmarshalEasyJSON(l *jlexer.Lexer) {
//	easyjson3486653aDecodeHw3Bench(l, v)
//}
