package casset

import (
	"reflect"
	"testing"
)

func TestElement_Next(t *testing.T) {
	testMemory := NewMemory[any]()
	type args struct {
		current IElement[any]
		v       interface{}
	}
	tests := []struct {
		name string
		args args
		want IElement[any]
	}{
		{
			name: "next with value number",
			args: args{
				current: testMemory.GetFront(),
				v:       10,
			},
			want: &Element[any]{
				nextElement: nil,
				prevElement: testMemory.GetFront(),
				memory:      testMemory,
				value:       10,
			},
		},
		{
			name: "next with exist",
			args: args{
				current: testMemory.GetFront(),
				v:       1000,
			},
			want: &Element[any]{
				nextElement: nil,
				prevElement: testMemory.GetFront(),
				memory:      testMemory,
				value:       10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.current.Next(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Element.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_Prev(t *testing.T) {
	testMemory := NewMemory[any]()
	type args struct {
		current IElement[any]
		v       interface{}
	}
	tests := []struct {
		name string
		args args
		want IElement[any]
	}{
		{
			name: "prev with value number",
			args: args{
				current: testMemory.GetFront(),
				v:       10,
			},
			want: &Element[any]{
				nextElement: testMemory.GetBack(),
				prevElement: nil,
				memory:      testMemory,
				value:       10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.current.Prev(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Element.Prev() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_Delete(t *testing.T) {
	testMemory := NewMemory[any]()
	current := testMemory.GetFront().Next(1).Next(2).Next(3).Next(4).Prev(nil).Prev(nil)

	if testMemory.Len() != 5 {
		t.Errorf("Len problem")
	}

	current = current.Delete()
	testMemory.GetFront().Delete()
	testMemory.GetBack().Delete()

	if testMemory.Len() != 2 {
		t.Errorf("Len problem after delete")
	}

	want := &Element[any]{
		nextElement: nil,
		prevElement: testMemory.GetFront(),
		memory:      testMemory,
		value:       3,
	}

	if !reflect.DeepEqual(current, want) {
		t.Errorf("Element.Prev() = %v, want %v", current, want)
	}

	testMemory.GetBack().Delete()

	want = &Element[any]{
		nextElement: nil,
		prevElement: nil,
		memory:      testMemory,
		value:       1,
	}

	if !reflect.DeepEqual(testMemory.GetFront(), want) {
		t.Errorf("Element.Prev() = %v, want %v", current, want)
	}

	testMemory.GetFront().Delete()

	if testMemory.Len() != 1 {
		t.Errorf("Len problem after delete")
	}

	testMemory.Clear().GetFront().SetValue(2)

	want = &Element[any]{
		nextElement: nil,
		prevElement: nil,
		memory:      testMemory,
		value:       2,
	}

	if !reflect.DeepEqual(testMemory.GetFront(), want) {
		t.Errorf("Element.Prev() = %v, want %v", testMemory.GetFront(), want)
	}
}

func TestMemory_Remove(t *testing.T) {
	type args struct {
		e1 func(IMemory[any], IElement[any]) IElement[any]
		e2 func(IMemory[any], IElement[any]) IElement[any]
	}

	tests := []struct {
		name        string
		fields      func() (IMemory[any], uint64)
		current     func(IMemory[any]) IElement[any]
		args        args
		wantFront   func(IMemory[any], IElement[any]) (IElement[any], IElement[any])
		wantCurrent func(IMemory[any], IElement[any]) (IElement[any], IElement[any])
		wantBack    func(IMemory[any], IElement[any]) (IElement[any], IElement[any])
		wantLen     uint64
	}{
		{
			name: "remove all",
			fields: func() (IMemory[any], uint64) {
				m := NewMemory[any]()

				return m, 5
			},
			current: func(m IMemory[any]) IElement[any] {
				return m.GetFront().Next(1).Next(2).Next(3).Next(4).GetPrevElement()
			},
			args: args{
				e1: func(m IMemory[any], c IElement[any]) IElement[any] { return nil },
				e2: func(m IMemory[any], c IElement[any]) IElement[any] { return nil },
			},
			wantFront: func(m IMemory[any], c IElement[any]) (IElement[any], IElement[any]) {
				return nil, nil
			},
			wantCurrent: func(m IMemory[any], c IElement[any]) (IElement[any], IElement[any]) {
				return nil, nil
			},
			wantBack: func(m IMemory[any], c IElement[any]) (IElement[any], IElement[any]) {
				return nil, nil
			},
			wantLen: 1,
		},
		{
			name: "basic test",
			fields: func() (IMemory[any], uint64) {
				m := NewMemory[any]()

				return m, 5
			},
			current: func(m IMemory[any]) IElement[any] {
				return m.GetFront().Next(1).Next(2).Next(3).Next(4).GetPrevElement()
			},
			args: args{
				e1: func(m IMemory[any], c IElement[any]) IElement[any] { return nil },
				e2: func(m IMemory[any], c IElement[any]) IElement[any] {
					return c.GetPrevElement()
				},
			},
			wantCurrent: func(m IMemory[any], c IElement[any]) (IElement[any], IElement[any]) {
				return &Element[any]{
					nextElement: c.GetNextElement(),
					prevElement: nil,
					memory:      m,
					value:       3,
				}, c
			},
			wantLen: 2,
		},
		{
			name: "reverse test",
			fields: func() (IMemory[any], uint64) {
				m := NewMemory[any]()
				return m, 5
			},
			current: func(m IMemory[any]) IElement[any] {
				return m.GetFront().Next(1).Next(2).Next(3).Next(4).GetPrevElement()
			},
			args: args{
				e1: func(m IMemory[any], c IElement[any]) IElement[any] {
					return c
				},
				e2: func(m IMemory[any], c IElement[any]) IElement[any] { return m.GetFront() },
			},
			wantLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fields != nil {
				m, length := tt.fields()

				currentFn := tt.current
				if currentFn == nil {
					currentFn = func(IMemory[any]) IElement[any] { return nil }
				}
				current := tt.current(m)

				if m.Len() != length {
					t.Errorf("Len problem create %d, want %v", m.Len(), length)
				}

				wants := map[string]func(IMemory[any], IElement[any]) (IElement[any], IElement[any]){
					"Front":   tt.wantFront,
					"Current": tt.wantCurrent,
					"Back":    tt.wantBack,
				}

				for name, wantFn := range wants {
					if wantFn == nil {
						continue
					}

					want, check := wantFn(m, current)
					m.RemoveRange(tt.args.e1(m, current), tt.args.e2(m, current))

					if !reflect.DeepEqual(check, want) {
						t.Errorf("%s = %+v, want %+v", name, check, want)
					}

					if m.Len() != tt.wantLen {
						t.Errorf("Len problem after delete %d, want %v", m.Len(), tt.wantLen)
					}
				}
			}
		})
	}
}
