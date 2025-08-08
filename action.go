package main

type Action interface {
}

type ActionHandler func(act Action, closed bool)
