package v2

type ClientPoller interface {
	ClientWithResponsesInterface

	OperationPoller(string, string) PollFunc
}
