package internal

func NewSavedEcho(value string) *SavedEcho {
    return &SavedEcho{ value }
}

type SavedEcho struct {
    value string
}


