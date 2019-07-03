import test from "ava"

test.todo("WaitingRoom.constructor accepts a capacity value")
test.todo("WaitingRoom.constructor rejects capacity value < 0")
test.todo("WaitingRoom.constructor accepts a capacity value and timeout")
test.todo("WaitingRoom.constructor rejects timeout < 0")
test.todo("WaitingRoom.on(timeout) is emitted when timeout")
test.todo("WaitingRoom.join accepts Entrant")
test.todo("WaitingRoom.on(complete) is emitted when full")
test.todo("WaitingRoom.join rejects Entrant if capacity is exceeded")
test.todo("WaitingRoom.join rejects Entrant if timeout")
