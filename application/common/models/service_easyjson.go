// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

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

func easyjsonCd93bc43DecodeSubdApplicationCommonModels(in *jlexer.Lexer, out *ServiceStatus) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "user":
			out.UsersCnt = int64(in.Int64())
		case "forum":
			out.ForumsCnt = int64(in.Int64())
		case "thread":
			out.ThreadsCnt = int64(in.Int64())
		case "post":
			out.PostsCnt = int64(in.Int64())
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
func easyjsonCd93bc43EncodeSubdApplicationCommonModels(out *jwriter.Writer, in ServiceStatus) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.UsersCnt))
	}
	{
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		out.Int64(int64(in.ForumsCnt))
	}
	{
		const prefix string = ",\"thread\":"
		out.RawString(prefix)
		out.Int64(int64(in.ThreadsCnt))
	}
	{
		const prefix string = ",\"post\":"
		out.RawString(prefix)
		out.Int64(int64(in.PostsCnt))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ServiceStatus) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonCd93bc43EncodeSubdApplicationCommonModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ServiceStatus) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonCd93bc43EncodeSubdApplicationCommonModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ServiceStatus) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonCd93bc43DecodeSubdApplicationCommonModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ServiceStatus) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonCd93bc43DecodeSubdApplicationCommonModels(l, v)
}
