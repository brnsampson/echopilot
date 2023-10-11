package server

type Service interface {
    Run() error
    Halt() error
    IsRunning() (bool, error)
    IsHalted() (bool, error)
}
