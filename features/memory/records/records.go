package records

import "github.com/brnsampson/echopilot/pkg/option"

func NewMemory(to_save string) *Memory {
    return &Memory { to_save }
}

type Memory struct {
    saved string
}

func (m *Memory) Display() string {
    return m.saved
}

type MemoryStore struct {
    saved []*Memory
}

func NewMemoryStore() *MemoryStore {
    memories := make([]*Memory, 0)

	return &MemoryStore { memories }
}

func (s *MemoryStore) Create(m *Memory) (*Memory, error) {
    s.saved = append(s.saved, m)
    return m, nil
}

//func (s *MemoryStore) Update(m *Memory) (*Memory, error) {
//    return m, nil
//}
//
//func (s *MemoryStore) Patch(m *Memory) (*Memory, error) {
//    return m, nil
//}
//
//func (s *MemoryStore) Get(m *Memory) *Memory {
//    return m
//}

func (s *MemoryStore) List(m option.Option[*Memory], page, offset int) []*Memory {
    // TODO: filter results based on passed memory option
    return s.saved
}

//func (s *MemoryStore) Delete(m *Memory) (*Memory, error) {
//    return m, nil
//}

