package ast

import (
	"fmt"
	"strings"
)

type ValueType int

const (
	BooleanType ValueType = iota
	FunctionType
	NativeFunctionType
	NilType
	NumberType
	StringType
)

type Value interface {
	GetType() ValueType
	IsEqualTo(other Value) bool
	IsTruthy() bool
}

type BooleanValue struct {
	Value bool
}

func NewBooleanValue(v bool) *BooleanValue {
	return &BooleanValue{Value: v}
}

func (b *BooleanValue) GetType() ValueType {
	return BooleanType
}

func (b *BooleanValue) IsEqualTo(other Value) bool {
	if other == nil || other.GetType() != b.GetType() {
		return false
	}
	return b.Value == other.(*BooleanValue).Value
}

func (b *BooleanValue) IsTruthy() bool {
	return b.Value
}

func (b *BooleanValue) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

type NilValue struct{}

func NewNilValue() *NilValue {
	return &NilValue{}
}

func (n *NilValue) GetType() ValueType {
	return NilType
}

func (n *NilValue) IsEqualTo(other Value) bool {
	if other == nil || other.GetType() != n.GetType() {
		return false
	}
	return true
}

func (n *NilValue) IsTruthy() bool {
	return false
}

func (n *NilValue) String() string {
	return "nil"
}

type NumberValue struct {
	Value float64
}

func NewNumberValue(value float64) *NumberValue {
	return &NumberValue{Value: value}
}

func (n *NumberValue) GetType() ValueType {
	return NumberType
}

func (n *NumberValue) IsEqualTo(other Value) bool {
	if other == nil || other.GetType() != NumberType {
		return false
	}
	return n.Value == other.(*NumberValue).Value
}

func (n *NumberValue) IsTruthy() bool {
	return n.Value != 0
}

func (n *NumberValue) String() string {
	numStr := strings.TrimRight(fmt.Sprintf("%f", n.Value), "0")
	lastIdx := len(numStr) - 1
	if numStr[lastIdx] == uint8('.') {
		numStr = numStr[:lastIdx]
	}
	return numStr
}

type StringValue struct {
	Value string
}

func NewStringValue(value string) *StringValue {
	return &StringValue{Value: value}
}

func (s *StringValue) GetType() ValueType {
	return StringType
}

func (s *StringValue) IsEqualTo(other Value) bool {
	if other == nil || other.GetType() != StringType {
		return false
	}
	return s.Value == other.(*StringValue).Value
}

func (s *StringValue) IsTruthy() bool {
	return true
}

func (s *StringValue) String() string {
	return s.Value
}
