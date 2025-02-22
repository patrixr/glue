package blueprint

// @auteur("Concepts")
//
// # Blueprints
//
// Glue **blueprints** are a way to define a series of actions to be taken in a structured manner.
// Not only do they allow for the execution of complex workflows, the also enable us to view, serialize and share those
// workflows with or without executing them.
//

type ActionFunc func() error

type BlueprintFunc func() Trace

type Trace struct {
	Name       string `json:"name"`
	Details    string `json:"details"`
	Error      error  `json:"error"`
	Annotation string `json:"annotation"`
}

type Results struct {
	Traces         []Trace `json:"traces"`
	Success        bool    `json:"success"`
	ErrorCount     int     `json:"error_count"`
	TimeElapsedSec int     `json:"time_elapsed"`
}

type Blueprint interface {
	Execute() Results
	Action(name string, details string, usertext string, fn ActionFunc)
	Add(blueprint Blueprint)
	PrettyPrint() string
}
