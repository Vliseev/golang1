// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package json

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson1982c6fcDecodeMailCourseGolangMailruCoursera3PerfomanceJson(in *jlexer.Lexer, out *Vasia) {
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
		case "Id":
			out.Id = int(in.Int())
		case "RealName":
			out.RealName = string(in.String())
		case "Login":
			out.Login = string(in.String())
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
func easyjson1982c6fcEncodeMailCourseGolangMailruCoursera3PerfomanceJson(out *jwriter.Writer, in Vasia) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Id))
	}
	{
		const prefix string = ",\"RealName\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.RealName))
	}
	{
		const prefix string = ",\"Login\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Login))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Vasia) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1982c6fcEncodeMailCourseGolangMailruCoursera3PerfomanceJson(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Vasia) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1982c6fcEncodeMailCourseGolangMailruCoursera3PerfomanceJson(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Vasia) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1982c6fcDecodeMailCourseGolangMailruCoursera3PerfomanceJson(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Vasia) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1982c6fcDecodeMailCourseGolangMailruCoursera3PerfomanceJson(l, v)
}
