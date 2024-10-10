package fsm

type StateOutputTuple struct {
	S State
	O string
}

type (
	State           string
	Event           string
	StateTransition map[Event]StateOutputTuple
	StateMap        map[State]StateTransition
)

// type FiniteStateMachine interface {
// 	Current() State
// 	Transition(Event) error
// 	Reset()
// }
//
// func FeedThrough(fsm FiniteStateMachine, events []Event) (State, error) {
// 	for _, e := range events {
// 		err := fsm.Transition(e)
// 		if err != nil {
// 			return "", err
// 		}
// 	}
// 	return fsm.Current(), nil
// }
//