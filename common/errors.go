package common

const (
	ERR_ACTION_ANSWER_INVALID_PAYLOAD = "attempted to %s answer with an invalid payload"
	ERR_ACTION_ANSWER_KEY             = "failed to %s answer with key %q: %v"
	ERR_ANSWER_CONFLICT               = "answer with key %q already exists"
	ERR_ANSWER_NOT_FOUND              = "answer with key %q not found"
	ERR_ANSWER_VALUE_EMPTY            = "answer value cannot be empty"
	ERR_EVENT_ANSWER_KEY              = "failed to get events with answer key %q: %v"
	ERR_EVENT_NOT_CREATED             = "could not create event of type %q with key %q: %v"
	ERR_EVENT_NOT_FOUND               = "event with key %q not found"
)
