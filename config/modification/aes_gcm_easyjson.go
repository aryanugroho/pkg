// +build  csall json

// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package modification

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

func easyjson4cfa51e5DecodeGithubComCorestoreioPkgConfigModification(in *jlexer.Lexer, out *AESGCMOptions) {
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
		case "Key":
			out.Key = string(in.String())
		case "KeyEnvironmentVariableName":
			out.KeyEnvironmentVariableName = string(in.String())
		case "Nonce":
			if in.IsNull() {
				in.Skip()
				out.Nonce = nil
			} else {
				out.Nonce = in.Bytes()
			}
		case "NonceEnvironmentVariableName":
			out.NonceEnvironmentVariableName = string(in.String())
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
func easyjson4cfa51e5EncodeGithubComCorestoreioPkgConfigModification(out *jwriter.Writer, in AESGCMOptions) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Key != "" {
		const prefix string = ",\"Key\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Key))
	}
	if in.KeyEnvironmentVariableName != "" {
		const prefix string = ",\"KeyEnvironmentVariableName\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.KeyEnvironmentVariableName))
	}
	if len(in.Nonce) != 0 {
		const prefix string = ",\"Nonce\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Base64Bytes(in.Nonce)
	}
	if in.NonceEnvironmentVariableName != "" {
		const prefix string = ",\"NonceEnvironmentVariableName\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.NonceEnvironmentVariableName))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AESGCMOptions) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4cfa51e5EncodeGithubComCorestoreioPkgConfigModification(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AESGCMOptions) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4cfa51e5EncodeGithubComCorestoreioPkgConfigModification(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AESGCMOptions) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4cfa51e5DecodeGithubComCorestoreioPkgConfigModification(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AESGCMOptions) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4cfa51e5DecodeGithubComCorestoreioPkgConfigModification(l, v)
}
