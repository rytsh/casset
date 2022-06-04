package casset

import (
	"testing"
)

func TestLen_Sub(t *testing.T) {
	type fields struct {
		current *Element
	}
	type args struct {
		val uint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint
	}{
		{
			name: "negatif",
			fields: fields{
				current: &Element{Value: uint(0)},
			},
			args: args{
				val: 5,
			},
			want: 0,
		},
		{
			name: "pozitif",
			fields: fields{
				current: &Element{Value: uint(200)},
			},
			args: args{
				val: 5,
			},
			want: 195,
		},
		{
			name: "more than max",
			fields: fields{
				current: func() *Element {
					c := &Element{
						Value: uint(3),
					}
					p := &Element{
						Value:       maxUint,
						NextElement: c,
					}
					c.PrevElement = p
					return c
				}(),
			},
			args: args{
				val: 5,
			},
			want: maxUint - 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Len{
				current: tt.fields.current,
			}
			l.Sub(tt.args.val)

			if tt.want != l.current.GetValue().(uint) {
				t.Errorf("want %d, current %d", tt.want, l.current.GetValue().(uint))
			}
		})
	}
}

func TestLen_Add(t *testing.T) {
	type fields struct {
		current *Element
	}
	type args struct {
		val uint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint
	}{
		{
			name: "pozitif",
			fields: fields{
				current: &Element{Value: uint(0)},
			},
			args: args{
				val: 5,
			},
			want: 5,
		},
		{
			name: "more than max",
			fields: fields{
				current: &Element{Value: maxUint - uint(3)},
			},
			args: args{
				val: 5,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Len{
				current: tt.fields.current,
			}
			l.Add(tt.args.val)

			if tt.want != l.current.GetValue().(uint) {
				t.Errorf("want %d, current %d", tt.want, l.current.GetValue().(uint))
			}
		})
	}
}

func TestLen_IsZero(t *testing.T) {
	type fields struct {
		current IElement
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "zero",
			fields: fields{
				current: &Element{Value: uint(0)},
			},
			want: true,
		},
		{
			name: "one",
			fields: fields{
				current: &Element{Value: uint(1)},
			},
			want: false,
		},
		{
			name: "more field",
			fields: fields{
				current: func() IElement { el := &Element{Value: uint(0)}; return el.Next(0) }(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Len{
				current: tt.fields.current,
			}
			if got := l.IsZero(); got != tt.want {
				t.Errorf("Len.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}
