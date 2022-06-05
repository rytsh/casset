package casset

import (
	"reflect"
	"testing"
)

func TestElement_Next(t *testing.T) {
	testMemory := NewMemory(nil)
	type args struct {
		current *Element
		v       interface{}
	}
	tests := []struct {
		name string
		args args
		want *Element
	}{
		{
			name: "next with value number",
			args: args{
				current: testMemory.Current.(*Element),
				v:       10,
			},
			want: &Element{
				NextElement: nil,
				PrevElement: testMemory.Front.(*Element),
				Memory:      testMemory,
				Value:       10,
			},
		},
		{
			name: "next with exist",
			args: args{
				current: testMemory.Front.(*Element),
				v:       1000,
			},
			want: &Element{
				NextElement: nil,
				PrevElement: testMemory.Front.(*Element),
				Memory:      testMemory,
				Value:       10,
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
	testMemory := NewMemory(nil)
	type args struct {
		current *Element
		v       interface{}
	}
	tests := []struct {
		name string
		args args
		want *Element
	}{
		{
			name: "prev with value number",
			args: args{
				current: testMemory.Current.(*Element),
				v:       10,
			},
			want: &Element{
				NextElement: testMemory.Back.(*Element),
				PrevElement: nil,
				Memory:      testMemory,
				Value:       10,
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
	testMemory := NewMemory(0)
	testMemory.Current = testMemory.Current.Next(1).Next(2).Next(3).Next(4).Prev(nil).Prev(nil)

	if testMemory.len.Cmp(5) != 0 {
		t.Errorf("Len problem")
	}

	testMemory.Current.Delete()
	testMemory.Front.Delete()
	testMemory.Back.Delete()

	if testMemory.len.Cmp(2) != 0 {
		t.Errorf("Len problem after delete")
	}

	want := &Element{
		NextElement: nil,
		PrevElement: testMemory.Front.(*Element),
		Memory:      testMemory,
		Value:       3,
	}

	if !reflect.DeepEqual(testMemory.Current, want) {
		t.Errorf("Element.Prev() = %v, want %v", testMemory.Current, want)
	}

	if err := testMemory.Back.Delete(); err != nil {
		t.Error(err)
	}

	want = &Element{
		NextElement: nil,
		PrevElement: nil,
		Memory:      testMemory,
		Value:       1,
	}

	if !reflect.DeepEqual(testMemory.Current, want) {
		t.Errorf("Element.Prev() = %v, want %v", testMemory.Current, want)
	}

	testMemory.Current.Delete()

	if testMemory.len.Cmp(0) != 0 {
		t.Errorf("Len problem after delete")
	}

	testMemory.Init(&Element{Memory: testMemory, Value: 2})

	want = &Element{
		NextElement: nil,
		PrevElement: nil,
		Memory:      testMemory,
		Value:       2,
	}

	if !reflect.DeepEqual(testMemory.Current, want) {
		t.Errorf("Element.Prev() = %v, want %v", testMemory.Current, want)
	}
}

func TestMemory_Remove(t *testing.T) {
	type args struct {
		e1 func(*Memory) IElement
		e2 func(*Memory) IElement
	}

	tests := []struct {
		name        string
		fields      func() (*Memory, uint)
		args        args
		wantFront   func(*Memory) (IElement, IElement)
		wantCurrent func(*Memory) (IElement, IElement)
		wantBack    func(*Memory) (IElement, IElement)
		wantLen     uint
	}{
		{
			name: "remove all",
			fields: func() (*Memory, uint) {
				m := NewMemory(0)
				m.Current = m.Current.Next(1).Next(2).Next(3).Next(4).GetPrevElement()
				return m, 5
			},
			args: args{
				e1: func(m *Memory) IElement { return nil },
				e2: func(m *Memory) IElement { return nil },
			},
			wantFront: func(m *Memory) (IElement, IElement) {
				return nil, nil
			},
			wantCurrent: func(m *Memory) (IElement, IElement) {
				return nil, nil
			},
			wantBack: func(m *Memory) (IElement, IElement) {
				return nil, nil
			},
			wantLen: 0,
		},
		{
			name: "basic test",
			fields: func() (*Memory, uint) {
				m := NewMemory(0)
				m.Current = m.Current.Next(1).Next(2).Next(3).Next(4).GetPrevElement()
				return m, 5
			},
			args: args{
				e1: func(m *Memory) IElement { return nil },
				e2: func(m *Memory) IElement { return m.Current.GetPrevElement() },
			},
			wantCurrent: func(m *Memory) (IElement, IElement) {
				return &Element{
					NextElement: m.Current.GetNextElement().(*Element),
					PrevElement: nil,
					Memory:      m,
					Value:       3,
				}, m.Current
			},
			wantLen: 2,
		},
		// {
		// 	name: "reverse test",
		// 	fields: func() (*Memory, uint) {
		// 		m := NewMemory(0)
		// 		m.Current = m.Current.Next(1).Next(2).Next(3).Next(4).GetPrevElement()
		// 		return m, 5
		// 	},
		// 	args: args{
		// 		e1: func(m *Memory) IElement { return m.Current },
		// 		e2: func(m *Memory) IElement { return m.Front },
		// 	},
		// 	wantCurrent: func(m *Memory) (IElement, IElement) {
		// 		return &Element{
		// 			NextElement: nil,
		// 			PrevElement: m.Current.GetPrevElement().GetPrevElement().(*Element),
		// 			Memory:      m,
		// 			Value:       2,
		// 		}, m.Current
		// 	},
		// 	wantLen: 2,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fields != nil {
				m, length := tt.fields()
				if m.len.Cmp(length) != 0 {
					t.Errorf("Len problem create %v, want %v", m.len.GetValueCurrent(), length)
				}

				wants := map[string]func(*Memory) (IElement, IElement){
					"Front":   tt.wantFront,
					"Current": tt.wantCurrent,
					"Back":    tt.wantBack,
				}

				for name, wantFn := range wants {
					if wantFn == nil {
						continue
					}

					want, check := wantFn(m)
					m.Remove(tt.args.e1(m), tt.args.e2(m))

					if !reflect.DeepEqual(check, want) {
						t.Errorf("%s = %+v, want %+v", name, check, want)
					}

					if m.len.Cmp(tt.wantLen) != 0 {
						t.Errorf("Len problem after delete %v, want %v", m.len.GetValueCurrent(), tt.wantLen)
					}
				}
			}
		})
	}
}
